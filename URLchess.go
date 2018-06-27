// +build js
//go:generate gopherjs build -m
package main

import (
	"URLchess/shf"
	"strings"

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
	body.Call("appendChild", model.Html.Board.Object)
	body.Call("appendChild", model.Html.ThrownOuts.Object)
	body.Call("appendChild", model.Html.GameStatus.Object)
	body.Call("appendChild", model.Html.MoveStatus.Object)
	body.Call("appendChild", model.Html.Notification.Object)

	movesString := strings.TrimPrefix(js.Global.Get("location").Get("hash").String(), "#")
	if len(movesString) > 0 {
		game, err := NewGame(movesString)
		if err != nil {
			//TODO use app to write error
			document.Call("write", "<div class=\"error\">"+err.Error()+"</div>")
			return
		}

		model.Game = game
	}

	model.RotateBoardForPlayer()
	if err := app.Update(); err != nil {
		if model.Html.Board.Object != nil {
			model.Html.Board.Object.Set("innerHTML", err.Error())
			model.Html.Board.Object.Get("classList").Call("add", "error")
		} else {
			document.Call("write", "<div id=\"board\" class=\"error\">"+err.Error()+"</div>")
		}
		return
	}
}
