package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewWrappedEnvironment(env *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: env,
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}
	return obj, ok
}

func (e *Environment) Set(key string, obj Object) Object {
	e.store[key] = obj
	return obj
}

