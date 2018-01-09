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

	app := app{movesString, g, move.Null, map[square.Square]func(square.Square, *js.Object){}}
	app.drawBoard()

	if err := app.updateBoard(); err != nil {
		document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
		return
	}
}

type app struct {
	movesString    string
	game           *game.Game
	nextMove       move.Move
	squaresHandler map[square.Square]func(sq square.Square, event *js.Object)
}

// Draws chess board, game-status & next-move elements to document
func (app app) drawBoard() {
	document := js.Global.Get("document")

	// is rotation supported?
	rotationSupported := false
	if div := document.Call("createElement", "div"); div != js.Undefined {
		if div.Get("style").Get("transform") != js.Undefined {
			rotationSupported = true
		}
		div.Call("remove")
	}

	rotateBoard180deg := rotationSupported && app.game.ActiveColor() == piece.Black

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

		if etr := document.Call("getElementById", "edging-top-right"); etr != nil {
			etr.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					if board := document.Call("getElementById", "board"); board != nil {
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

		if etr := document.Call("getElementById", "edging-bottom-left"); etr != nil {
			etr.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					if board := document.Call("getElementById", "board"); board != nil {
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
		document.Call("write", "<div id=\"promotion-overlay\">")
		document.Call("write", "</div>")
		// promotion overlay
	}

	{
		// board
		document.Call("write", "</div>")
	}

	document.Call("write", `<div id="next-move">
<a id="next-move-link" href=""></a> <- copy this link and send to your oponent
</div>`)

	document.Call("write", `<div id="game-status">
		<p>Moving player: <span id="moving-player"></span></p>
		<p>Game status: <span id="game-progress"><span></p>
</div>`)

	// map click events to grid squares
	for i := int(63); i >= 0; i-- {
		sq := square.Square(i)
		if sqElm := document.Call("getElementById", sq.String()); sqElm != nil {
			sqElm.Call(
				"addEventListener",
				"click",
				app.squareHandler,
				false,
			)
		}
	}
}

// Fills board grid with markers and pieces, updates status and next-move elements,
// assign handler functions to grid squares
func (app *app) updateBoard() error {
	{
		// clear playground

		//TODO clear board, status

		// hide promotion overlay
		if promotionOverlay := js.Global.Get("document").Call("getElementById", "promotion-overlay"); promotionOverlay != nil {
			promotionOverlay.Get("classList").Call("remove", "show")
		}

		// clear next-move
		if nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link"); nextMoveLink != nil {
			nextMoveLink.Set("innerHTML", "")
			nextMoveLink.Set("href", "")
		}
		if nextMove := js.Global.Get("document").Call("getElementById", "next-move"); nextMove != nil {
			nextMove.Get("classList").Call("remove", "show")
		}

		// clear grid squares handler, ...
		for i := int(63); i >= 0; i-- {
			delete(app.squaresHandler, square.Square(i))
		}
	}

	drawPosition := app.game.Position()
	// precalculate next move markers and stuff
	var nextMoveError error
	nextMoveMarkerClasses := map[square.Square][]string{}
	if app.nextMove == move.Null {
		// no next move

		// set handlers for moving player pieces
		for i := int(63); i >= 0; i-- {
			sq := square.Square(i)
			pc := drawPosition.OnSquare(sq)
			if pc.Color == app.game.Position().ActiveColor {
				app.squaresHandler[sq] = func(sq square.Square, _ *js.Object) {
					app.nextMove.Source = sq
					if err := app.updateBoard(); err != nil {
						js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
						return
					}
				}
			}
		}
	} else {
		// some move, legal or illegal or incomplete
		if _, ok := app.game.Position().LegalMoves()[app.nextMove]; ok {
			// legal move

			// fill marker classes to squares
			nextMoveMarkerClasses[app.nextMove.Source] = []string{"next-move", "next-move-from"}
			nextMoveMarkerClasses[app.nextMove.Destination] = []string{"next-move", "next-move-to"}

			//set drawing position to next move
			drawPosition = app.game.Position().MakeMove(app.nextMove)

			nextMoveString, err := encodeMove(app.nextMove) // encoding.go
			if err != nil {
				nextMoveError = errors.New("Next move encoding error: " + err.Error())
			} else {
				location := js.Global.Get("location")
				url := location.Get("origin").String() + location.Get("pathname").String() + "?" + app.movesString + nextMoveString

				if nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link"); nextMoveLink != nil {
					nextMoveLink.Set("innerHTML", url)
					nextMoveLink.Set("href", url)
				}

				if nextMove := js.Global.Get("document").Call("getElementById", "next-move"); nextMove != nil {
					nextMove.Get("classList").Call("add", "show")
				}
			}
		} else {
			// illegal move

			// set handlers for moving player pieces
			for i := int(63); i >= 0; i-- {
				sq := square.Square(i)
				pc := drawPosition.OnSquare(sq)
				if pc.Color == app.game.Position().ActiveColor {
					app.squaresHandler[sq] = func(sq square.Square, _ *js.Object) {
						app.nextMove.Source = sq
						if err := app.updateBoard(); err != nil {
							js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
							return
						}
					}
				}
			}

			// illegal move. but why?
			if app.nextMove.Source == square.NoSquare {
				// from not filled
				// this should not happen
				nextMoveError = errors.New("next move is not null, but has no from square filled")
			} else {
				// from filled

				// mark from square
				nextMoveMarkerClasses[app.nextMove.Source] = []string{"next-move", "next-move-from"}

				// remove from handlers, to unhiglight if clicking on the same piece again
				delete(app.squaresHandler, app.nextMove.Source)

				// from filled. is it legal?
				legalFromMoves := map[move.Move]struct{}{}
				for move, _ := range app.game.Position().LegalMoves() {
					if move.Source == app.nextMove.Source {
						legalFromMoves[move] = struct{}{}
					}
				}
				if len(legalFromMoves) > 0 {
					// from is legal

					// from is legal, what about others?
					if app.nextMove.Destination == square.NoSquare {
						// to not filled

						// mark possible to squares
						for move, _ := range legalFromMoves {
							if nextMoveMarkerClasses[move.Destination] == nil {
								nextMoveMarkerClasses[move.Destination] = []string{"next-move", "next-move-possible-to"}
							}
						}

						// add handlers to possible next move
						for move, _ := range legalFromMoves {
							app.squaresHandler[move.Destination] = func(sq square.Square, _ *js.Object) {
								app.nextMove.Destination = sq
								if err := app.updateBoard(); err != nil {
									js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
									return
								}
							}
						}
					} else {
						// to filled. is it legal?
						legalFromToMoves := map[move.Move]struct{}{}
						for move, _ := range legalFromMoves {
							if move.Destination == app.nextMove.Destination {
								legalFromToMoves[move] = struct{}{}
							}
						}
						if len(legalFromToMoves) > 0 {
							// to is also legal

							// mark from square
							nextMoveMarkerClasses[app.nextMove.Destination] = []string{"next-move", "next-move-to"}

							// to is also legal. but the whole move is illegal. there have to be a promotion behind it
							if app.nextMove.Promote == piece.None {
								// jop, promote not filled, do something about it
								//TODO promotion
								if promotionOverlay := js.Global.Get("document").Call("getElementById", "promotion-overlay"); promotionOverlay != nil {
									innerHTML := ""
									for i := 0; i <= int(piece.King); i++ {
										if _, ok := legalFromToMoves[move.Move{
											Source:      app.nextMove.Source,
											Destination: app.nextMove.Destination,
											Promote:     piece.Type(i),
										}]; ok {
											innerHTML += piecesString[app.game.Position().ActiveColor][piece.Type(i)]
										}
									}
									promotionOverlay.Set("innerHTML", innerHTML)

									//TODO add handlers for promotion pieces
									promotionOverlay.Get("classList").Call("add", "show")
								}
							} else {
								// promote is filled, but is illegal
								nextMoveError = errors.New("next move promotion is illegal! from: " + app.nextMove.Source.String() + ", to: " + app.nextMove.Destination.String() + ", promote: " + app.nextMove.Promote.String())
							}
						} else {
							// to is illegal
							nextMoveError = errors.New("next move to square is illegal! from: " + app.nextMove.Source.String() + ", to: " + app.nextMove.Destination.String())
						}
					}
				} else {
					// from is illegal
					nextMoveError = errors.New("next move from square is illegal! from: " + app.nextMove.Source.String())
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
				lm := app.game.Position().LastMove
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
		gameProgressElement.Set("innerHTML", app.game.Status().String())
	}

	if gameMovingPlayerElement := js.Global.Get("document").Call("getElementById", "moving-player"); gameMovingPlayerElement != nil {
		if app.game.Status() != game.InProgress {
			// game has ended
			gameMovingPlayerElement.Set("innerHTML", "")
		} else {
			// game is in progress
			gameMovingPlayerElement.Set("innerHTML", piecesString[app.game.ActiveColor()][piece.King])
		}
	}

	return nil
}

func (app *app) squareHandler(event *js.Object) {
	elm := event.Get("currentTarget")
	if elm == nil || elm == js.Undefined {
		js.Global.Call("alert", "no current target element")
		return
	}
	elmId := elm.Call("getAttribute", "id").String()
	if strings.TrimSpace(elmId) == "" {
		js.Global.Call("alert", "no id attribute in target")
		return
	}

	sq := square.Parse(elmId)
	handler, ok := app.squaresHandler[sq]
	if !ok {
		//js.Global.Call("alert", "no "+elmId+" square handler")
		app.nextMove = move.Null
		if err := app.updateBoard(); err != nil {
			js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
			return
		}
		return
	}

	//js.Global.Call("alert", "handle "+elmId+" square")
	handler(sq, event)
}
