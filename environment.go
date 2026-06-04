package main

import (
	"fmt"

	"github.com/daitonium/gloxy/ast"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
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
	if !ok {
		panic(RuntimeError{
			token: name,
			msg:   "Undefined variable '" + name.Lexeme + "'.",
		})
	}
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}
	fmt.Println("val found")
	return val

}

func (e *Environment) Assign(name ast.Token, value any) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	}
	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}
	panic(RuntimeError{
		token: name,
		msg:   "Undefined variable '" + name.Lexeme + "'.",
	})

}
