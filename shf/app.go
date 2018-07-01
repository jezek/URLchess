// Simple Html Framework
package shf

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Element struct {
	object interface{}
}

var Window = &Element{js.Global}

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
func (t *Tools) HashChange(function func(*Event) error) error {
	return t.app.HashChange(function)
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
func (app *App) HashChange(function func(*Event) error) error {
	return app.elventListener("hashchange", Window, function)
}
func (app *App) Click(target *Element, function func(*Event) error) error {
	return app.elventListener("click", target, function)
}
func (app *App) elventListener(eventName string, target *Element, function func(*Event) error) error {
	if app == nil {
		return errors.New("App is nil")
	}
	if eventName == "" {
		return errors.New("no event name")
	}
	if target == nil {
		return errors.New("no target")
	}

	if app.events == nil {
		app.events = map[string]map[*Element]func(*js.Object){}
	}

	_, ok := app.events[eventName]
	if !ok {
		app.events[eventName] = map[*Element]func(*js.Object){}
	}

	if registeredFunc, ok := app.events[eventName][target]; ok {
		target.Call("removeEventListener", eventName, registeredFunc, false)
		delete(app.events[eventName], target)
		//js.Global.Call("alert", "unregistered event: "+e.target.String()+":"+e.target.Get("id").String())
	}

	if function == nil {
		return nil
	}

	jsEventCallback := func(event *js.Object) {
		//TODO handle errors via app.ErrorCallback
		if err := function(&Event{event}); err != nil {
			js.Global.Call("alert", eventName+" event function returned error: "+err.Error())
			return
		}
		if err := app.Update(); err != nil {
			js.Global.Call("alert", "after "+eventName+" event app dom update error: "+err.Error())
			return
		}
	}

	target.Call("addEventListener", eventName, jsEventCallback, false)
	app.events[eventName][target] = jsEventCallback
	//js.Global.Call("alert", "registered event: "+target.String()+":"+target.Get("id").String())
	return nil
}
