package constraints

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"

	versions "github.com/CodeClarityCE/utility-node-semver/versions"

	"slices"
)

func LexConstraint(constraint string) (tokens []Token, literals []string) {
	lexer := newLexer(strings.NewReader(constraint))
	return lexer.ScanWhole()
}

var operatorsStarts = []rune{'=', '<', '>', '!', '&', '|', '-', '^', '~', '(', ')'}
var versionUnsafe = []rune{'=', '<', '>', '!', '&', '|', '^', '~', '(', ')'}

// eof represents a marker rune for the end of the reader.
var eof = rune(0)

type Lexer struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func newLexer(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

func (lexer Lexer) ScanWhole() (tokens []Token, literals []string) {
	tokens = []Token{SOF}
	literals = []string{""}

	for {
		token, literal := lexer.ScanOne()

		if token != WS {
			tokens = append(tokens, token)
			literals = append(literals, literal)
		}

		if token == EOF {
			break
		}
	}

	// Trim the token list from prefix WS and suffix WS
	// i.e. remove whitespace directly behind SOF and whitespace directly infront of EOF (if any)
	if len(tokens) > 1 {
		if tokens[1] == WS {
			tokens = tokens[2:]
			tokens = append([]Token{SOF}, tokens...)
			literals = literals[2:]
			literals = append([]string{""}, literals...)
		}
	}
	if len(tokens) > 1 {
		if tokens[len(tokens)-2] == WS {
			tokens = tokens[0 : len(tokens)-2]
			tokens = append(tokens, EOF)
			literals = literals[0 : len(literals)-2]
			literals = append(literals, "")
		}
	}

	return removeSuperfluousAndOp(augmentMissingEqualityOp(tokens, literals))
}

func (lexer Lexer) ScanOne() (token Token, literal string) {
	// Read the next rune
	ch := lexer.read()

	// If we see whitespace then consume all contiguous whitespace.
	if isWhitespace(ch) {
		lexer.unread()
		return lexer.scanWhitespace()
	} else if isDigit(ch) {
		lexer.unread()
		return lexer.scanVersionExpression()
	} else if isOperatorStart(ch) {
		lexer.unread()
		return lexer.scanOperator()
	} else if isLetter(ch) || ch == '*' {
		lexer.unread()
		return lexer.scanVersionExpression()
	}

	switch ch {
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (lexer Lexer) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := lexer.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			lexer.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (lexer Lexer) scanVersionExpression() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := lexer.read(); ch == eof {
			break
		} else if isWhitespace(ch) || isVersionUnsafe(ch) {
			lexer.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return VERSION_EXPRESSION, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (lexer Lexer) scanOperator() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := lexer.read(); ch == eof {
			break
		} else if !isOperatorStart(ch) {
			lexer.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
			if ch == '(' || ch == ')' {
				break
			}
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "&&":
		return AND, buf.String()
	case "||":
		return OR, buf.String()
	case "=":
		return EQ, buf.String()
	case "<=":
		return LE, buf.String()
	case ">=":
		return GE, buf.String()
	case "<":
		return LT, buf.String()
	case ">":
		return GT, buf.String()
	case "!":
		// return NOT, buf.String()
		return ILLEGAL, buf.String()
	case "~":
		return TILDE, buf.String()
	case "^":
		return CARET, buf.String()
	case "(":
		return OPEN_PARENTHESIS, buf.String()
	case ")":
		return CLOSE_PARENTHESIS, buf.String()
	case "-":
		return HYPHEN, buf.String()
	}

	// Otherwise return as a regular identifier.
	return ILLEGAL, buf.String()
}

func isVersionUnsafe(ch rune) bool {
	return slices.Contains(versionUnsafe, ch)
}

func isOperatorStart(ch rune) bool { return slices.Contains(operatorsStarts, ch) }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '-' }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// unread places the previously read rune back on the reader.
func (lexer Lexer) unread() {
	err := lexer.r.UnreadRune()
	if err != nil {
		log.Printf("Error during unread rune: %s", err)
	}
}

func (lexer Lexer) read() rune {
	ch, _, err := lexer.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func removeSuperfluousAndOp(tokens []Token, literals []string) ([]Token, []string) {

	tokensToReturn := []Token{}
	literalsToReturn := []string{}

	if slices.Contains(tokens, ILLEGAL) {
		return tokens, literals
	}

	for idx, token := range tokens {

		literal := literals[idx]

		if token == AND && idx >= 2 {

			prevPrevToken := tokens[idx-2]
			prevToken := tokens[idx-1]

			prevPrevPrevToken := EOF
			if idx >= 3 {
				prevPrevPrevToken = tokens[idx-3]
			}

			if IsRangeToken(prevPrevToken) && prevToken == VERSION_EXPRESSION && (len(tokens)-idx) >= 2 && (prevPrevPrevToken == EOF || IsJoinToken(prevPrevPrevToken)) {

				nextToken := tokens[idx+1]
				nextNextToken := tokens[idx+2]

				if IsRangeToken(nextToken) && nextNextToken == VERSION_EXPRESSION {
					continue
				}

			}

		}

		tokensToReturn = append(tokensToReturn, token)
		literalsToReturn = append(literalsToReturn, literal)

	}

	return tokensToReturn, literalsToReturn

}

func augmentMissingEqualityOp(tokens []Token, literals []string) ([]Token, []string) {

	tokensToReturn := []Token{}
	literalsToReturn := []string{}

	if slices.Contains(tokens, ILLEGAL) {
		return tokens, literals
	}

	for idx, token := range tokens {

		literal := literals[idx]

		previousToken := SOF
		if idx >= 1 {
			previousToken = tokens[idx-1]
		}

		nextToken := EOF
		if idx+1 <= len(tokens)-1 {
			nextToken = tokens[idx+1]
		}

		if token == VERSION_EXPRESSION && versions.IsStaticVersion(literal) && (!IsRangeToken(previousToken) && !IsEqualityToken(previousToken)) && previousToken != HYPHEN && nextToken != HYPHEN {
			tokensToReturn = append(tokensToReturn, EQ)
			literalsToReturn = append(literalsToReturn, "=")
		}
		tokensToReturn = append(tokensToReturn, token)
		literalsToReturn = append(literalsToReturn, literal)

	}

	return tokensToReturn, literalsToReturn

}
