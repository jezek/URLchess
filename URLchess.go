package main

//go:generate gopherjs build -m

import (
	"URLchess/shf"

	"github.com/andrewbackes/chess/game"
	"github.com/gopherjs/gopherjs/js"
)

const Version = "0.9"

func main() {
	//js.Global.Call("alert", "main")

	document := js.Global.Get("document")
	if document == js.Undefined {
		return
	}

	model := &Model{}

	// is rotation supported?
	if div := js.Global.Get("document").Call("createElement", "div"); div != js.Undefined {
		if div.Get("style").Get("transform") != js.Undefined {
			model.rotationSupported = true
		}
		div.Call("remove")
	}
	// is execCommand supported?
	if exec := js.Global.Get("document").Get("execCommand"); exec != nil && exec != js.Undefined {
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

	if hash := js.Global.Get("location").Get("hash").String(); len(hash) > 0 {
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
}
