package ast

import (
	"github.com/stretchr/testify/assert"
	"gtihub.com/yudai2929/monkey-lang/token"
	"testing"
)

func TestProgram_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "myVar"}, Value: "myVar"},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "anotherVar"}, Value: "anotherVar"},
			},
		},
	}

	assert.Equal(t, "let myVar = anotherVar;", program.String(), "program.String() wrong")
}
