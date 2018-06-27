// Simple Html Framework
package shf

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Updater interface {
	Update(*Tools) error
}

type Tools struct {
	app     *App
	created []*js.Object
}

func (t *Tools) Update(updaters ...Updater) error {
	for _, updater := range updaters {
		if err := updater.Update(t); err != nil {
			return err
		}
	}
	return nil
}
func (t *Tools) Click(target *js.Object, function func(*js.Object)) error {
	return t.app.Click(target, function)
}
func (t *Tools) CreateElement(etype string) *js.Object {
	elm := js.Global.Get("document").Call("createElement", etype)
	t.created = append(t.created, elm)
	return elm
}
func (t *Tools) Created(elm *js.Object) bool {
	if elm == nil {
		return false
	}

	for _, createdElm := range t.created {
		if elm.Call("isSameNode", createdElm).Bool() {
			return true
		}
	}
	return false
}

func Create(model Updater) (*App, error) {
	app := &App{
		model,
		map[string][]event{},
	}
	if err := app.Update(); err != nil {
		return nil, err
	}
	return app, nil
}

type event struct {
	target   *js.Object
	function func(*js.Object)
}

type App struct {
	model  Updater
	events map[string][]event
}

func (app *App) Update() error {
	if err := app.model.Update(&Tools{app, []*js.Object{}}); err != nil {
		return err
	}
	return nil
}
func (app *App) Click(target *js.Object, function func(*js.Object)) error {
	if app == nil {
		return errors.New("App is nil")
	}
	if target == nil {
		return errors.New("no target provided for click event")
	}
	if app.events == nil {
		app.events = map[string][]event{}
	}

	jsEventName := "click"
	_, ok := app.events[jsEventName]
	if !ok {
		app.events[jsEventName] = []event{}
	}
	for i, e := range app.events[jsEventName] {
		if target.Call("isSameNode", e.target).Bool() {
			e.target.Call("removeEventListener", jsEventName, e.function, false)
			app.events[jsEventName][i].target = nil
			app.events[jsEventName][i].function = nil
			//js.Global.Call("alert", "unregistered event: "+e.target.String()+":"+e.target.Get("id").String())
		}
	}
	//TODO remove nil targets

	if function == nil {
		return nil
	}

	jsEventCallback := func(event *js.Object) {
		function(event)
		if err := app.Update(); err != nil {
			js.Global.Call("alert", "after "+jsEventName+" event app dom update error: "+err.Error())
		}
	}

	target.Call("addEventListener", jsEventName, jsEventCallback, false)
	app.events[jsEventName] = append(app.events[jsEventName], event{target, jsEventCallback})
	//js.Global.Call("alert", "registered event: "+target.String()+":"+target.Get("id").String())
	return nil
}
