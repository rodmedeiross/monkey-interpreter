package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	return obj, ok
}

func (e *Environment) Set(key string, obj Object) Object {
	e.store[key] = obj
	return obj
}

