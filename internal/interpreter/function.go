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
	declartion *ast.Function
}

func NewFunction(declaration *ast.Function) *Function {
	return &Function{
		declartion: declaration,
	}
}

func (f Function) Arity() int {
	return len(f.declartion.Params)
}

func (f Function) Call(interpreter *Interpreter, arguments []any) any {
	environment := env.New(interpreter.Globals)
	for i, param := range f.declartion.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declartion.Body, environment)
	return nil
}

func (f Function) String() string {
	return "<fn " + f.declartion.Name.Lexeme + ">"
}
