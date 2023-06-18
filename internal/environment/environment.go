package environment

import (
	"github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func New(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) GetAt(distance int, lexeme string) any {
	return e.ancestor(distance).values[lexeme]
}

func (e *Environment) Get(name token.Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	panic(error.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) AssignAt(distance int, name token.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) Assign(name token.Token, value any) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	panic(error.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}
