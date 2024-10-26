package constraints

type Token string

const (
	// Special tokens
	ILLEGAL Token = "ILLEGAL"
	EOF     Token = "EOF"
	SOF     Token = "SOF"
	WS      Token = "WS"

	// Keywords
	EQ                 Token = "EQ"                // =
	LT                 Token = "LT"                // <
	LE                 Token = "LE"                // <=
	GT                 Token = "GT"                // >
	GE                 Token = "GE"                // >=
	NOT                Token = "NOT"               // !
	TILDE              Token = "TILDE"             // ~
	CARET              Token = "CARET"             // ^
	OPEN_PARENTHESIS   Token = "OPEN_PARENTHESIS"  // (
	CLOSE_PARENTHESIS  Token = "CLOSE_PARENTHESIS" // )
	AND                Token = "AND"               // &&
	OR                 Token = "OR"                // ||
	HYPHEN             Token = "HYPHEN"            // -
	STAR               Token = "STAR"              // *
	ANY                Token = "ANY"               // ANY
	WILDCARD_X         Token = "WILDCARD_X"        // X
	CONSTRAINT         Token = "CONSTRAINT"
	UNKNOW_IDENTIFIER  Token = "UNKNOW_IDENTIFIER"
	VERSION_EXPRESSION Token = "VERSION_EXPRESSION"
)

func IsOperatorToken(token Token) bool {
	return token == EQ || token == LT || token == LE || token == GT || token == GE || token == TILDE || token == CARET || token == OPEN_PARENTHESIS || token == CLOSE_PARENTHESIS || token == AND || token == OR || token == HYPHEN
}

func IsRangeToken(token Token) bool {
	return token == LT || token == LE || token == GT || token == GE || token == TILDE || token == CARET || token == HYPHEN
}

func IsEqualityToken(token Token) bool {
	return token == EQ
}

func IsJoinToken(token Token) bool {
	return token == OR || token == AND
}

func (token Token) toString() string {
	switch token {
	case EQ:
		return "="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	case TILDE:
		return "~"
	case CARET:
		return "^"
	case AND:
		return "&&"
	case OR:
		return "||"
	case HYPHEN:
		return "-"
	default:
		return ""
	}
}
