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

// Build URLchess using gopherjs and generate assets/URLchess.js and assets/URLchess.js.map files.
//go:generate sh -c "GOPHERJS_GOROOT=\"$(go1.19.13 env GOROOT)\" GOOS=js GOARCH=ecmascript gopherjs build -m -o assets/URLchess.js"

// Build URLchess and generate assets/URLchess.wasm file and copy go wasm runtime to assets.
//go:generate sh -c "GOOS=js GOARCH=wasm go build -o assets/URLchess.wasm"
//go:generate sh -c "cp \"$(go env GOROOT)\"/misc/wasm/wasm_exec.js ./assets"

// Build URLchess and generate assets/URLchess.tinygo.wasm file using tinygo and copy tinygo wasm runtime to assets.
//go:generate sh -c "tinygo build -o assets/URLchess.tinygo.wasm -target wasm ."
//go:generate sh -c "cp \"$(tinygo env TINYGOROOT)\"/targets/wasm_exec.js assets/wasm_exec.tinygo.js"

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

	//TODO jezek - Make it so this is not needed and the board is rotated upon initialization.
	model.RotateBoardForPlayer()

	// If game ended, notify the player.
	if st := model.ChessGame.game.Status(); st != game.InProgress {
		model.showEndGameNotification(app.Tools())
		//TODO jezek - Update only elements needed for showing notification. Or better, make it so the notification is shown upon init ant this is not needed.
	}

	model.Html.Cover.GameStatus.rebuild(app.Tools())

	//TODO jezek - Update only status move body.
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
