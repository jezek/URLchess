// +build js
package main

import (
	"github.com/andrewbackes/chess/game"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	document := js.Global.Get("document")
	document.Call("write", "URLchess<br/>")

	g := game.New()
	document.Call("write", "<pre>"+g.String()+"</pre>")

	document.Call("write", "bye")
}
