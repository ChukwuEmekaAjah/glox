package interpreter

// Scanner represents a source code scanner
type Scanner struct {
	source  string
	tokens  []Token
	current uint
	start   uint
	line    uint
}

// NewScanner creates a new Scanner instance
func NewScanner(sourceCode string) Scanner {
	return Scanner{source: sourceCode, tokens: make([]Token, 0), current: 0, start: 0, line: 1}
}

// ScanTokens scans the tokens
func (scanner *Scanner) ScanTokens() []Token {

	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, Token{tokenType: "EOF", lexeme: "", literal: nil, line: scanner.line})
	return scanner.tokens
}

// isAtEnd tells us if we've reached the end of the source code string
func (scanner Scanner) isAtEnd() bool {
	return scanner.current >= uint(len(scanner.source))
}

// scanToken looks for lexemes inside the source string
func (scanner *Scanner) scanToken() {
	c := scanner.advance()

	switch c {
	case '(':
		scanner.addToken("LEFT_PAREN", nil)
	case ')':
		scanner.addToken("RIGHT_PAREN", nil)
	case '{':
		scanner.addToken("LEFT_BRACE", nil)
	case '}':
		scanner.addToken("RIGHT_BRACE", nil)
	case ',':
		scanner.addToken("COMMA", nil)
	case '.':
		scanner.addToken("DOT", nil)
	case '-':
		scanner.addToken("MINUS", nil)
	case '+':
		scanner.addToken("PLUS", nil)
	case ';':
		scanner.addToken("SEMICOLON", nil)
	case '*':
		scanner.addToken("STAR", nil)
	case '!':
		if scanner.match('=') {
			scanner.addToken("BANG_EQUAL", nil)
		} else {
			scanner.addToken("BANG", nil)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken("EQUAL_EQUAL", nil)
		} else {
			scanner.addToken("EQUAL", nil)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken("LESS_EQUAL", nil)
		} else {
			scanner.addToken("LESS", nil)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken("GREATER_EQUAL", nil)
		} else {
			scanner.addToken("GREATER", nil)
		}
	case '/':
		if scanner.match('/') {
			// A comment goes until the end of the line
			for scanner.peek(0) != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else if scanner.match('*') {
			// A multiline comment is here and goes until the next */
			for scanner.peek(0) != '*' && scanner.peek(1) != '/' && !scanner.isAtEnd() {
				scanner.advance()
			}
			scanner.advance()
			scanner.advance()
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		scanner.line++
	default:
		ReportError(scanner.line, "Unexpected character.")
	}
}

// advance returns the next character in the source code
func (scanner *Scanner) advance() byte {
	scanner.current++
	return scanner.source[scanner.current-1]
}

func (scanner *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, Token{tokenType: tokenType, lexeme: text, literal: literal, line: scanner.line})
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.source[scanner.current] != byte(expected) {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) peek(position uint) byte {
	if scanner.isAtEnd() {
		return byte(0)
	}

	return scanner.source[scanner.current+position]
}
