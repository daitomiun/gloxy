package main

import (
	"github.com/daitonium/gloxy/ast"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name ast.Token) any {
	if val, ok := e.Values[name.Lexeme]; ok {
		if val == nil {
			panic(RuntimeError{
				token: name,
				msg:   "Variable '" + name.Lexeme + "' is uninitialized.",
			})
		}
		return val
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(RuntimeError{
		token: name,
		msg:   "Undefined variable '" + name.Lexeme + "'.",
	})

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
