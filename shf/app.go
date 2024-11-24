// Simple Html Framework
package shf

import (
	"URLchess/shf/js"
	"errors"
	"strconv"
	"time"
)

type Initializer interface {
	Init(*Tools) error
}
type Updater interface {
	Update(*Tools) error
}
type Destroyer interface {
	Destroy(*Tools)
}

type Tools struct {
	app *App
}

func (t *Tools) Initialize(updaters ...Updater) error {
	for _, updater := range updaters {
		if updater == nil {
			continue
		}
		if initializer, ok := updater.(Initializer); ok {
			if err := initializer.Init(t); err != nil {
				return err
			}
		}
	}
	return nil
}
func (t *Tools) Update(updaters ...Updater) error {
	for _, updater := range updaters {
		if updater == nil {
			continue
		}
		if err := updater.Update(t); err != nil {
			return err
		}
	}
	return nil
}
func (t *Tools) Destroy(elements ...Element) {
	for _, element := range elements {
		if element == nil {
			continue
		}
		if destroyer, ok := element.(Destroyer); ok {
			destroyer.Destroy(t)
		}
		t.app.DestroyElement(element)
	}
}

func (t *Tools) Input(target Element, function func(e Event) error) error {
	return t.app.Input(target, function)
}
func (t *Tools) Click(target Element, function func(e Event) error) error {
	return t.app.Click(target, function)
}
func (t *Tools) DblClick(target Element, function func(e Event) error) error {
	return t.app.DblClick(target, function)
}
func (t *Tools) ClickRemove(target Element) error {
	return t.app.Click(target, nil)
}
func (t *Tools) DblClickRemove(target Element) error {
	return t.app.DblClick(target, nil)
}
func (t *Tools) HashChange(function func(HashChangeEvent) error) error {
	return t.app.HashChange(function)
}
func (t *Tools) CreateElement(etype string) Element {
	return t.app.CreateElement(etype)
}
func (t *Tools) CreateTextNode(text string) js.Object {
	return js.Global().Get("document").Call("createTextNode", text)
}
func (t *Tools) Destroylement(elm Element) {
	t.app.DestroyElement(elm)
}
func (t *Tools) Created(elm Element) bool {
	return t.app.ElementCreated(elm)
}
func (t *Tools) Timer(duration time.Duration, callback func()) int {
	return t.app.Timer(duration, callback)
}

func Create(model Updater) (*App, error) {
	app := &App{
		model,
		nil,
		nil,
		map[Element]struct{}{},
	}
	app.tools = &Tools{app}

	if err := app.Initialize(); err != nil {
		return nil, err
	}
	if err := app.Update(); err != nil {
		return nil, err
	}
	return app, nil
}

type App struct {
	model   Updater
	tools   *Tools
	events  map[string]map[Element]js.Func
	created map[Element]struct{}
}

func (app *App) Tools() *Tools { return app.tools }
func (app *App) Initialize() error {
	if err := app.tools.Initialize(app.model); err != nil {
		return err
	}
	return nil
}
func (app *App) Update() error {
	if err := app.tools.Update(app.model); err != nil {
		return err
	}
	return nil
}
func (app *App) CreateElement(etype string) Element {
	elm := &element{CreateElementObject(etype)}
	if _, ok := app.created[elm]; ok {
		js.Global().Call("alert", "*App.CreateElement: an element can not be created twice the same. Why is this happening?")
		return nil
	}
	app.created[elm] = struct{}{}
	return elm
}
func (app *App) ElementCreated(elm Element) bool {
	_, ok := app.created[elm]
	return ok
}
func (app *App) DestroyElement(elm Element) {
	if !app.ElementCreated(elm) {
		return
	}
	delete(app.created, elm)
	DestroyElementObject(elm.Object())
}
func (app *App) HashChange(function func(HashChangeEvent) error) error {
	return app.elventListener("hashchange", Window, func(e Event) error {
		hce := &hashChangeEvent{e}
		if err := function(hce); err != nil {
			return err
		}
		return nil
	})
}
func (app *App) Input(target Element, function func(e Event) error) error {
	return app.elventListener("input", target, function)
}
func (app *App) Click(target Element, function func(e Event) error) error {
	return app.elventListener("click", target, function)
}
func (app *App) DblClick(target Element, function func(e Event) error) error {
	return app.elventListener("dblclick", target, function)
}
func (app *App) elventListener(eventName string, target Element, function func(e Event) error) error {
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
		app.events = map[string]map[Element]js.Func{}
	}

	_, ok := app.events[eventName]
	if !ok {
		app.events[eventName] = map[Element]js.Func{}
	}

	if registeredFunc, ok := app.events[eventName][target]; ok {
		target.Call("removeEventListener", eventName, registeredFunc, false)
		delete(app.events[eventName], target)
		registeredFunc.Release()
		//js.Global().Call("alert", "unregistered event: "+e.target.String()+":"+e.target.Get("id").String())
	}

	if function == nil {
		return nil
	}

	jsEventCallback := js.FuncOf(func(_ js.Object, args []js.Object) any {
		//TODO handle errors via app.ErrorCallback
		if len(args) < 1 {
			js.Global().Call("alert", eventName+" event function called with no arguments")
			return errors.New(eventName + " event function called with no arguments")
		}
		if err := function(args[0]); err != nil {
			js.Global().Call("alert", eventName+" event function returned error: "+err.Error())
			return err
		}
		if err := app.Update(); err != nil {
			js.Global().Call("alert", "after "+eventName+" event app dom update error: "+err.Error())
			return err
		}
		return nil
	})

	target.Call("addEventListener", eventName, jsEventCallback, false)
	app.events[eventName][target] = jsEventCallback
	//js.Global().Call("alert", "registered event: "+target.String()+":"+target.Get("id").String())
	return nil
}

func CreateElementObject(etype string) js.Object {
	return js.Global().Get("document").Call("createElement", etype)
}
func DestroyElementObject(o js.Object) {
	o.Call("remove")
}
func (app *App) Timer(duration time.Duration, callback func()) int {
	ms := int(duration / time.Millisecond) // in milliseconds
	timeoutId := 0
	timeoutId = js.Global().Call("setTimeout", js.FuncOf(func(this js.Object, args []js.Object) any {
		callback()
		//js.Global().Call("alert", "timer "+strconv.Itoa(timeoutId)+" gone off after: "+strconv.Itoa(ms))
		if err := app.Update(); err != nil {
			js.Global().Call("alert", "after timer "+strconv.Itoa(timeoutId)+" app dom update error: "+err.Error())
			return err
		}
		return nil
	}), ms).Int()
	//js.Global().Call("alert", "timer "+strconv.Itoa(timeoutId)+" set to: "+strconv.Itoa(ms))
	return timeoutId
}
