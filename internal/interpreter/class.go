package interpreter

type Class struct {
	name    string
	methods map[string]*Function
}

func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{
		name:    name,
		methods: methods,
	}
}

func (c *Class) FindMethod(name string) *Function {
	if _, ok := c.methods[name]; ok {
		return c.methods[name]
	}

	return nil
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewInstance(c)
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (c *Class) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}
