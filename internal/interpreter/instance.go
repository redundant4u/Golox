package interpreter

import (
	"github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (i *Instance) Get(name token.Token) any {
	if field, ok := i.fields[name.Lexeme]; ok {
		return field
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i)
	}

	panicMsg := "Undefined property '" + name.Lexeme + "'."
	error.ReportRuntimeError(name, panicMsg)
	panic(panicMsg)
}

func (i *Instance) Set(name token.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}
