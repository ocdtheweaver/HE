package lexer

import (
"testing"
)

func TestLexerBasic(t *testing.T) {
input := `module ui {
    fn navbar(items) [
        print "Hello"
    ]
}`

lex := New(input)

expectedTypes := []TokenType{
MODULE, IDENT, LBRACE,
FN, IDENT, LPAREN, IDENT, RPAREN, LBRACK,
PRINT, STRING, RBRACK,
RBRACE, EOF,
}

i := 0
for {
tok := lex.Next()
if tok.Type == EOF {
break
}

if i >= len(expectedTypes) {
t.Errorf("More tokens than expected")
break
}

if tok.Type != expectedTypes[i] {
t.Errorf("Token %d: expected %v, got %v (lit: %q)", i, expectedTypes[i], tok.Type, tok.Lit)
}
i++
}
}

func TestLexerComments(t *testing.T) {
input := `~ This is a comment ~
print "Hello" ~ inline comment ~`

lex := New(input)

tokens := []Token{}
for {
tok := lex.Next()
if tok.Type == EOF {
break
}
tokens = append(tokens, tok)
}

// Should only have print and string tokens, no comment tokens
if len(tokens) != 2 {
t.Errorf("Expected 2 tokens after skipping comments, got %d", len(tokens))
}

if tokens[0].Type != PRINT {
t.Errorf("First token should be PRINT, got %v", tokens[0].Type)
}

if tokens[1].Type != STRING {
t.Errorf("Second token should be STRING, got %v", tokens[1].Type)
}
}
