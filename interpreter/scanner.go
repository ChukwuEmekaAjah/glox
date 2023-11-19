package interpreter

// Scanner represents a source code scanner
type Scanner struct {
	source string
}

// Token represents a token after lexing the source code
type Token struct {
}

// NewScanner creates a new Scanner instance
func NewScanner(sourceCode string) Scanner {
	return Scanner{source: sourceCode}
}

// ScanTokens scans the tokens
func (scanner *Scanner) ScanTokens() []Token {
	return []Token{}
}
