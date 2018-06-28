// Simple Html Framework
package shf

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Element struct {
	object interface{}
}

func (e *Element) Object() *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object)
}

func (e *Element) Get(key string) *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object).Get(key)
}

func (e *Element) Set(key string, value interface{}) {
	if e == nil || e.object == nil {
		return
	}
	e.object.(*js.Object).Set(key, value)
}

func (e *Element) Delete(key string) {
	if e == nil || e.object == nil {
		return
	}
	e.object.(*js.Object).Delete(key)
}

func (e *Element) Call(name string, args ...interface{}) *js.Object {
	if e == nil || e.object == nil {
		return js.Undefined
	}
	return e.object.(*js.Object).Call(name, args...)
}

type Event struct {
	*js.Object
}

type Initializer interface {
	Init(*Tools) error
}
type Updater interface {
	Update(*Tools) error
}

type Tools struct {
	app     *App
	created map[*Element]bool
}

func (t *Tools) Update(updaters ...Updater) error {
	for _, updater := range updaters {
		if updater == nil {
			continue
		}
		if initializer, ok := updater.(Initializer); ok {
			if !t.app.initialized[initializer] {
				if err := initializer.Init(t); err != nil {
					return err
				}
				t.app.initialized[initializer] = true
			}
		}
		if err := updater.Update(t); err != nil {
			return err
		}
	}
	return nil
}
func (t *Tools) Click(target *Element, function func(*Event) error) error {
	return t.app.Click(target, function)
}

func (t *Tools) CreateElement(etype string) *Element {
	elm := &Element{js.Global.Get("document").Call("createElement", etype)}
	if t.created[elm] {
		js.Global.Call("alert", "Tools.ElementCreate: an element can not be created twice the same. Why is this happening?")
		return nil
	}
	t.created[elm] = true
	return elm
}
func (t *Tools) Created(elm *Element) bool {
	if elm == nil {
		return false
	}

	return t.created[elm]
}

func Create(model Updater) (*App, error) {
	app := &App{
		model,
		nil,
		map[Initializer]bool{},
	}
	if err := app.Update(); err != nil {
		return nil, err
	}
	return app, nil
}

type App struct {
	model       Updater
	events      map[string]map[*Element]func(*js.Object)
	initialized map[Initializer]bool
}

func (app *App) Update() error {
	tools := &Tools{app, map[*Element]bool{}}
	if err := tools.Update(app.model); err != nil {
		return err
	}
	return nil
}
func (app *App) Click(target *Element, function func(*Event) error) error {
	if app == nil {
		return errors.New("App is nil")
	}
	if target == nil {
		return errors.New("no target provided for click event")
	}
	if app.events == nil {
		app.events = map[string]map[*Element]func(*js.Object){}
	}

	jsEventName := "click"
	_, ok := app.events[jsEventName]
	if !ok {
		app.events[jsEventName] = map[*Element]func(*js.Object){}
	}

	if registeredFunc, ok := app.events[jsEventName][target]; ok {
		target.Call("removeEventListener", jsEventName, registeredFunc, false)
		delete(app.events[jsEventName], target)
		//js.Global.Call("alert", "unregistered event: "+e.target.String()+":"+e.target.Get("id").String())
	}

	if function == nil {
		return nil
	}

	jsEventCallback := func(event *js.Object) {
		//TODO handle errors via app.ErrorCallback
		if err := function(&Event{event}); err != nil {
			js.Global.Call("alert", jsEventName+" event function returned error: "+err.Error())
			return
		}
		if err := app.Update(); err != nil {
			js.Global.Call("alert", "after "+jsEventName+" event app dom update error: "+err.Error())
			return
		}
	}

	target.Call("addEventListener", jsEventName, jsEventCallback, false)
	app.events[jsEventName][target] = jsEventCallback
	//js.Global.Call("alert", "registered event: "+target.String()+":"+target.Get("id").String())
	return nil
}
