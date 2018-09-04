package meander

type Facade interface {
	Public() interface{}
}

func Public(o interface{}) interface{} {
	// check if the argument has facade interface
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
