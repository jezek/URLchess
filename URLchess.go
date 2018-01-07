// +build js
package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"github.com/gopherjs/gopherjs/js"
)

var piecesString map[piece.Color]map[piece.Type]string = map[piece.Color]map[piece.Type]string{
	piece.White: map[piece.Type]string{
		piece.King:   "<span class=\"piece white king\">♔</span>",
		piece.Queen:  "<span class=\"piece white queen\">♕</span>",
		piece.Bishop: "<span class=\"piece white bishop\">♗</span>",
		piece.Knight: "<span class=\"piece white knight\">♘</span>",
		piece.Rook:   "<span class=\"piece white rook\">♖</span>",
		piece.Pawn:   "<span class=\"piece white pawn\">♙</span>",
	},
	piece.Black: map[piece.Type]string{
		piece.King:   "<span class=\"piece black king\">♚</span>",
		piece.Queen:  "<span class=\"piece black queen\">♛</span>",
		piece.Bishop: "<span class=\"piece black bishop\">♝</span>",
		piece.Knight: "<span class=\"piece black knight\">♞</span>",
		piece.Rook:   "<span class=\"piece black rook\">♜</span>",
		piece.Pawn:   "<span class=\"piece black pawn\">♟</span>",
	},
}

func main() {
	document := js.Global.Get("document")
	if document == js.Undefined {
		return
	}

	document.Call("write", "<h1>URLchess</h1>")

	defer func() {
		js.Global.Get("document").Call("write", `<div style="margin-top: 2em;border-top: 1px solid black; padding-top:1em; font-size:0.8em;">
	URLchess by jEzEk. Source on <a href="https://github.com/jezek/URLchess">github</a>.
</div>`)
	}()

	location := js.Global.Get("location")
	movesString := ""
	{
		href := location.Get("href").String()
		if i := strings.Index(href, "?"); i != -1 {
			movesString = strings.TrimSpace(href[i+1:])
		}
	}

	moves, err := decodeMoves(movesString) //encoding.go
	if err != nil {
		document.Call("write", "Error decoding moves: "+err.Error())
		return
	}

	g := game.New()

	{
		var err error
		for i, move := range moves {
			if g.Status() != game.InProgress {
				err = errors.New("Too many moves in url string! " + strconv.Itoa(i+1) + " moves are enough")
				break
			}

			_, merr := g.MakeMove(move)
			if merr != nil {
				err = errors.New("Errorneous move number " + strconv.Itoa(i+1) + ": " + merr.Error())
				break
			}
		}

		// is rotation supported?
		rotationSupported := false
		if div := js.Global.Get("document").Call("createElement", "div"); div != js.Undefined {
			if div.Get("style").Get("transform") != js.Undefined {
				rotationSupported = true
			}
			div.Call("remove")
		}

		rotateBoard180deg := rotationSupported && g.ActiveColor() == piece.Black

		{
			//board
			class := ""
			if rotateBoard180deg {
				class = " class=\"rotated180deg\""
			}
			document.Call("write", "<div id=\"board\""+class+">")
		}

		{
			//edging-top
			document.Call("write", "<div id=\"edging-top\">")
			for i := 0; i < 8; i++ {
				document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
			}
			document.Call("write", "</div>")
		}

		if rotationSupported {
			//edging-top-right
			document.Call("write", "<div id=\"edging-top-right\">")
			document.Call("write", "↻")
			document.Call("write", "</div>")

			if etr := js.Global.Get("document").Call("getElementById", "edging-top-right"); etr != js.Undefined {
				etr.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if board := js.Global.Get("document").Call("getElementById", "board"); board != js.Undefined {
							if board.Get("classList").Call("contains", "rotated180deg").Bool() {
								board.Get("classList").Call("remove", "rotated180deg")
							} else {
								board.Get("classList").Call("add", "rotated180deg")
							}
						}
					},
				)
			}
		}

		{
			//edging-left
			document.Call("write", "<div id=\"edging-left\">")
			for i := 8; i > 0; i-- {
				document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
			}
			document.Call("write", "</div>")
		}

		{
			//grid
			document.Call("write", "<div id=\"grid\">")
			squareTones := []string{"light-square", "dark-square"}
			for i := int(63); i >= 0; i-- {
				document.Call("write", "<div id=\""+square.Square(i).String()+"\" class=\""+squareTones[(i%8+i/8)%2]+"\"></div>")
			}
			document.Call("write", "</div>")
		}

		{
			//edging-right
			document.Call("write", "<div id=\"edging-right\">")
			for i := 8; i > 0; i-- {
				document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
			}
			document.Call("write", "</div>")
		}

		if rotationSupported {
			//edging-bottom-left
			document.Call("write", "<div id=\"edging-bottom-left\">")
			document.Call("write", "↻") //↶↷↺↻
			document.Call("write", "</div>")

			if etr := js.Global.Get("document").Call("getElementById", "edging-bottom-left"); etr != js.Undefined {
				etr.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if board := js.Global.Get("document").Call("getElementById", "board"); board != js.Undefined {
							if board.Get("classList").Call("contains", "rotated180deg").Bool() {
								board.Get("classList").Call("remove", "rotated180deg")
							} else {
								board.Get("classList").Call("add", "rotated180deg")
							}
						}
					},
				)
			}

		}

		{
			//edging-bottom
			document.Call("write", "<div id=\"edging-bottom\">")
			for i := 0; i < 8; i++ {
				document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
			}
			document.Call("write", "</div>")
		}

		{
			//board
			document.Call("write", "</div>")
		}

		if e := fillChessBoardGrid(g.Position()); e != nil {
			err = e
		}

		if err != nil {
			document.Call("write", "<div>"+err.Error()+"</div>")
			return
		}
	}

	if g.Status() != game.InProgress {
		document.Call("write", "<div style=\"margin-top: 1em\">Game has ended. "+g.Status().String()+"</div>")
		return
	}
	document.Call("write", "<div id=\"game-status\">")
	document.Call("write", "<p>Player on the move: "+piecesString[g.ActiveColor()][piece.King]+"</p>")
	document.Call("write", "</div>") //game-status

	document.Call("write", `<div id="next-move">
<p>Your next move: <input id="move-input"/> eg. e2e4 (or e7e8q for promotion to queen) and press [ENTER]</p>
<a id="next-move-link" href=""></a><span id="next-move-error"></span>
</div>`)

	moveInput := js.Global.Get("document").Call("getElementById", "move-input")
	if moveInput == js.Undefined {
		document.Call("write", "Next move input element not found")
	} else {
		moveInput.Call(
			"addEventListener",
			"keyup",
			func(event *js.Object) {
				if keycode := event.Get("keyCode").Int(); keycode == 13 {
					nextMoveError := js.Global.Get("document").Call("getElementById", "next-move-error")
					if nextMoveError == js.Undefined {
						document.Call("write", "Next move error element not found")
						return
					}
					nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link")
					if nextMoveLink == js.Undefined {
						nextMoveError.Set("innerHTML", "Next move link element not found")
						return
					}

					nextMoveError.Set("innerHTML", "")
					nextMoveLink.Set("innerHTML", "")
					nextMoveLink.Set("href", "")

					nextMovePCM := strings.TrimSpace(moveInput.Get("value").String())
					if nextMovePCM == "" {
						return
					}

					nextMove := move.Parse(nextMovePCM)
					if nextMove == move.Null {
						nextMoveError.Set("innerHTML", "Next move is not in PCN format")
						return
					}

					if _, ok := g.LegalMoves()[nextMove]; ok == false {
						nextMoveError.Set("innerHTML", "Next move is not a legal move")
						return
					}

					nextMoveString, err := encodeMove(nextMove) //encoding.go
					if err != nil {
						nextMoveError.Set("innerHTML", "Next move encoding error: "+err.Error())
						return
					}

					url := location.Get("origin").String() + location.Get("pathname").String() + "?" + movesString + nextMoveString
					nextMoveLink.Set("innerHTML", url)
					nextMoveLink.Set("href", url)
					nextMoveError.Set("innerHTML", " <- copy this link and send to your oponent")
				}
			},
			false,
		)
		moveInput.Call("focus")
	}
}

func fillChessBoardGrid(position *position.Position) error {
	//fill board grid with markers and pieces
	for i := int(63); i >= 0; i-- {
		sq := square.Square(i)
		sqElm := js.Global.Get("document").Call("getElementById", sq.String())
		if sqElm == js.Undefined {
			return errors.New("Can't find square element: " + sq.String())
		}
		innerSquare := ""
		lm := position.LastMove
		if lm != move.Null && (lm.Source == sq || lm.Destination == sq) {
			//last-move from or to marker is on square
			dir := "from"
			if lm.Destination == sq {
				dir = "to"
			}
			innerSquare += "<span class=\"marker last-move " + dir + "\"></span>"
		}
		pc := position.OnSquare(sq)
		if pieceString, ok := piecesString[pc.Color][pc.Type]; ok {
			innerSquare += pieceString
		}

		sqElm.Set("innerHTML", innerSquare)
	}
	return nil
}
