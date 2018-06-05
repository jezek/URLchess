// +build js
//go:generate gopherjs build -m
package main

import (
	"URLchess/app"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
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

	location := js.Global.Get("location")
	movesString := ""
	{
		hash := location.Get("hash").String()
		{ // in old times the url was like "host/?moves", now we use host/#moves so fix prviious url
			search := strings.TrimPrefix(location.Get("search").String(), "?")
			if len(search) > 0 && len(hash) == 0 {
				js.Global.Call("alert", "URL is in format for URLchess version <=0.6, transforming...")
				location.Set("href", strings.Replace(location.Get("href").String(), "?", "#", 1))
				return
			}
		}
		movesString = strings.TrimPrefix(hash, "#")
	}

	app, err := app.NewHtml(movesString)
	if err != nil {
		document.Call("write", err.Error())
	}
	app.DrawBoard()

	//js.Global.Call("alert", "calling update board from main")
	if err := app.UpdateBoard(); err != nil {
		document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
		return
	}
}
