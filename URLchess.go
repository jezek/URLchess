// +build js
//go:generate gopherjs build -m
package main

import (
	"URLchess/app"
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

	htmlApp := app.HtmlApp{}
	htmlApp.SetRootElement(document.Get("body"))

	if err := htmlApp.UpdateDom(); err != nil {
		document.Call("write", err.Error())
	}

	movesString := strings.TrimPrefix(js.Global.Get("location").Get("hash").String(), "#")
	if len(movesString) > 0 {
		//js.Global.Call("alert", movesString)
		game, err := app.NewGame(movesString)
		if err != nil {
			document.Call("write", err.Error())
		}

		htmlApp.Game = game
	}

	htmlApp.RotateBoardForPlayer()

	if err := htmlApp.UpdateDom(); err != nil {
		document.Call("write", err.Error())
	}
}
