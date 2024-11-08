package main

import (
	"URLchess/shf"
	"URLchess/shf/js"
	"os"

	"github.com/andrewbackes/chess/game"
)

const Version = "0.10"

// Apply version to *.html files.
//go:generate sh -c "sed -i 's/?v[0-9.]\\+/?v'\"$(grep '^const Version = ' URLchess.go | grep -oP '\\d+\\.\\d+(\\.\\d+)?')\"'/' *.html"

// Build URLchess and generate assets/URLchess.wasm file.
//go:generate env GOOS=js GOARCH=wasm go build -o assets/URLchess.wasm
// Build URLchess using gopherjs and generate assets/URLchess.js and assets/URLchess.js.map files.
//go:generate env GOOS=js GOARCH=ecmascript gopherjs build -m -o assets/URLchess.js

func main() {
	if js.WRAPS == "" {
		println("Not supposed to be run as file, use 'go generate' and run index.html in browser.")
		os.Exit(-1)
	}

	document := js.Global().Get("document")
	if js.IsUndefined(document) {
		println("ERROR: Undefined document")
		return
	}

	model := &Model{}

	// is rotation supported?
	if div := js.Global().Get("document").Call("createElement", "div"); !js.IsUndefined(div) {
		if !js.IsUndefined(div.Get("style").Get("transform")) {
			model.rotationSupported = true
		}
		div.Call("remove")
	}
	// is execCommand supported?
	if exec := js.Global().Get("document").Get("execCommand"); !js.IsUndefined(exec) {
		model.execSupported = true
	}

	app, err := shf.Create(model)
	if err != nil {
		document.Call("write", "<div id=\"board\" class=\"error\">"+err.Error()+"</div>")
		return
	}

	body := document.Get("body")
	body.Call("appendChild", model.Html.Header.Element.Object())
	body.Call("appendChild", model.Html.Board.Element.Object())
	body.Call("appendChild", model.Html.ThrownOuts.Element.Object())
	body.Call("appendChild", model.Html.Cover.Element.Object())
	body.Call("appendChild", model.Html.Footer.Element.Object())
	body.Call("appendChild", model.Html.Export.Element.Object())
	body.Call("appendChild", model.Html.Notification.Element.Object())

	if hash := js.Global().Get("location").Get("hash").String(); len(hash) > 0 {
		game, err := NewGame(hash)
		if err != nil {
			//TODO - Use app to write error.
			document.Call("write", "<div class=\"error\">"+err.Error()+"</div>")
			return
		}

		model.Game = game
	}

	// If game ended, notify the player.
	if st := model.Game.game.Status(); st != game.InProgress {
		model.showEndGameNotification(app.Tools())
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

	if js.WRAPS == "syscall/js" {
		select {}
	}
}
