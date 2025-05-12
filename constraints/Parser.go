package constraints

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	utils "github.com/CodeClarityCE/utility-node-semver/utils"
	version "github.com/CodeClarityCE/utility-node-semver/versions"

	"slices"
)

var (
	ErrInvalidConstraint            = errors.New("invalid Semantic Version")
	ErrEmptyConstraint              = errors.New("empty constraint")
	ErrIllegalCharacterInConstraint = errors.New("illegal character in constraint")
	ErrInvalidVersion               = errors.New("invalid Version")
)

type Range struct {
	StartOp      Token
	StartVersion version.Semver
	EndOp        Token
	EndVersion   version.Semver
}

type JoinOp string

const (
	CONJUNCTON  JoinOp = "and"
	DISJUNCTION JoinOp = "or"
)

type SubConstraint struct {
	Ranges Range
}

type Constraint struct {
	Original string
	Ranges   []Range
	Join     []JoinOp
}

func (c *Constraint) String() string {
	stringRep := ""
	joinLength := len(c.Join)
	for idx, r := range c.Ranges {
		if idx < joinLength {
			stringRep = stringRep + r.String() + string(c.Join[idx])
		} else {
			stringRep = stringRep + r.String()
		}
	}
	return stringRep
}

// Parses a given node semver constraint string into a constraint object
func ParseConstraint(constraintString string) (Constraint, error) {
	tokens, literals := LexConstraint(constraintString)

	// Check for illegal tokens
	if slices.Contains(tokens, ILLEGAL) || slices.Contains(tokens, OPEN_PARENTHESIS) || slices.Contains(tokens, CLOSE_PARENTHESIS) {
		indicies := []int{}
		indicies = append(indicies, utils.GetAllIndicies(tokens, ILLEGAL)...)
		indicies = append(indicies, utils.GetAllIndicies(tokens, OPEN_PARENTHESIS)...)
		indicies = append(indicies, utils.GetAllIndicies(tokens, CLOSE_PARENTHESIS)...)
		err := newIllegalTokenInConstraint(fmt.Sprintf("Found illegal token in the tokenized string.\n\tHere: %s", getConstraintErrorString(tokens, literals, indicies)))
		log.Println(err)
		return Constraint{}, err
	}

	// Check if the constraint is correctly composed
	err := validateConstraintComposition(tokens, literals)
	if err != nil {
		log.Println(err)
		return Constraint{}, err
	}

	constraint := Constraint{
		Original: constraintString,
	}

	// Each sub constraint is "desugared" into a simple range. e.g. >= 5.0.0 =< 6.0.0
	// This is because all other operators: ~, ^, .x, Any, *, are simply syntactic sugar for a range

	currentTokenList := []Token{}
	currentLiteralList := []string{}
	for idx, token := range tokens {
		literal := literals[idx]
		if token == OR || token == AND {
			subConstraint, err := parseSubConstraint(currentTokenList, currentLiteralList)
			if err != nil {
				return Constraint{}, err
			}
			if token == AND {
				constraint.Join = append(constraint.Join, CONJUNCTON)
			}
			if token == OR {
				constraint.Join = append(constraint.Join, DISJUNCTION)
			}
			constraint.Ranges = append(constraint.Ranges, subConstraint)
			currentTokenList = []Token{}
			currentLiteralList = []string{}
		} else if token == EOF {
			subConstraint, err := parseSubConstraint(currentTokenList, currentLiteralList)
			if err != nil {
				return Constraint{}, err
			}
			constraint.Ranges = append(constraint.Ranges, subConstraint)
			currentTokenList = []Token{}
			currentLiteralList = []string{}
		} else {
			if token != SOF {
				currentTokenList = append(currentTokenList, token)
				currentLiteralList = append(currentLiteralList, literal)
			}
		}
	}

	return constraint, nil
}

func validateConstraintComposition(tokens []Token, literals []string) error {

	joinOperatorsProcessed := []Token{}
	tokensProcessed := []Token{}
	previousToken := SOF
	for idx, token := range tokens {

		// A constraint that is empty is not valid
		// if token == EOF && previousToken == SOF {
		// 	return EmptyConstraint
		// }

		// A version constraint cannot start with a join token
		//   e.g. && 5.0.0 || 2.0.0
		if IsJoinToken(token) && previousToken == SOF {
			return newErrInvalidConstraint(fmt.Sprintf("Found join operator (&&, ||) at the start of the constraint.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		}

		// A version constraint cannot end with a join token
		//   e.g. 1.x || >= 5.0.0 &&
		if token == EOF && IsJoinToken(previousToken) {
			return newErrInvalidConstraint(fmt.Sprintf("Found join operator (&&, ||) at the end of the constraint.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		}

		// A constraint cannot contain two subsequent range tokens or equality tokens
		if (IsRangeToken(token) || IsEqualityToken(token)) && (IsRangeToken(previousToken) || IsEqualityToken(previousToken)) {
			return newErrInvalidConstraint(fmt.Sprintf("Found two subsequent range operators (>=, >, <, <=, -, ~, ^) or equality operators (!,=).\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		}

		// A constraint cannot contain two subsequent join tokens
		if IsJoinToken(token) && IsJoinToken(previousToken) {
			return newErrInvalidConstraint(fmt.Sprintf("Found two subsequent join operators (&&, ||).\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		}

		// A constraint cannot contain two subsequent static version expressions
		//
		// where static version means a simple version that does not have .x .* or ANY
		if (token == VERSION_EXPRESSION && version.IsStaticVersion(literals[idx])) && (previousToken == VERSION_EXPRESSION && version.IsStaticVersion(literals[idx-1])) {
			return newErrInvalidConstraint(fmt.Sprintf("Found two subsequent static version expression that are not joined by join operator (&& , ||), or a range operator (>=, >, <, <=, -, ~, ^).\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		}

		// We cannot have a version expression without a preceeding operator token (join, equality or range)
		//
		// Exception: hyphens
		//   e.g. 5.0.0 - 7.0.0 is a valid expression even through 5.0.0 does not have an operator in front of it
		//
		// Exception: .x, ANY, * expressions
		//   e.g 1.x is a valid expression even though 1.x does not have an operator in front of it
		//           that is because 1.x is really just syntactic sugar for >= 1.0.0
		if token == VERSION_EXPRESSION && !IsOperatorToken(previousToken) {

			previousToken := SOF
			if idx >= 1 {
				previousToken = tokens[idx-1]
			}

			nextToken := EOF
			if idx+1 <= len(tokens)-1 {
				nextToken = tokens[idx+1]
			}

			if previousToken != HYPHEN && nextToken != HYPHEN && version.IsStaticVersion(literals[idx]) {
				return newErrInvalidConstraint(fmt.Sprintf("Found a version expression without a preceeding (>=, >, <, <=, -, ~, ^) range, equality (=, !) or join operator (&&, ||).\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
			}
		}

		// Two subsequent version expressions must be joined by join or an range token
		//  e.g. 4.0.0 || 5.0.0
		//  e.g. >= 2.0.0 && =< 4.0.0
		//  e.g. >= 2.0.0 =< 4.0.0
		if len(tokensProcessed) >= 3 && token == VERSION_EXPRESSION && tokensProcessed[idx-2] == VERSION_EXPRESSION {
			if (!IsRangeToken(previousToken) && !IsJoinToken(previousToken)) && !IsRangeToken(tokensProcessed[idx-3]) && !IsJoinToken(tokensProcessed[idx-3]) {
				return newErrInvalidConstraint(fmt.Sprintf("Found two subsequent version expressions that are not joined by join (&&, ||)) or range operator (>=, >, <, <=, -, ~, ^).\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
			}
		}

		// Two (or more) subsequent "static" version expressions cannot be joined by a conjunction
		// if the operator of both version expressions is an equlity operator
		//  e.g. = 4.0.0 && = 5.0.0 is not allowed
		//  e.g. ! 4.0.0 && = 5.0.0 is not allowed
		//  but  = 4.0.0 || = 5.0.0 is allowed
		//
		// where static version means a simple version that does not have .x .* or ANY
		if token == VERSION_EXPRESSION && (previousToken == EQ || previousToken == NOT) {
			idxEq := slices.Index(tokensProcessed[0:len(tokensProcessed)-1], EQ)

			if idxEq == -1 {
				idxEq = slices.Index(tokensProcessed[0:len(tokensProcessed)-1], NOT)
			}

			if idxEq != -1 {
				versionLiteral := literals[idx]
				versionLiteralBefore := literals[idxEq+1]
				if version.IsStaticVersion(versionLiteral) && version.IsStaticVersion(versionLiteralBefore) {
					if len(joinOperatorsProcessed) > 0 {
						if joinOperatorsProcessed[len(joinOperatorsProcessed)-1] == AND {
							return newErrInvalidConstraint(fmt.Sprintf("Found two static versions joined by &&, which is a logical impossiblity.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idxEq + 1, idx})))
						}
					}
				}
			}
		}

		// // A valid join must consist of two versions, a LHS version and a RHS version
		// if IsJoinToken(token) {

		// 	// Join expression found, where there is no LHS version
		// 	if previousToken != VERSION_EXPRESSION {
		// 		return newErrInvalidConstraint(fmt.Sprintf("Found join (&&, ||) expression that does not have a LHS version.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		// 	}

		// 	// Join expression found, where there is no RHS version
		// 	nextToken := EOF
		// 	if idx+1 <= len(tokens)-1 {
		// 		nextToken = tokens[idx+1]
		// 	}
		// 	if nextToken != VERSION_EXPRESSION {
		// 		return newErrInvalidConstraint(fmt.Sprintf("Found join (&&, ||) expression that does not have a RHS version.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
		// 	}

		// }

		// A valid hyphenated range must consist of two static or partial versions, a start version and an end version joined in the middle by a hyphen
		if token == HYPHEN {

			// Hyphen found but no static or partial start version
			if previousToken != VERSION_EXPRESSION {
				return newErrInvalidConstraint(fmt.Sprintf("Found hyphenated range that is not prefixed with a version.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
			}

			// Hyphen and start version found, but start version is not a static or partial version
			if !version.IsPartialVersion(literals[idx-1]) && !version.IsStaticVersion(literals[idx-1]) {
				return newErrInvalidConstraint(fmt.Sprintf("Found hyphenated range that is not prefixed with a static or partial version.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
			}

			// Hyphen and start version found, but start version is not a static or partial version
			if idx > 3 && (IsRangeToken(tokens[idx-2]) || IsEqualityToken(tokens[idx-2])) {
				return newErrInvalidConstraint(fmt.Sprintf("Found out of place range (>=, >, <, <=, -, ~, ^) or equality token (=, !) in hyphenated range.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx, idx - 4})))
			}

			// Hyphen found but no end version
			nextToken := EOF
			if idx+1 <= len(tokens)-1 {
				nextToken = tokens[idx+1]
			}

			if nextToken != VERSION_EXPRESSION {
				return newErrInvalidConstraint(fmt.Sprintf("Found hyphenated range that is not complete. End version is missing.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
			}

			// Hyphen, start and end version found, but end version is not a static or partial version
			if nextToken == VERSION_EXPRESSION {
				if !version.IsPartialVersion(literals[idx+1]) && !version.IsStaticVersion(literals[idx+1]) {
					return newErrInvalidConstraint(fmt.Sprintf("Found hyphenated range that is not suffixed with a static or partial version. End version is not static.\n\tHere: %s", getConstraintErrorString(tokens, literals, []int{idx})))
				}
			}

		}

		if IsJoinToken(token) {
			joinOperatorsProcessed = append(joinOperatorsProcessed, token)
		}

		tokensProcessed = append(tokensProcessed, token)
		previousToken = token
	}

	return nil

}

func parseSubConstraint(tokens []Token, literals []string) (Range, error) {
	switch {
	case isHyphenatedRange(tokens):
		// https://github.com/npm/codeclarity.io/node-semver#hyphen-ranges-xyz---abc
		semverRange, err := parseHyphen(literals)
		return semverRange, err
	case isCaretRange(tokens):
		// https://github.com/npm/codeclarity.io/node-semver#caret-ranges-123-025-004
		semverRange, err := parseCaretRange(literals)
		return semverRange, err
	case isTildeRange(tokens):
		// https://github.com/npm/codeclarity.io/node-semver#tilde-ranges-123-12-1
		semverRange, err := parseTildeRange(literals)
		return semverRange, err
	case isRange(tokens):
		semverRange, err := parseRange(tokens, literals)
		return semverRange, err
	case isStaticRange(tokens, literals):
		semverRange, err := parseStaticRange(tokens, literals)
		return semverRange, err
	case isPartialRange(tokens, literals):
		semverRange, err := parseXRange(literals)
		return semverRange, err
	case isXRange(tokens, literals):
		// https://github.com/npm/codeclarity.io/node-semver#x-ranges-12x-1x-12-
		semverRange, err := parseXRange(literals)
		return semverRange, err
	default:
		return Range{}, errors.New("unknown sub constraint type")
	}

}

// func getSubConstraints(tokens []Token, literals []string, delimiter Token, omitEOFAndSOF bool) ([][]Token, [][]string) {
// 	tokensToReturn := [][]Token{}
// 	literalsToReturn := [][]string{}

// 	tokenGroup := []Token{}
// 	literalGroup := []string{}

// 	lengthTokens := len(tokens)

// 	for idx, token := range tokens {

// 		if token == SOF {
// 			if !omitEOFAndSOF {
// 				tokenGroup = append(tokenGroup, SOF)
// 				literalGroup = append(literalGroup, "")
// 			}
// 		} else if token == EOF {
// 			if !omitEOFAndSOF {
// 				tokenGroup = append(tokenGroup, EOF)
// 				literalGroup = append(literalGroup, "")
// 			}
// 		} else if token != delimiter {
// 			tokenGroup = append(tokenGroup, token)
// 			literalGroup = append(literalGroup, literals[idx])
// 		}

// 		if token == delimiter || idx == lengthTokens-1 {
// 			tokensToReturn = append(tokensToReturn, tokenGroup)
// 			literalsToReturn = append(literalsToReturn, literalGroup)
// 			tokenGroup = []Token{}
// 			literalGroup = []string{}
// 		}

// 	}
// 	return tokensToReturn, literalsToReturn
// }

func isHyphenatedRange(tokenList []Token) bool {
	return len(tokenList) == 3 && tokenList[0] == VERSION_EXPRESSION && tokenList[1] == HYPHEN && tokenList[2] == VERSION_EXPRESSION
}

func isXRange(tokenList []Token, literalList []string) bool {
	return len(tokenList) == 0 || (len(tokenList) == 1 && tokenList[0] == VERSION_EXPRESSION) || literalList[0] == "*" || literalList[0] == "ANY" || strings.Contains(literalList[0], ".x") || strings.Contains(literalList[0], ".X") || strings.Contains(literalList[0], ".*")
}

func isCaretRange(tokenList []Token) bool {
	return len(tokenList) == 2 && tokenList[0] == CARET && tokenList[1] == VERSION_EXPRESSION
}

func isTildeRange(tokenList []Token) bool {
	return len(tokenList) == 2 && tokenList[0] == TILDE && tokenList[1] == VERSION_EXPRESSION
}

func isRange(tokenList []Token) bool {
	if len(tokenList) == 2 && IsRangeToken(tokenList[0]) && tokenList[1] == VERSION_EXPRESSION {
		return true
	}
	if len(tokenList) == 4 && IsRangeToken(tokenList[0]) && tokenList[1] == VERSION_EXPRESSION && IsRangeToken(tokenList[2]) && tokenList[3] == VERSION_EXPRESSION {
		return true
	}
	return false
}

func isStaticRange(tokenList []Token, literalList []string) bool {
	return len(tokenList) == 2 && IsEqualityToken(tokenList[0]) && tokenList[1] == VERSION_EXPRESSION && version.IsStaticVersion(literalList[1])
}

func isPartialRange(tokenList []Token, literalList []string) bool {
	return len(tokenList) == 1 && tokenList[0] == VERSION_EXPRESSION && version.IsPartialVersion(literalList[0])
}

func parseHyphen(literalList []string) (Range, error) {

	// https://github.com/npm/codeclarity.io/node-semver#hyphen-ranges-xyz---abcs.
	// Specifies an inclusive set.

	parsedRange := Range{}
	parsedRange.StartOp = GE
	parsedRange.EndOp = LE

	startVersion := version.GetVersionPart(literalList[0])
	endVersion := version.GetVersionPart(literalList[2])

	startMetaDataPart := version.GetMetaDataPart(literalList[0])
	startPreReleasePart := version.GetPreReleasePart(literalList[0])

	endMetaDataPart := version.GetMetaDataPart(literalList[2])
	endPreReleasePart := version.GetPreReleasePart(literalList[2])

	startVersionString := ""
	endVersionString := ""

	// Specifies an inclusive set.
	//   e.g 1.2.3 - 2.3.4 	:= 	>=1.2.3 <=2.3.4
	if version.IsStaticVersion(startVersion) && version.IsStaticVersion(endVersion) {
		startVersionString = getVersionStringFromMetaPreReleaseParts(startVersion, startMetaDataPart, startPreReleasePart)
		endVersionString = getVersionStringFromMetaPreReleaseParts(endVersion, endMetaDataPart, endPreReleasePart)
	}

	// If a partial version is provided as the first version in the inclusive range, then the missing pieces are replaced with zeroes.
	//   e.g. 1.2 - 2.3.4 	:= 	>=1.2.0 <=2.3.4
	if version.IsPartialVersion(startVersion) && version.IsStaticVersion(endVersion) {
		partialStartVersion := startVersion
		if strings.Count(partialStartVersion, ".") == 0 {
			partialStartVersion += ".0.0"
		}
		if strings.Count(partialStartVersion, ".") == 1 {
			partialStartVersion += ".0"
		}

		startVersionString = getVersionStringFromMetaPreReleaseParts(partialStartVersion, startMetaDataPart, startPreReleasePart)
		endVersionString = getVersionStringFromMetaPreReleaseParts(endVersion, endMetaDataPart, endPreReleasePart)
	}

	// If a partial version is provided as the second version in the inclusive range, then all versions that start with the supplied parts
	// of the tuple are accepted, but nothing that would be greater than the provided tuple parts.
	//   e.g. 1.2.3 - 2.3 	:= 	>=1.2.3 <2.4.0-0
	//   e.g. 1.2.3 - 2 	:=	>=1.2.3 <3.0.0-0
	if version.IsStaticVersion(startVersion) && version.IsPartialVersion(endVersion) {

		partialEndVersion := endVersion
		coercedPartialEndRange := partialEndVersion
		parsedRange.EndOp = LT

		if strings.Count(partialEndVersion, ".") == 0 {
			parsed, err := strconv.Atoi(partialEndVersion)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := int(parsed)
			major++
			coercedPartialEndRange = fmt.Sprintf("%d.0.0", major)
		}

		if strings.Count(partialEndVersion, ".") == 1 {
			parsed, err := strconv.Atoi(strings.Split(partialEndVersion, ".")[1])
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := strings.Split(partialEndVersion, ".")[0]
			minor := int(parsed)
			minor++
			coercedPartialEndRange = fmt.Sprintf("%s.%d.0", major, minor)
		}

		startVersionString = getVersionStringFromMetaPreReleaseParts(startVersion, startMetaDataPart, startPreReleasePart)
		endVersionString = getVersionStringFromMetaPreReleaseParts(coercedPartialEndRange, endMetaDataPart, endPreReleasePart)
	}

	if version.IsPartialVersion(startVersion) && version.IsPartialVersion(endVersion) {

		partialEndVersion := endVersion
		coercedPartialEndRange := partialEndVersion
		parsedRange.EndOp = LT
		partialStartVersion := startVersion
		if strings.Count(partialStartVersion, ".") == 0 {
			partialStartVersion += ".0.0"
		}
		if strings.Count(partialStartVersion, ".") == 1 {
			partialStartVersion += ".0"
		}

		if strings.Count(partialEndVersion, ".") == 0 {
			parsed, err := strconv.Atoi(partialEndVersion)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := int(parsed)
			major++
			coercedPartialEndRange = fmt.Sprintf("%d.0.0", major)
		}

		if strings.Count(partialEndVersion, ".") == 1 {
			parsed, err := strconv.Atoi(strings.Split(partialEndVersion, ".")[1])
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := strings.Split(partialEndVersion, ".")[0]
			minor := int(parsed)
			minor++
			coercedPartialEndRange = fmt.Sprintf("%s.%d.0", major, minor)
		}

		startVersionString = getVersionStringFromMetaPreReleaseParts(partialStartVersion, startMetaDataPart, startPreReleasePart)
		endVersionString = getVersionStringFromMetaPreReleaseParts(coercedPartialEndRange, endMetaDataPart, endPreReleasePart)

	}

	semver, err := version.ParseSemver(startVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.StartVersion = semver
	semver, err = version.ParseSemver(endVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.EndVersion = semver

	return parsedRange, nil
}

func parseTildeRange(literalList []string) (Range, error) {

	// https://github.com/npm/codeclarity.io/node-semver#tilde-ranges-123-12-1
	// Allows patch-level changes if a minor version is specified on the comparator. Allows minor-level changes if not.

	parsedRange := Range{}
	parsedRange.StartOp = GE

	versionString := version.GetVersionPart(literalList[1])
	metaDataPart := version.GetMetaDataPart(literalList[1])
	preReleasePart := version.GetPreReleasePart(literalList[1])

	split := strings.Split(versionString, ".")
	majorString := ""
	minorString := ""
	patchString := ""

	startVersionString := ""
	endVersionString := ""

	if len(split) == 3 {
		patchString = strings.ToLower(split[2])
	}
	if len(split) >= 2 {
		minorString = strings.ToLower(split[1])
	}
	if len(split) >= 1 {
		majorString = strings.ToLower(split[0])
	}

	// ~*, ~Any is a special case := >= 0.0.0
	if utils.ContainsOnly(split, "*") ||
		utils.ContainsOnly(split, "x") ||
		utils.ContainsOnly(split, "X") ||
		versionString == "ANY" {
		semver, err := version.ParseSemver("0.0.0")
		if err != nil {
			return Range{}, err
		}
		parsedRange.StartVersion = semver
		return parsedRange, nil
	}

	parsedRange.EndOp = LT

	if patchString == "x" || patchString == "*" {
		patchString = ""
		versionString = fmt.Sprintf("%s.%s", majorString, minorString)
	}

	if minorString == "x" || minorString == "*" {
		minorString = ""
		versionString = majorString
	}

	// ~1.2.3 	:= >=1.2.3 <1.(2+1).0 	:= >=1.2.3 <1.3.0-0
	// ~0.2.3 	:= >=0.2.3 <0.(2+1).0 	:= >=0.2.3 <0.3.0-0
	if strings.Count(versionString, ".") == 2 {
		parsed, err := strconv.Atoi(minorString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		minor := int(parsed)
		minor++
		startVersionString = getVersionStringFromMetaPreReleaseParts(versionString, metaDataPart, preReleasePart)
		endVersionString = fmt.Sprintf("%s.%d.0", majorString, minor)
	}

	// ~1.2 	:= >=1.2.0 <1.(2+1).0 	:= >=1.2.0 <1.3.0-0
	// ~0.2 	:= >=0.2.0 <0.(2+1).0 	:= >=0.2.0 <0.3.0-0
	if strings.Count(versionString, ".") == 1 {
		parsed, err := strconv.Atoi(minorString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		minor := int(parsed)
		minor++
		startVersionString = getVersionStringFromMetaPreReleaseParts(fmt.Sprintf("%s.0", versionString), metaDataPart, preReleasePart)
		endVersionString = fmt.Sprintf("%s.%d.0", majorString, minor)
	}

	// ~1 	:= >=1.0.0 <(1+1).0.0 	:= >=1.0.0 <2.0.0-0
	// ~0 	:= >=0.0.0 <(0+1).0.0 	:= >=0.0.0 <1.0.0-0
	if strings.Count(versionString, ".") == 0 {
		parsed, err := strconv.Atoi(versionString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		major := int(parsed)
		major++
		startVersionString = getVersionStringFromMetaPreReleaseParts(fmt.Sprintf("%s.0.0", versionString), metaDataPart, preReleasePart)
		endVersionString = fmt.Sprintf("%d.0.0", major)
	}

	semver, err := version.ParseSemver(startVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.StartVersion = semver
	semver, err = version.ParseSemver(endVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.EndVersion = semver

	return parsedRange, nil
}

// ^*      -->  (any)
// ^1.2.3  -->  >=1.2.3 <2.0.0
// ^1.2    -->  >=1.2.0 <2.0.0
// ^1      -->  >=1.0.0 <2.0.0
// ^0.2.3  -->  >=0.2.3 <0.3.0
// ^0.2    -->  >=0.2.0 <0.3.0
// ^0.0.3  -->  >=0.0.3 <0.0.4
// ^0.0    -->  >=0.0.0 <0.1.0
// ^0      -->  >=0.0.0 <1.0.0
func parseCaretRange(literalList []string) (Range, error) {

	// https://github.com/npm/codeclarity.io/node-semver#caret-ranges-123-025-004
	// Allows changes that do not modify the left-most non-zero element in the [major, minor, patch] tuple

	parsedRange := Range{}
	parsedRange.StartOp = GE

	versionString := version.GetVersionPart(literalList[1])
	metaDataPart := version.GetMetaDataPart(literalList[1])
	preReleasePart := version.GetPreReleasePart(literalList[1])

	split := strings.Split(versionString, ".")
	majorString := ""
	minorString := ""
	patchString := ""

	startVersionString := ""
	endVersionString := ""

	if len(split) == 3 {
		patchString = strings.ToLower(split[2])
	}
	if len(split) >= 2 {
		minorString = strings.ToLower(split[1])
	}
	if len(split) >= 1 {
		majorString = strings.ToLower(split[0])
	}

	// ^*, ^Any is a special case := >= 0.0.0
	if utils.ContainsOnly(split, "*") ||
		utils.ContainsOnly(split, "x") ||
		utils.ContainsOnly(split, "X") ||
		versionString == "ANY" {
		parsedRange.StartOp = GE
		semver, err := version.ParseSemver("0.0.0")
		if err != nil {
			return Range{}, err
		}
		parsedRange.StartVersion = semver
		return parsedRange, nil
	}

	if patchString == "*" || patchString == "" {
		patchString = "x"
	}
	if minorString == "*" || minorString == "" {
		minorString = "x"
	}
	if majorString == "*" || majorString == "" {
		majorString = "x"
	}

	parsedRange.EndOp = LT

	// versionString = getVersionStringFromParts(majorString, minorString, patchString, metaDataPart, preReleasePart)

	// A missing minor and patch values will desugar to zero, but also
	// allow flexibility within those values, even if the major version is zero.
	//   e.g. ^1.x 	:= 	>=1.0.0 <2.0.0-0
	//   e.g. ^0.x 	:= 	>=0.0.0 <1.0.0-0
	if minorString == "x" && patchString == "x" {
		parsed, err := strconv.Atoi(majorString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		major := int(parsed)
		major++

		startVersionString = getVersionStringFromParts(majorString, "0", "0", metaDataPart, preReleasePart)
		endVersionString = getVersionStringFromParts(fmt.Sprintf("%d", major), "0", "0", "", "")

	} else if patchString == "x" {

		startVersionString = fmt.Sprintf("%s.%s.0", majorString, minorString)

		// When parsing caret ranges, a missing patch value desugars to the number 0, but will
		// allow flexibility within that value, even if the major and minor versions are both 0.
		//   e.g. ^1.2.x 	:= 	>=1.2.0 <2.0.0-0
		//   e.g. ^0.0.x 	:= 	>=0.0.0 <0.1.0-0
		//   e.g. ^0.0 		:= 	>=0.0.0 <0.1.0-0
		if majorString != "0" {
			parsed, err := strconv.Atoi(majorString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := int(parsed)
			major++
			endVersionString = getVersionStringFromParts(fmt.Sprintf("%d", major), "0", "0", "", "")
		} else if minorString != "0" {
			parsed, err := strconv.Atoi(minorString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			minor := int(parsed)
			minor++
			endVersionString = getVersionStringFromParts(majorString, fmt.Sprintf("%d", minor), "0", "", "")
		} else {
			parsed, err := strconv.Atoi(minorString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			minor := int(parsed)
			minor++
			endVersionString = getVersionStringFromParts(majorString, fmt.Sprintf("%d", minor), "0", "", "")
		}
	} else {

		startVersionString = getVersionStringFromParts(majorString, minorString, patchString, metaDataPart, preReleasePart)

		// ^1.2.3 	:= 	>=1.2.3 <2.0.0-0
		// ^0.2.3 	:= 	>=0.2.3 <0.3.0-0
		// ^0.0.3 	:= 	>=0.0.3 <0.0.4-0
		if majorString != "0" {
			parsed, err := strconv.Atoi(majorString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			major := int(parsed)
			major++
			endVersionString = getVersionStringFromParts(fmt.Sprintf("%d", major), "0", "0", "", "")
		} else if minorString != "0" {
			parsed, err := strconv.Atoi(minorString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			minor := int(parsed)
			minor++
			endVersionString = getVersionStringFromParts(majorString, fmt.Sprintf("%d", minor), "0", "", "")
		} else {
			parsed, err := strconv.Atoi(patchString)
			if err != nil {
				return Range{}, ErrInvalidVersion
			}
			patch := int(parsed)
			patch++
			endVersionString = getVersionStringFromParts(majorString, minorString, fmt.Sprintf("%d", patch), "", "")
		}
	}

	semver, err := version.ParseSemver(startVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.StartVersion = semver
	semver, err = version.ParseSemver(endVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.EndVersion = semver

	return parsedRange, nil

}

func parseStaticRange(tokenList []Token, literalList []string) (Range, error) {

	parsedRange := Range{}
	parsedRange.StartOp = tokenList[0]

	startVersionString := parsePartsOrZero(literalList[1])

	semver, err := version.ParseSemver(startVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.StartVersion = semver

	return parsedRange, nil

}

func parsePartsOrZero(versionString string) string {

	originalVersion := versionString
	versionString = version.GetVersionPart(originalVersion)
	metaDataPart := version.GetMetaDataPart(originalVersion)
	preReleasePart := version.GetPreReleasePart(originalVersion)

	split := strings.Split(versionString, ".")
	majorString := "0"
	minorString := "0"
	patchString := "0"

	if len(split) == 3 {
		patchString = strings.ToLower(split[2])
	}
	if len(split) >= 2 {
		minorString = strings.ToLower(split[1])
	}
	if len(split) >= 1 {
		majorString = strings.ToLower(split[0])
	}

	if metaDataPart != "" && preReleasePart != "" {
		return fmt.Sprintf("%s.%s.%s+%s-%s", majorString, minorString, patchString, metaDataPart, preReleasePart)
	} else if metaDataPart != "" {
		return fmt.Sprintf("%s.%s.%s+%s", majorString, minorString, patchString, metaDataPart)
	} else if preReleasePart != "" {
		return fmt.Sprintf("%s.%s.%s-%s", majorString, minorString, patchString, preReleasePart)
	} else {
		return fmt.Sprintf("%s.%s.%s", majorString, minorString, patchString)
	}

}

func parseRange(tokenList []Token, literalList []string) (Range, error) {

	parsedRange := Range{}

	if len(literalList) == 2 {
		parsedRange.StartOp = tokenList[0]
		startVersionString := parsePartsOrZero(literalList[1])
		startVersionString = strings.ReplaceAll(startVersionString, "x", "0")
		startVersionString = strings.ReplaceAll(startVersionString, "X", "0")
		startVersionString = strings.ReplaceAll(startVersionString, "*", "0")
		semver, err := version.ParseSemver(startVersionString)
		if err != nil {
			return Range{}, err
		}
		parsedRange.StartVersion = semver
		return parsedRange, nil
	} else if len(literalList) == 4 {
		parsedRange.StartOp = tokenList[0]
		startVersionString := parsePartsOrZero(literalList[1])
		startVersionString = strings.ReplaceAll(startVersionString, "x", "0")
		startVersionString = strings.ReplaceAll(startVersionString, "X", "0")
		startVersionString = strings.ReplaceAll(startVersionString, "*", "0")
		semver, err := version.ParseSemver(startVersionString)
		if err != nil {
			return Range{}, err
		}
		parsedRange.StartVersion = semver
		parsedRange.EndOp = tokenList[2]
		endVersionString := parsePartsOrZero(literalList[3])
		endVersionString = strings.ReplaceAll(endVersionString, "x", "0")
		endVersionString = strings.ReplaceAll(endVersionString, "X", "0")
		endVersionString = strings.ReplaceAll(endVersionString, "*", "0")
		semver, err = version.ParseSemver(endVersionString)
		if err != nil {
			return Range{}, err
		}
		parsedRange.EndVersion = semver
		return parsedRange, nil
	}

	return Range{}, errors.New("invalid range")

}

func parseXRange(literalList []string) (Range, error) {

	// https://github.com/npm/codeclarity.io/node-semver#x-ranges-12x-1x-12-
	// Any of X, x, or * may be used to "stand in" for one of the numeric values in the [major, minor, patch] tuple.
	//   e.g. * 	:= 	>=0.0.0
	//   e.g. 1.x 	:= 	>=1.0.0 <2.0.0-0
	//   e.g. 1.2.x := 	>=1.2.0 <1.3.0-0
	//   e.g. "" 	:= 	>=0.0.0
	//   e.g. 1 	:=  >=1.0.0 <2.0.0-0
	//   e.g. 1.2 	:=	>=1.2.0 <1.3.0-0

	parsedRange := Range{}
	parsedRange.StartOp = GE
	versionString := ""

	if len(literalList) > 0 {
		versionString = version.GetVersionPart(literalList[0])
	}

	split := strings.Split(versionString, ".")
	majorString := ""
	minorString := ""
	patchString := ""

	if len(split) == 3 {
		patchString = strings.ToLower(split[2])
	}
	if len(split) >= 2 {
		minorString = strings.ToLower(split[1])
	}
	if len(split) >= 1 {
		majorString = strings.ToLower(split[0])
	}

	startVersionString := ""
	endVersionString := ""

	if versionString == "" || versionString == "*" || versionString == "ANY" ||
		utils.ContainsOnly(split, "*") ||
		utils.ContainsOnly(split, "x") ||
		utils.ContainsOnly(split, "X") {
		semver, err := version.ParseSemver("0.0.0")
		if err != nil {
			return Range{}, err
		}
		parsedRange.StartVersion = semver
		return parsedRange, nil
	}

	if minorString == "" || minorString == "x" || minorString == "*" {
		startVersionString = fmt.Sprintf("%s.0.0", majorString)
		parsed, err := strconv.Atoi(majorString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		major := int(parsed)
		major++
		endVersionString = fmt.Sprintf("%d.0.0", major)
	} else if patchString == "" || patchString == "x" || patchString == "*" {
		startVersionString = fmt.Sprintf("%s.%s.0", majorString, minorString)
		parsed, err := strconv.Atoi(minorString)
		if err != nil {
			return Range{}, ErrInvalidVersion
		}
		minor := int(parsed)
		minor++
		endVersionString = fmt.Sprintf("%s.%d.0", majorString, minor)
	}

	semver, err := version.ParseSemver(startVersionString)
	if err != nil {
		return Range{}, err
	}
	parsedRange.StartVersion = semver

	if endVersionString != "" {
		semver, err = version.ParseSemver(endVersionString)
		if err != nil {
			return Range{}, err
		}
		parsedRange.EndOp = LT
		parsedRange.EndVersion = semver
	}

	return parsedRange, nil

}

func getConstraintErrorString(tokens []Token, literals []string, errorPositions []int) string {
	formattedStrings := []string{}
	for idx, token := range tokens {
		if !slices.Contains(errorPositions, idx) {
			formattedStrings = append(formattedStrings, fmt.Sprintf("`%s<'%s'>`", token, literals[idx]))
		} else {
			formattedStrings = append(formattedStrings, fmt.Sprintf("===> `%s<'%s'>` <===", token, literals[idx]))
		}
	}
	return fmt.Sprintf("%s\n\tIn: %s", strings.Join(formattedStrings, " • "), strings.Join(literals, ""))
}

func newErrInvalidConstraint(message string) error {
	return fmt.Errorf("invalid Version Constraint: %s", message)
}

func newIllegalTokenInConstraint(message string) error {
	return fmt.Errorf("invalid Version Constraint: %s", message)
}

func (parsedRange Range) String() string {
	return fmt.Sprintf("%s %s %s %s", parsedRange.StartOp.toString(), parsedRange.StartVersion.String(), parsedRange.EndOp.toString(), parsedRange.EndVersion.String())

}

func getVersionStringFromParts(majorString string, minorString string, patchString string, metaDataPart string, preReleasePart string) string {
	if metaDataPart != "" && preReleasePart != "" {
		return fmt.Sprintf("%s.%s.%s+%s-%s", majorString, minorString, patchString, metaDataPart, preReleasePart)
	} else if metaDataPart != "" {
		return fmt.Sprintf("%s.%s.%s+%s", majorString, minorString, patchString, metaDataPart)
	} else if preReleasePart != "" {
		return fmt.Sprintf("%s.%s.%s-%s", majorString, minorString, patchString, preReleasePart)
	} else {
		return fmt.Sprintf("%s.%s.%s", majorString, minorString, patchString)
	}
}

func getVersionStringFromMetaPreReleaseParts(version string, metaDataPart string, preReleasePart string) string {
	if metaDataPart != "" && preReleasePart != "" {
		return fmt.Sprintf("%s+%s-%s", version, metaDataPart, preReleasePart)
	} else if metaDataPart != "" {
		return fmt.Sprintf("%s+%s", version, metaDataPart)
	} else if preReleasePart != "" {
		return fmt.Sprintf("%s-%s", version, preReleasePart)
	} else {
		return version
	}
}
