package main

import "github.com/daitonium/gloxy/ast"

type Environment struct {
	Values map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name ast.Token) any {
	if _, ok := e.Values[name.Lexeme]; ok {
		return e.Values[name.Lexeme]
	}
	panic(RuntimeError{
		token: name,
		msg:   "Undefined variable '" + name.Lexeme + "'.",
	})
}
