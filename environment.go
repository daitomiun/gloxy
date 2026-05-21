package main

import (
	"fmt"

	"github.com/daitonium/gloxy/ast"
)

type Environment struct {
	Values map[string]any
}

func (e *Environment) Define(name string, value any) {
	fmt.Println("define value")
	e.Values[name] = value
	fmt.Printf("deine values -> %s \n", e.Values)
}

func (e *Environment) Get(name ast.Token) any {
	fmt.Printf("lexeme -> %s \n", name.Lexeme)
	val, ok := e.Values[name.Lexeme]
	fmt.Printf("values -> %s \n", e.Values)
	fmt.Printf("val -> %s \n", val)
	if ok {
		fmt.Println("val found")
		return val
	}
	panic(RuntimeError{
		token: name,
		msg:   "Undefined variable '" + name.Lexeme + "'.",
	})
}
