package lexer

type TokenType string

type LangToken struct {
	Type    TokenType // LangToken type
	Literal string    // LangToken literal
	Line    int       // 0-based line number where the token starts
	Pos     int       // 1-based column number where the token starts
	Length  int       // Length of the token literal in runes
}

const (
	TokenTypeUndefined        TokenType = "Undefined"
	TokenTypeEOF              TokenType = "EOF"
	TokenTypeIdentifier       TokenType = "Identifier"
	TokenTypeNumber           TokenType = "Number"
	TokenTypeString           TokenType = "String"
	TokenTypeAssignment       TokenType = "Assignment"
	TokenTypePlus             TokenType = "Plus"
	TokenTypeMinus            TokenType = "Minus"
	TokenTypeMultiply         TokenType = "Multiply"
	TokenTypeDot              TokenType = "Dot"
	TokenTypeDivide           TokenType = "Divide"
	TokenTypeEqual            TokenType = "Equal"
	TokenTypeLessThan         TokenType = "LessThan"
	TokenTypeLessThanEqual    TokenType = "LessThanEqual"
	TokenTypeGreaterThan      TokenType = "GreaterThan"
	TokenTypeLeftParenthesis  TokenType = "LeftParenthesis"
	TokenTypeRightParenthesis TokenType = "RightParenthesis"
	TokenTypeLeftBrace        TokenType = "LeftBrace"
	TokenTypeRightBrace       TokenType = "RightBrace"
	TokenTypeLeftBracket      TokenType = "LeftBracket"
	TokenTypeRightBracket     TokenType = "RightBracket"
	TokenTypeComma            TokenType = "Comma"
	TokenTypeSemicolon        TokenType = "Semicolon"
	TokenTypeColon            TokenType = "Colon"
	TokenTypeQuestionMark     TokenType = "QuestionMark"
	TokenTypeArrow            TokenType = "Arrow"
	TokenTypeLambdaArrow      TokenType = "LambdaArrow"
	TokenTypeIf               TokenType = "If"
	TokenTypeThen             TokenType = "Then"
	TokenTypeElse             TokenType = "Else"
	TokenTypeFor              TokenType = "For"
	TokenTypeWhile            TokenType = "While"
	TokenTypeIn               TokenType = "In"
	TokenTypeRange            TokenType = "Range"
	TokenTypeDo               TokenType = "Do"
	TokenTypeSwitch           TokenType = "Switch"
	TokenTypeCase             TokenType = "Case"
	TokenTypeDefault          TokenType = "Default"
	TokenTypeData             TokenType = "Data"
	TokenTypeType             TokenType = "Type"
	TokenTypeAssembly         TokenType = "Assembly"
	TokenTypeMain             TokenType = "Main"
	TokenTypeComment          TokenType = "Comment"
	TokenTypeMultiLineComment TokenType = "MultiLineComment"
	TokenTypeReturn           TokenType = "Return"
	TokenTypeSyscall          TokenType = "Syscall"
	TokenTypeImport           TokenType = "Import"
)

const TokenTypeFunction TokenType = "Function"

const TokenTypeLet TokenType = "Let"

// Keywords is a map of reserved keywords to their corresponding token types.
var Keywords = map[string]TokenType{
	"function": TokenTypeFunction,
	"let":      TokenTypeLet,
	"if":       TokenTypeIf,
	"in":       TokenTypeIn,
	"range":    TokenTypeRange,
	"->":       TokenTypeLambdaArrow,
	"else":     TokenTypeElse,
	"then":     TokenTypeThen,
	"for":      TokenTypeFor,
	"while":    TokenTypeWhile,
	"do":       TokenTypeDo,
	"switch":   TokenTypeSwitch,
	"case":     TokenTypeCase,
	"default":  TokenTypeDefault,
	"data":     TokenTypeData,
	"type":     TokenTypeType,
	"return":   TokenTypeReturn,
	"asm":      TokenTypeAssembly,
	"syscall":  TokenTypeSyscall,
	"import":   TokenTypeImport,
	// Add more keywords here
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := Keywords[ident]; ok {
		return tokType
	}
	return TokenTypeIdentifier
}

func newToken(tokenType TokenType, ch rune) LangToken {
	return LangToken{Type: tokenType, Literal: string(ch)}
}
