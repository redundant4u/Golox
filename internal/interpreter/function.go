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
	declaration   *ast.Function
	clousure      *env.Environment
	isInitializer bool
}

func NewFunction(declaration *ast.Function, closure *env.Environment, isInitializer bool) *Function {
	return &Function{
		declaration:   declaration,
		clousure:      closure,
		isInitializer: isInitializer,
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (returnValue any) {
	environment := env.New(f.clousure)
	for i, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case Return:
				if f.isInitializer {
					returnValue = f.clousure.GetAt(0, "this")
				} else {
					returnValue = err.Value
				}
			default:
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(f.declaration.Body, environment)
	if f.isInitializer {
		return f.clousure.GetAt(0, "this")
	}
	return nil
}

func (f *Function) Bind(instance *Instance) *Function {
	environment := env.New(f.clousure)
	environment.Define("this", instance)
	return &Function{
		declaration:   f.declaration,
		clousure:      environment,
		isInitializer: f.isInitializer,
	}
}

func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
