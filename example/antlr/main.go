package main

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/lsytj0413/ena/example/antlr/hello"
)

var _ = antlr.NewInputStream

// antlr4 -Dlanguage=Go -o hello -package hello Hello.g4
// antlr4 Hello.g4
// javac *.java
// grun Hello tokens -tokens < test.txt
func main() {
	is := antlr.NewInputStream("hello part")

	lexer := hello.NewHelloLexer(is)
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}

		fmt.Printf("%s (%q)\n", lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}
}
