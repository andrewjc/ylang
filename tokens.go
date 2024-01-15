package compiler

type TokenType string

type LangToken struct {
	Type    TokenType // LangToken type
	Literal string    // LangToken literal
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
	TokenTypeDivide           TokenType = "Divide"
	TokenTypeLessThan         TokenType = "LessThan"
	TokenTypeLessThanEqual    TokenType = "LessThanEqual"
	TokenTypeGreaterThan      TokenType = "GreaterThan"
	TokenTypeLeftParenthesis  TokenType = "LeftParenthesis"
	TokenTypeRightParenthesis TokenType = "RightParenthesis"
	TokenTypeLeftBrace        TokenType = "LeftBrace"
	TokenTypeRightBrace       TokenType = "RightBrace"
	TokenTypeComma            TokenType = "Comma"
	TokenTypeSemicolon        TokenType = "Semicolon"
	TokenTypeColon            TokenType = "Colon"
	TokenTypeQuestionMark     TokenType = "QuestionMark"
	TokenTypeArrow            TokenType = "Arrow"
	TokenTypeLambdaArrow      TokenType = "LambdaArrow"
	TokenTypeIf               TokenType = "If"
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
	"for":      TokenTypeFor,
	"while":    TokenTypeWhile,
	"do":       TokenTypeDo,
	"switch":   TokenTypeSwitch,
	"case":     TokenTypeCase,
	"default":  TokenTypeDefault,
	"data":     TokenTypeData,
	"type":     TokenTypeType,
	"asm":      TokenTypeAssembly,
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
