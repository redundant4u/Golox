package interpreter

import (
	"github.com/redundant4u/Golox/internal/ast"
	env "github.com/redundant4u/Golox/internal/environment"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
}

type Function struct {
	declaration *ast.Function
	clousure    *env.Environment
}

func NewFunction(declaration *ast.Function, closure *env.Environment) *Function {
	return &Function{
		declaration: declaration,
		clousure:    closure,
	}
}

func (f Function) Arity() int {
	return len(f.declaration.Params)
}

func (f Function) Call(interpreter *Interpreter, arguments []any) (returnValue any) {
	environment := env.New(f.clousure)
	for i, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(Return); ok {
				returnValue = err.Value
			} else {
				panic(r)
			}
		}
	}()
	interpreter.executeBlock(f.declaration.Body, environment)
	returnValue = nil

	return
}

func (f Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
