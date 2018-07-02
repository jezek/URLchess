// +build js
//go:generate gopherjs build -m
package main

import (
	"URLchess/shf"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	//js.Global.Call("alert", "main")

	document := js.Global.Get("document")
	if document == js.Undefined {
		return
	}

	document.Call("write", "<div id=\"header\">URLchess</div>")

	defer func() {
		js.Global.Get("document").Call("write", `<div id="footer">
		<a href="https://jezek.github.io/URLchess">URLchess</a> by jEzEk. Source on <a href="https://github.com/jezek/URLchess">github</a>.
</div>`)
	}()

	model := &Model{}

	// is rotation supported?
	if div := js.Global.Get("document").Call("createElement", "div"); div != js.Undefined {
		if div.Get("style").Get("transform") != js.Undefined {
			model.rotationSupported = true
		}
		div.Call("remove")
	}

	app, err := shf.Create(model)
	if err != nil {
		document.Call("write", "<div id=\"board\" class=\"error\">"+err.Error()+"</div>")
		return
	}

	body := document.Get("body")
	body.Call("appendChild", model.Html.Board.Element.Object())
	body.Call("appendChild", model.Html.ThrownOuts.Element.Object())
	body.Call("appendChild", model.Html.Cover.Element.Object())
	body.Call("appendChild", model.Html.Notification.Element.Object())

	if hash := js.Global.Get("location").Get("hash").String(); len(hash) > 0 {
		game, err := NewGame(hash)
		if err != nil {
			//TODO use app to write error
			document.Call("write", "<div class=\"error\">"+err.Error()+"</div>")
			return
		}

		model.Game = game
	}

	model.RotateBoardForPlayer()
	if err := app.Update(); err != nil {
		if model.Html.Board.Element != nil {
			model.Html.Board.Element.Set("innerHTML", err.Error())
			model.Html.Board.Element.Get("classList").Call("add", "error")
		} else {
			document.Call("write", "<div id=\"board\" class=\"error\">"+err.Error()+"</div>")
		}
		return
	}
}
