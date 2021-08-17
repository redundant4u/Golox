package lox

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) interface{} {
	if val, ok := e.values[name.Lexeme]; ok == true {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	errmsg := "Undefined variable."
	RuntimeError(name, errmsg)
	panic(errmsg)
}

func (e *Environment) Assign(name Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok == true {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	errmsg := "Undefined variable."
	RuntimeError(name, errmsg)
	panic(errmsg)
}
