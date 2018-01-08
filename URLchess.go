// +build js
package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
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

	moves, err := decodeMoves(movesString) // encoding.go
	if err != nil {
		document.Call("write", "Error decoding moves: "+err.Error())
		return
	}

	g := game.New()

	{
		// apply game moves
		for i, move := range moves {
			if g.Status() != game.InProgress {
				document.Call("write", "Too many moves in url string! "+strconv.Itoa(i+1)+" moves are enough")
				return
			}

			_, merr := g.MakeMove(move)
			if merr != nil {
				document.Call("write", "Errorneous move number "+strconv.Itoa(i+1)+": "+merr.Error())
				return
			}
		}
	}

	{
		// draw chess board

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
			// board
			class := ""
			if rotateBoard180deg {
				class = " class=\"rotated180deg\""
			}
			document.Call("write", "<div id=\"board\""+class+">")
		}

		{
			// edging-top
			document.Call("write", "<div id=\"edging-top\">")
			for i := 0; i < 8; i++ {
				document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
			}
			document.Call("write", "</div>")
		}

		if rotationSupported {
			// edging-top-right
			document.Call("write", "<div id=\"edging-top-right\">")
			document.Call("write", "↻")
			document.Call("write", "</div>")

			if etr := js.Global.Get("document").Call("getElementById", "edging-top-right"); etr != nil {
				etr.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if board := js.Global.Get("document").Call("getElementById", "board"); board != nil {
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
			// edging-left
			document.Call("write", "<div id=\"edging-left\">")
			for i := 8; i > 0; i-- {
				document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
			}
			document.Call("write", "</div>")
		}

		{
			// grid
			document.Call("write", "<div class=\"grid\">")
			squareTones := []string{"light-square", "dark-square"}
			for i := int(63); i >= 0; i-- {
				document.Call("write", "<div id=\""+square.Square(i).String()+"\" class=\""+squareTones[(i%8+i/8)%2]+"\"></div>")
			}
			document.Call("write", "</div>")
		}

		{
			// edging-right
			document.Call("write", "<div id=\"edging-right\">")
			for i := 8; i > 0; i-- {
				document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
			}
			document.Call("write", "</div>")
		}

		if rotationSupported {
			// edging-bottom-left
			document.Call("write", "<div id=\"edging-bottom-left\">")
			document.Call("write", "↻") // ↶↷↺↻
			document.Call("write", "</div>")

			if etr := js.Global.Get("document").Call("getElementById", "edging-bottom-left"); etr != nil {
				etr.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if board := js.Global.Get("document").Call("getElementById", "board"); board != nil {
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
			// edging-bottom
			document.Call("write", "<div id=\"edging-bottom\">")
			for i := 0; i < 8; i++ {
				document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
			}
			document.Call("write", "</div>")
		}

		{
			// board
			document.Call("write", "</div>")
		}
		document.Call("write", `<div id="game-status">
		<p>Moving player: <span id="moving-player"></span></p>
		<p>Game status: <span id="game-progress"><span></p>
</div>`)

		if err := fillChessBoard(g, move.Null); err != nil {
			document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
			return
		}
	}

	document.Call("write", `<div id="next-move">
<p>Your next move: <input id="move-input"/> eg. e2e4 (or e7e8q for promotion to queen) and press [ENTER]</p>
<a id="next-move-link" href=""></a><span id="next-move-error"></span>
</div>`)

	moveInput := js.Global.Get("document").Call("getElementById", "move-input")
	if moveInput == nil {
		document.Call("write", "Next move input element not found")
	} else {
		moveInput.Call(
			"addEventListener",
			"keyup",
			func(event *js.Object) {
				if keycode := event.Get("keyCode").Int(); keycode == 13 {
					nextMoveErrorElement := js.Global.Get("document").Call("getElementById", "next-move-error")
					if nextMoveErrorElement == nil {
						document.Call("write", "Next move error element not found")
						return
					}
					nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link")
					if nextMoveLink == nil {
						nextMoveErrorElement.Set("innerHTML", "Next move link element not found")
						return
					}

					nextMoveErrorElement.Set("innerHTML", "")
					nextMoveLink.Set("innerHTML", "")
					nextMoveLink.Set("href", "")

					nextMovePCM := strings.TrimSpace(moveInput.Get("value").String())
					if nextMovePCM == "" {
						return
					}

					nextMove := move.Parse(nextMovePCM)
					if nextMove == move.Null {
						nextMoveErrorElement.Set("innerHTML", "Next move is not in PCN format")
						return
					}

					if _, ok := g.LegalMoves()[nextMove]; ok == false {
						nextMoveErrorElement.Set("innerHTML", "Next move is not a legal move")
						return
					}

					nextMoveString, err := encodeMove(nextMove) // encoding.go
					if err != nil {
						nextMoveErrorElement.Set("innerHTML", "Next move encoding error: "+err.Error())
						return
					}

					url := location.Get("origin").String() + location.Get("pathname").String() + "?" + movesString + nextMoveString
					nextMoveLink.Set("innerHTML", url)
					nextMoveLink.Set("href", url)
					nextMoveErrorElement.Set("innerHTML", " <- copy this link and send to your oponent")
				}
			},
			false,
		)
		moveInput.Call("focus")
	}
}

func fillChessBoard(g *game.Game, nextMove move.Move) error {
	// fill board grid with markers and pieces

	drawPosition := g.Position()
	// precalculate next move markers and stuff
	var nextMoveError error
	nextMoveMarkerClasses := map[square.Square][]string{}
	if nextMove == move.Null {
		// no next move
	} else {
		// some move, legal or illegal or incomplete
		if _, ok := g.Position().LegalMoves()[nextMove]; ok {
			// legal move

			// fill marker classes to squares
			nextMoveMarkerClasses[nextMove.Source] = []string{"next-move", "next-move-from"}
			nextMoveMarkerClasses[nextMove.Destination] = []string{"next-move", "next-move-to"}

			//set drawing position to next move
			drawPosition = g.Position().MakeMove(nextMove)

			//TODO generate next move link + stuff

		} else {
			// illegal move. but why?
			if nextMove.Source == square.NoSquare {
				// from not filled
			} else {
				// from filled. is it legal?
				legalFromMoves := map[move.Move]struct{}{}
				for move, _ := range g.Position().LegalMoves() {
					if move.Source == nextMove.Source {
						legalFromMoves[move] = struct{}{}
					}
				}
				if len(legalFromMoves) > 0 {
					// from is legal

					// mark from square
					nextMoveMarkerClasses[nextMove.Source] = []string{"next-move", "next-move-from"}

					// from is legal, what about others?
					if nextMove.Destination == square.NoSquare {
						// to not filled

						// mark possible to squares
						for move, _ := range legalFromMoves {
							if nextMoveMarkerClasses[move.Destination] == nil {
								nextMoveMarkerClasses[move.Destination] = []string{"next-move", "next-move-possible-to"}
							}
						}
					} else {
						// to filled. is it legal?
						legalFromToMoves := map[move.Move]struct{}{}
						for move, _ := range legalFromMoves {
							if move.Destination == nextMove.Destination {
								legalFromToMoves[move] = struct{}{}
							}
						}
						if len(legalFromToMoves) > 0 {
							// to is also legal

							// mark from square
							nextMoveMarkerClasses[nextMove.Destination] = []string{"next-move", "next-move-to"}

							// to is also legal. but the whole move is illegal. there have to be a promotion behind it
							if nextMove.Promote == piece.None {
								// jop, promote not filled, do something about it
								//TODO promotion
							} else {
								// promote is filled, but is illegal
								nextMoveError = errors.New("next move promotion is illegal! from: " + nextMove.Source.String() + ", to: " + nextMove.Destination.String() + ", promote: " + nextMove.Promote.String())
							}
						} else {
							// to is illegal
							nextMoveError = errors.New("next move to square is illegal! from: " + nextMove.Source.String() + ", to: " + nextMove.Destination.String())
						}
					}
				} else {
					// from is illegal
					nextMoveError = errors.New("next move from square is illegal! from: " + nextMove.Source.String())
				}
			}
		}
	}

	for i := int(63); i >= 0; i-- {
		sq := square.Square(i)
		sqElm := js.Global.Get("document").Call("getElementById", sq.String())
		if sqElm == nil {
			return errors.New("Can't find square element: " + sq.String())
		}

		marker := ""
		{
			// marker
			markerClasses := []string{"marker"}

			{
				// last move
				lm := g.Position().LastMove
				if lm != move.Null && (lm.Source == sq || lm.Destination == sq) {
					// last-move from or to marker is on square
					dir := "from"
					if lm.Destination == sq {
						dir = "to"
					}
					markerClasses = append(markerClasses, "last-move", "last-move-"+dir)
				}
			}

			{
				// next move
				// use precomputed classes
				if m, ok := nextMoveMarkerClasses[sq]; ok {
					markerClasses = append(markerClasses, m...)
				}
			}

			if len(markerClasses) > 1 {
				marker = "<span class=\"" + strings.Join(markerClasses, " ") + "\"></span>"
			}
		}

		pc := drawPosition.OnSquare(sq)
		sqElm.Set("innerHTML", marker+piecesString[pc.Color][pc.Type])
	}

	if nextMoveError != nil {
		return nextMoveError
	}

	// fill game status
	if gameProgressElement := js.Global.Get("document").Call("getElementById", "game-progress"); gameProgressElement != nil {
		gameProgressElement.Set("innerHTML", g.Status().String())
	}

	if gameMovingPlayerElement := js.Global.Get("document").Call("getElementById", "moving-player"); gameMovingPlayerElement != nil {
		if g.Status() != game.InProgress {
			// game has ended
			gameMovingPlayerElement.Set("innerHTML", "")
		} else {
			// game is in progress
			gameMovingPlayerElement.Set("innerHTML", piecesString[g.ActiveColor()][piece.King])
		}
	}

	return nil
}
