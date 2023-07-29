package environment

import (
	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Environment struct {
	values    map[string]any
	Enclosing *Environment
}

func New(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		Enclosing: enclosing,
	}
}

func (env *Environment) Define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) ancestor(distance int) *Environment {
	environment := env
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}

func (env *Environment) GetAt(distance int, lexeme string) any {
	return env.ancestor(distance).values[lexeme]
}

func (env *Environment) Get(name token.Token) any {
	if value, ok := env.values[name.Lexeme]; ok {
		return value
	}

	if env.Enclosing != nil {
		return env.Enclosing.Get(name)
	}

	panic(e.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}

func (env *Environment) AssignAt(distance int, name token.Token, value any) {
	env.ancestor(distance).values[name.Lexeme] = value
}

func (env *Environment) Assign(name token.Token, value any) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	if env.Enclosing != nil {
		env.Enclosing.Assign(name, value)
		return
	}

	panic(e.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}
