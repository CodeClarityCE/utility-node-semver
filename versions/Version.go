package versions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Semver struct {
	Major         int
	Minor         int
	Patch         int
	PreReleaseTag string
	MetaData      string
	
	// Composer-specific fields
	IsDev     bool   // true for dev versions (dev-master, 1.0.x-dev)
	DevBranch string // branch name for dev-* versions
	Stability string // stability flag (@stable, @RC, @beta, @alpha, @dev)
}

var (
	ErrInvalidVersionParts = errors.New("invalid version parts")
	ErrInvalidPreRelease   = errors.New("invalid pre release")
	ErrInvalidMetaData     = errors.New("invalid meta data")
)

// ParseSemverWithEcosystem parses a semver string into a semver object for specified ecosystem
func ParseSemverWithEcosystem(versionLiteral string, ecosystem string) (Semver, error) {
	if versionLiteral == "" {
		return Semver{}, ErrInvalidVersionParts
	}

	semver := Semver{}

	// Handle Composer-specific dev versions
	if ecosystem == "composer" {
		if strings.HasPrefix(versionLiteral, "dev-") {
			semver.IsDev = true
			semver.DevBranch = strings.TrimPrefix(versionLiteral, "dev-")
			// Dev versions are considered maximum versions for their branch
			semver.Major = 999
			semver.Minor = 999
			semver.Patch = 999
			return semver, nil
		}

		if strings.HasSuffix(versionLiteral, "-dev") {
			semver.IsDev = true
			versionLiteral = strings.TrimSuffix(versionLiteral, "-dev")
		}

		// Handle stability flags (@stable, @RC, etc.)
		if strings.Contains(versionLiteral, "@") {
			parts := strings.Split(versionLiteral, "@")
			if len(parts) == 2 {
				versionLiteral = parts[0]
				semver.Stability = parts[1]
			}
		}
	}

	// Remove 'v' prefix if present
	versionLiteral = strings.TrimPrefix(versionLiteral, "v")

	versionParts := strings.Split(GetVersionPart(versionLiteral), ".")
	
	// Handle x.y.x patterns for Composer
	if ecosystem == "composer" {
		for i, part := range versionParts {
			if part == "x" || part == "X" {
				versionParts[i] = "999" // Treat wildcards as high numbers
			}
		}
	}

	if len(versionParts) != 3 {
		if len(versionParts) == 2 {
			for i, part := range versionParts {
				parsed, err := strconv.Atoi(part)
				if err != nil {
					return Semver{}, ErrInvalidVersionParts
				}

				switch i {
				case 0:
					semver.Major = parsed
				case 1:
					semver.Minor = parsed
				}

				semver.Patch = 0
			}
		} else {
			return Semver{}, ErrInvalidVersionParts
		}
	} else {
		for i, part := range versionParts {
			parsed, err := strconv.Atoi(part)
			if err != nil {
				return Semver{}, ErrInvalidVersionParts
			}

			switch i {
			case 0:
				semver.Major = parsed
			case 1:
				semver.Minor = parsed
			case 2:
				semver.Patch = parsed
			}
		}
	}

	semver.MetaData = GetMetaDataPart(versionLiteral)
	semver.PreReleaseTag = GetPreReleasePart(versionLiteral)

	return semver, nil
}

// Parses a semver string into a semver object
// DEPRECATED: Use ParseSemverWithEcosystem for new code. Defaults to NodeJS ecosystem.
func ParseSemver(versionLiteral string) (Semver, error) {
	return ParseSemverWithEcosystem(versionLiteral, "nodejs")
}

// Compares versions v1 and v2, and returns true if v1 >= v2 and false otherwise
// This comparision is done according to the semver 2.0 spec with Composer extensions
func (v1 Semver) GE(v2 Semver, ignorePreRelease bool) bool {
	// Handle dev version comparisons for Composer
	if v1.IsDev && v2.IsDev {
		// Both dev versions - compare by branch name
		if v1.DevBranch != v2.DevBranch {
			return v1.DevBranch >= v2.DevBranch
		}
		return true // Same dev branch is equal
	}

	if v1.IsDev && !v2.IsDev {
		// Dev version is higher than stable version of same base
		return true
	}

	if !v1.IsDev && v2.IsDev {
		// Stable version is lower than dev version
		return false
	}

	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor > v2.Minor
	}

	if v1.Patch != v2.Patch {
		return v1.Patch > v2.Patch
	}

	// both are equal in terms of [major, minor, patch] at this point
	if ignorePreRelease {
		return true
	}

	// As per semver spec:
	//   When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
	//
	//   5.0.0 >= 5.0.0-beta.5 returns true since 5.0.0 is considered greater than 5.0.0-beta.5
	if v1.PreReleaseTag == "" && v2.PreReleaseTag != "" {
		return true
	} else if v1.PreReleaseTag != "" && v2.PreReleaseTag == "" {
		return false
	}

	preReleaseComparison := comparePreRelease(v1.PreReleaseTag, v2.PreReleaseTag)
	return preReleaseComparison >= 0
}

// Compares versions v1 and v2, and returns true if v1 > v2 and false otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) GT(v2 Semver, ignorePreRelease bool) bool {
	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor > v2.Minor
	}

	if v1.Patch != v2.Patch {
		return v1.Patch > v2.Patch
	}

	// both are equal in terms of [major, minor, patch] at this point
	if ignorePreRelease {
		return false
	}

	// As per semver spec:
	//   When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
	//
	//   5.0.0 > 5.0.0-beta.5 returns true since 5.0.0 is considered greater than 5.0.0-beta.5
	if v1.PreReleaseTag == "" && v2.PreReleaseTag != "" {
		return true
	} else if v1.PreReleaseTag != "" && v2.PreReleaseTag == "" {
		return false
	}

	preReleaseComparison := comparePreRelease(v1.PreReleaseTag, v2.PreReleaseTag)
	return preReleaseComparison == 1
}

// Compares versions v1 and v2, and returns true if v1 <= v2 and false otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) LE(v2 Semver, ignorePreRelease bool) bool {
	if v1.Major != v2.Major {
		return v1.Major < v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor < v2.Minor
	}

	if v1.Patch != v2.Patch {
		return v1.Patch < v2.Patch
	}

	// both are equal in terms of [major, minor, patch] at this point
	if ignorePreRelease {
		return true
	}

	// As per semver spec:
	//   When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
	//
	//   5.0.0 < 5.0.0-beta.5 returns false since 5.0.0 is considered greater than 5.0.0-beta.5
	if v1.PreReleaseTag == "" && v2.PreReleaseTag != "" {
		return false
	} else if v1.PreReleaseTag != "" && v2.PreReleaseTag == "" {
		return true
	}

	preReleaseComparison := comparePreRelease(v1.PreReleaseTag, v2.PreReleaseTag)
	return preReleaseComparison <= 0
}

// Compares versions v1 and v2, and returns true if v1 < v2 and false otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) LT(v2 Semver, ignorePreRelease bool) bool {
	if v1.Major != v2.Major {
		return v1.Major < v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor < v2.Minor
	}

	if v1.Patch != v2.Patch {
		return v1.Patch < v2.Patch
	}

	// both are equal in terms of [major, minor, patch] at this point
	if ignorePreRelease {
		return false
	}

	// As per semver spec:
	//   When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
	//
	//   5.0.0 < 5.0.0-beta.5 returns false since 5.0.0 is considered greater than 5.0.0-beta.5
	if v1.PreReleaseTag == "" && v2.PreReleaseTag != "" {
		return false
	} else if v1.PreReleaseTag != "" && v2.PreReleaseTag == "" {
		return true
	}

	preReleaseComparison := comparePreRelease(v1.PreReleaseTag, v2.PreReleaseTag)
	return preReleaseComparison == -1
}

// Compares versions v1 and v2, and returns true if v1 = v2 and false otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) EQ(v2 Semver, ignorePreRelease bool) bool {
	if v1.Major != v2.Major || v1.Minor != v2.Minor || v1.Patch != v2.Patch {
		return false
	}

	if ignorePreRelease {
		return true
	}

	if (v1.PreReleaseTag == "" && v2.PreReleaseTag != "") || (v1.PreReleaseTag != "" && v2.PreReleaseTag == "") {
		return false
	}

	preReleaseComparison := comparePreRelease(v1.PreReleaseTag, v2.PreReleaseTag)
	return preReleaseComparison == 0
}

// Compares versions v1 and v2, and returns true if v1 != v2 and false otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) NEQ(v2 Semver, includePreReleases bool) bool {
	return !v1.EQ(v2, includePreReleases)
}

// Compares versions v1 and v2, and returns 0 if v1 = v2, returns 1 if v1 > v2 and 0 otherwise
// This comparision is done according to the semver 2.0 spec
func (v1 Semver) Compare(v2 Semver, ignorePreRelease bool) int {
	if v1.EQ(v2, ignorePreRelease) {
		return 0
	}
	if v1.GT(v2, ignorePreRelease) {
		return 1
	}
	return -1
}

// Returns a string representation of the parsed semver
// Note: the metadata and prerelease info might not be in the same order
// as the original string that was parsed
func (v Semver) String() string {

	versionString := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.MetaData != "" {
		versionString += fmt.Sprintf("+%s", v.MetaData)
	}

	if v.PreReleaseTag != "" {
		versionString += fmt.Sprintf("-%s", v.PreReleaseTag)
	}

	return versionString

}

func GetMetaDataPart(versionLiteral string) string {
	if !strings.Contains(versionLiteral, "+") {
		return ""
	}
	if versionLiteral != "" {
		versionsArr := strings.SplitN(versionLiteral, "+", 2)
		if len(versionsArr) > 1 {
			versionLiteral = versionsArr[1]
		} else {
			return ""
		}
	}
	if versionLiteral != "" {
		versionsArr := strings.SplitN(versionLiteral, "-", 2)
		if len(versionsArr) > 0 {
			versionLiteral = versionsArr[0]
		} else {
			return ""
		}
	}
	return versionLiteral
}

func GetPreReleasePart(versionLiteral string) string {
	if !strings.Contains(versionLiteral, "-") {
		return ""
	}

	versionsArr := strings.SplitN(versionLiteral, string('-'), 2)
	if len(versionsArr) < 2 {
		return ""
	}

	versionLiteral = versionsArr[1]

	versionsArr = strings.SplitN(versionLiteral, string('+'), 2)
	if len(versionsArr) == 0 {
		return ""
	}

	versionLiteral = versionsArr[0]

	return versionLiteral
}

func GetVersionPart(versionLiteral string) string {
	if !strings.Contains(versionLiteral, "-") && !strings.Contains(versionLiteral, "+") {
		return versionLiteral
	}
	versionsArr := strings.SplitN(versionLiteral, string('-'), 2)
	if len(versionsArr) > 0 {
		versionLiteral = versionsArr[0]
	}
	versionsArr = strings.SplitN(versionLiteral, string('+'), 2)
	if len(versionsArr) > 0 {
		versionLiteral = versionsArr[0]
	}
	return versionLiteral
}

func comparePreRelease(preRelease1 string, preRelease2 string) int {
	// According to the semver spec
	// "Precedence for two pre-release versions with the same major, minor, and patch version
	// MUST be determined by comparing each dot separated identifier from left to right
	// until a difference is found as follows:
	//   (1) Identifiers consisting of only digits are compared numerically.
	//   (2) Identifiers with letters or hyphens are compared lexically in ASCII sort order.
	//   (3) Numeric identifiers always have lower precedence than non-numeric identifiers.
	//   (4) A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal."

	// The key word being: "Precedence ... MUST be determined by comparing each dot separated identifier"
	parts1 := strings.Split(preRelease1, ".")
	parts2 := strings.Split(preRelease2, ".")
	lenParts1 := len(parts1)
	lenParts2 := len(parts2)

	iterationCount := lenParts1
	if lenParts2 > lenParts1 {
		iterationCount = lenParts2
	}

	for i := 0; i < iterationCount; i++ {
		iden1 := ""
		iden2 := ""

		if i < lenParts1 {
			iden1 = parts1[i]
		}
		if i < lenParts2 {
			iden2 = parts2[i]
		}

		comp := comparePreReleaseIdentifier(iden1, iden2)
		if comp != 0 {
			return comp
		}
	}

	return 0
}

func comparePreReleaseIdentifier(preRelease1Identifier string, preRelease2Identifier string) int {
	if preRelease1Identifier == preRelease2Identifier {
		return 0
	}

	// According to the semver spec:
	// "A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal."
	// Thus trivially if one is empty but the other is not the other is greater
	if preRelease1Identifier == "" && preRelease2Identifier != "" {
		return -1
	}

	if preRelease1Identifier != "" && preRelease2Identifier == "" {
		return 1
	}

	isIntPre1Id, err1 := strconv.Atoi(preRelease1Identifier)
	isIntPre2Id, err2 := strconv.Atoi(preRelease2Identifier)

	if err1 == nil && err2 == nil {
		// If both are ints we compare them numerically
		if isIntPre1Id > isIntPre2Id {
			return 1
		} else if isIntPre1Id < isIntPre2Id {
			return -1
		}
		return 0
	}

	// According to the semver spec:
	// "Numeric identifiers always have lower precedence than non-numeric identifiers."
	// Thus trivially if one is an int but the other is not, then the other is greater
	if err1 == nil && err2 != nil {
		return -1
	}

	if err1 != nil && err2 == nil {
		return 1
	}

	// If neither are ints, then according to the semver spec:
	// "Identifiers with letters or hyphens are compared lexically in ASCII sort order."
	return strings.Compare(preRelease1Identifier, preRelease2Identifier)

}

func IsStaticVersion(versionLiteral string) bool {
	versionLiteral = GetVersionPart(versionLiteral)
	return strings.Count(versionLiteral, ".") == 2 && !IsWildCardVersion(versionLiteral)
}

func IsPartialVersion(versionLiteral string) bool {
	versionLiteral = GetVersionPart(versionLiteral)
	return strings.Count(versionLiteral, ".") < 2 && !IsWildCardVersion(versionLiteral)
}

func IsWildCardVersion(versionLiteral string) bool {
	versionLiteral = GetVersionPart(versionLiteral)
	return strings.Count(versionLiteral, "x") > 0 || strings.Count(versionLiteral, "X") > 0 || strings.Count(versionLiteral, "*") > 0 || versionLiteral == "ANY" || versionLiteral == "*"
}
