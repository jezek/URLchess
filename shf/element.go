package shf

import "github.com/gopherjs/gopherjs/js"

type Element interface {
	Event
	Object() *js.Object
}

var Window Element = &element{js.Global}

type element struct {
	object interface{}
}

func (e *element) Object() *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object)
}

func (e *element) Get(key string) *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object).Get(key)
}

func (e *element) Set(key string, value interface{}) {
	if e == nil || e.object == nil {
		return
	}
	e.object.(*js.Object).Set(key, value)
}

func (e *element) Delete(key string) {
	if e == nil || e.object == nil {
		return
	}
	e.object.(*js.Object).Delete(key)
}

func (e *element) Call(name string, args ...interface{}) *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object).Call(name, args...)
}
