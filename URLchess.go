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

var playablePiecesType = []piece.Type{piece.Pawn, piece.Rook, piece.Knight, piece.Bishop, piece.Queen, piece.King}
var promotablePiecesType = []piece.Type{piece.Rook, piece.Knight, piece.Bishop, piece.Queen}

var pieceTypesToName = map[piece.Type]string{
	piece.Pawn:   "pawn",
	piece.Rook:   "rook",
	piece.Knight: "knight",
	piece.Bishop: "bishop",
	piece.Queen:  "queen",
	piece.King:   "king",
}

var pieceNamesToType map[string]piece.Type = func() map[string]piece.Type {
	res := make(map[string]piece.Type, len(pieceTypesToName))
	for k, v := range pieceTypesToName {
		if _, ok := res[v]; ok {
			panic("error creating pieceNamesToType map: duplicate name: " + v)
		}
		res[v] = k
	}
	return res
}()

var piecesToString map[piece.Piece]string = func() map[piece.Piece]string {
	res := make(map[piece.Piece]string, len(piece.Colors)*len(pieceTypesToName))
	for _, pieceColor := range piece.Colors {
		for pieceType, pieceName := range pieceTypesToName {
			p := piece.New(pieceColor, pieceType)
			res[p] = "<span class=\"piece " + strings.ToLower(pieceColor.String()) + " " + pieceName + "\">" + p.Figurine() + "</span>"
		}
	}
	return res
}()

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
	gtos := make(gameThrowOuts, len(moves))

	{ // apply game moves
		for i, move := range moves {
			if g.Status() != game.InProgress {
				document.Call("write", "Too many moves in url string! "+strconv.Itoa(i+1)+" moves are enough")
				return
			}

			// position before move
			pbm := g.Position()

			_, merr := g.MakeMove(move)
			if merr != nil {
				document.Call("write", "Errorneous move number "+strconv.Itoa(i+1)+": "+merr.Error())
				return
			}

			// throw outs for this move
			tos := throwOuts{}

			// copy previous move throw outs
			if i > 0 {
				for p, c := range gtos[i-1] {
					tos[p] = c
				}
			}

			// was a piece thrown out = move destination contains some piece
			if p := pbm.OnSquare(move.To()); p.Type != piece.None {
				if _, ok := tos[p]; !ok {
					tos[p] = 0
				}
				tos[p]++
			}

			gtos[i] = tos
		}
	}

	app := app{movesString, g, gtos, move.Null, map[square.Square]func(square.Square, *js.Object){}}
	app.drawBoard()

	//js.Global.Call("alert", "calling update board from main")
	if err := app.updateBoard(); err != nil {
		document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
		return
	}
}

type app struct {
	movesString    string
	game           *game.Game
	gameGc         gameThrowOuts
	nextMove       move.Move
	squaresHandler map[square.Square]func(sq square.Square, event *js.Object)
}

// Draws chess board, game-status and next-move elements to document.
// Also sets event listeners for grid and undo next move
func (app *app) drawBoard() {
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

	{ // board
		class := ""
		if rotateBoard180deg {
			class = " class=\"rotated180deg\""
		}
		document.Call("write", "<div id=\"board\""+class+">")
	}

	{ // edging-top-left
		document.Call("write", "<div id=\"edging-top-left\" class=\"edging corner\">")
		document.Call("write", "</div>")
	}

	{ // edging-top
		document.Call("write", "<div id=\"edging-top\" class=\"edging horizontal\">")
		for i := 0; i < 8; i++ {
			document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
		}
		document.Call("write", "</div>")
	}

	{ // edging-top-right
		document.Call("write", "<div id=\"edging-top-right\" class=\"edging corner\">")
		document.Call("write", "</div>")

		if edging := document.Call("getElementById", "edging-top-right"); edging != nil && rotationSupported {
			edging.Set("innerHTML", "↻")
			edging.Call(
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

	{ // edging-left
		document.Call("write", "<div id=\"edging-left\" class=\"edging vertical\">")
		for i := 8; i > 0; i-- {
			document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
		}
		document.Call("write", "</div>")
	}

	{ // grid
		document.Call("write", "<div class=\"grid\">")
		squareTones := []string{"light-square", "dark-square"}
		for i := int(63); i >= 0; i-- {
			document.Call("write", "<div id=\""+square.Square(i).String()+"\" class=\""+squareTones[(i%8+i/8)%2]+"\"></div>")
		}
		document.Call("write", "</div>")
	}

	{ // edging-right
		document.Call("write", "<div id=\"edging-right\" class=\"edging vertical\">")
		for i := 8; i > 0; i-- {
			document.Call("write", "<div>"+strconv.Itoa(i)+"</div>")
		}
		document.Call("write", "</div>")
	}

	{ // edging-bottom-left
		document.Call("write", "<div id=\"edging-bottom-left\" class=\"edging corner\">")
		document.Call("write", "</div>")

		if edging := document.Call("getElementById", "edging-bottom-left"); edging != nil && rotationSupported {
			edging.Set("innerHTML", "↻")
			edging.Call(
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

	{ // edging-bottom
		document.Call("write", "<div id=\"edging-bottom\" class=\"edging horizontal\">")
		for i := 0; i < 8; i++ {
			document.Call("write", "<div>"+string(rune('a'+i))+"</div>")
		}
		document.Call("write", "</div>")
	}

	{ // edging-bottom-right
		document.Call("write", "<div id=\"edging-bottom-right\" class=\"edging corner\">")
		document.Call("write", "</div>")
	}

	{ // promotion overlay
		document.Call("write", "<div id=\"promotion-overlay\">")
		for _, pieceType := range promotablePiecesType {
			document.Call("write", "<span id=\"promote-to-"+pieceTypesToName[pieceType]+"\" class=\"piece\" piece=\""+pieceTypesToName[pieceType]+"\"></span>")
		}
		document.Call("write", "</div>")
	}

	{ // board
		document.Call("write", "</div>")
	}

	{ // thrown out pieces
		document.Call("write", "<div id=\"thrown-outs-container\">")
		for _, c := range piece.Colors {
			document.Call("write", "<div id=\"thrown-outs-"+strings.ToLower(c.String())+"\" class=\"thrown-outs\"></div>")
		}
		document.Call("write", "</div>")
	}

	{ // next move
		hintString := "copy this link and send it to your oponent"

		if exec := js.Global.Get("document").Get("execCommand"); exec != nil && exec != js.Undefined {
			hintString = "this link has been copied to clipboard, send it to your oponent"
		}

		document.Call("write", `<div id="next-move" class="hidden">
	<p class="link">
		Next move URL link:
		<input id="next-move-input" readonly="readonly" value=""/>
		<span class="hint">`+hintString+`</span>
	</p>
  <p class="actions">
		<a id="next-move-link" href="">make move</a>
		<a id="next-move-back" href="">undo move</a>
	</p>
</div>`)
	}

	{ // game status
		document.Call("write", `<div id="game-status">
	<p id="game-status-text">... loading ...</p>
	<p id="game-status-player">`+piecesToString[piece.New(piece.White, piece.King)]+piecesToString[piece.New(piece.Black, piece.King)]+`</p>
</div>`)
	}

	{ // event listeners

		// next move back
		if back := document.Call("getElementById", "next-move-back"); back != nil {
			back.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					event.Call("preventDefault")
					app.nextMove = move.Null
					//js.Global.Call("alert", "calling update board from next-move-back listener")
					if err := app.updateBoard(); err != nil {
						document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
						return
					}
				},
				false,
			)
		}
		// next move link ckick
		if ebl := document.Call("getElementById", "edging-bottom-left"); ebl != nil {
			if link := document.Call("getElementById", "next-move-link"); link != nil {
				link.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if board := document.Call("getElementById", "board"); board != nil {

							// rotate only if needed
							shouldBeRotatedInNextMove := app.game.Position().ActiveColor != piece.Black
							isRotated := board.Get("classList").Call("contains", "rotated180deg").Bool()

							if shouldBeRotatedInNextMove != isRotated {
								//rotate, wait, cange location
								event.Call("preventDefault")
								if isRotated {
									board.Get("classList").Call("remove", "rotated180deg")
								} else {
									board.Get("classList").Call("add", "rotated180deg")
								}

								url := link.Call("getAttribute", "href")
								js.Global.Call(
									"setTimeout",
									func() {
										js.Global.Get("location").Set("href", url)
									},
									800,
								)
							}
							if nm := document.Call("getElementById", "next-move"); nm != nil {
								nm.Get("classList").Call("add", "hidden")
							}

						}
					},
					false,
				)
			}
		}

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

		// promotion overlay
		if promotionOverlay := document.Call("getElementById", "promotion-overlay"); promotionOverlay != nil {
			promotionOverlay.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					event.Call("preventDefault")
					app.nextMove = move.Null
					//js.Global.Call("alert", "calling update board from promotion-overlay listener")
					if err := app.updateBoard(); err != nil {
						document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
						return
					}
				},
				false,
			)
			// add handlers for promotion pieces
			for _, pt := range promotablePiecesType {
				if promotionPiece := js.Global.Get("document").Call("getElementById", "promote-to-"+pieceTypesToName[pt]); promotionPiece != nil {
					promotionPiece.Call(
						"addEventListener",
						"click",
						func(event *js.Object) {
							event.Call("stopPropagation")
							if elm := event.Get("currentTarget"); elm != nil {
								pieceName := elm.Call("getAttribute", "piece").String()
								if pt, ok := pieceNamesToType[pieceName]; ok {
									app.nextMove.Promote = piece.Type(pt)
									//js.Global.Call("alert", "promote to: "+app.nextMove.Promote.String())
									//js.Global.Call("alert", "calling update board from promotion-piece "+pieceName+" listener")
									if err := app.updateBoard(); err != nil {
										js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
										return
									}
								}
							}
						},
						false,
					)
				}
			}
		}
	}
}

// Fills board grid with markers and pieces, updates status and next-move elements,
// assign handler functions to grid squares
func (app *app) updateBoard() error {
	//js.Global.Call("alert", "update: nextMove: "+app.nextMove.String())
	{ // clear playground

		//TODO clear board

		// clear thrown out pieces
		for _, c := range piece.Colors {
			if tos := js.Global.Get("document").Call("getElementById", "thrown-outs-"+strings.ToLower(c.String())); tos != nil {
				tos.Set("innerHTML", "")
			}
		}

		// clear status
		if gameStatusText := js.Global.Get("document").Call("getElementById", "game-status-text"); gameStatusText != nil {
			gameStatusText.Set("innerHTML", "")
		}
		if gameStatusPlayer := js.Global.Get("document").Call("getElementById", "game-status-player"); gameStatusPlayer != nil {
			gameStatusPlayer.Set("innerHTML", "")
		}

		// hide promotion overlay
		if promotionOverlay := js.Global.Get("document").Call("getElementById", "promotion-overlay"); promotionOverlay != nil {
			promotionOverlay.Get("classList").Call("remove", "show")

			//TODO clear promotion pieces elements
		}

		// clear next-move
		if nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link"); nextMoveLink != nil {
			nextMoveLink.Set("href", "")
		}
		if nextMoveInput := js.Global.Get("document").Call("getElementById", "next-move-input"); nextMoveInput != nil {
			nextMoveInput.Set("value", "")
		}
		if nextMove := js.Global.Get("document").Call("getElementById", "next-move"); nextMove != nil {
			nextMove.Get("classList").Call("add", "hidden")
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
		//js.Global.Call("alert", "no next move")

		// set handlers for moving player pieces
		for i := int(63); i >= 0; i-- {
			sq := square.Square(i)
			pc := drawPosition.OnSquare(sq)
			if pc.Color == app.game.Position().ActiveColor {
				app.squaresHandler[sq] = func(sq square.Square, _ *js.Object) {
					app.nextMove.Source = sq
					//js.Global.Call("alert", "calling update board from square handler moving piece")
					if err := app.updateBoard(); err != nil {
						js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
						return
					}
				}
			}
		}
	} else {
		// some move, legal or illegal or incomplete
		//js.Global.Call("alert", "some move, legal or illegal or incomplete")

		color := strings.ToLower(app.game.Position().ActiveColor.String())
		if _, ok := app.game.Position().LegalMoves()[app.nextMove]; ok {
			// legal move
			//js.Global.Call("alert", "legal move")

			// fill marker classes to squares
			nextMoveMarkerClasses[app.nextMove.Source] = []string{"next-move", "next-move-" + color, "next-move-from"}
			nextMoveMarkerClasses[app.nextMove.Destination] = []string{"next-move", "next-move-" + color, "next-move-to"}

			//set drawing position to next move
			drawPosition = app.game.Position().MakeMove(app.nextMove)

			nextMoveString, err := encodeMove(app.nextMove) // encoding.go
			if err != nil {
				nextMoveError = errors.New("Next move encoding error: " + err.Error())
			} else {
				location := js.Global.Get("location")
				url := location.Get("origin").String() + location.Get("pathname").String() + "?" + app.movesString + nextMoveString

				if nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link"); nextMoveLink != nil {
					nextMoveLink.Set("href", url)
				}
				if nextMoveInput := js.Global.Get("document").Call("getElementById", "next-move-input"); nextMoveInput != nil {
					nextMoveInput.Set("value", url)
				}

				if nextMove := js.Global.Get("document").Call("getElementById", "next-move"); nextMove != nil {
					nextMove.Get("classList").Call("remove", "hidden")
				}

				// select input text & copy to clipboard
				if nextMoveInput := js.Global.Get("document").Call("getElementById", "next-move-input"); nextMoveInput != nil {
					nextMoveInput.Call("select")
					js.Global.Get("document").Call("execCommand", "Copy")
					nextMoveInput.Call("blur")
				}
			}
		} else {
			// illegal move
			//js.Global.Call("alert", "illegal move")

			// set handlers for moving player pieces
			for i := int(63); i >= 0; i-- {
				sq := square.Square(i)
				pc := drawPosition.OnSquare(sq)
				if pc.Color == app.game.Position().ActiveColor {
					app.squaresHandler[sq] = func(sq square.Square, _ *js.Object) {
						app.nextMove.Source = sq
						//js.Global.Call("alert", "calling update board from square handler moving piece illegal move")
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
				//js.Global.Call("alert", "from filled")

				// mark from square
				nextMoveMarkerClasses[app.nextMove.Source] = []string{"next-move", "next-move-" + color, "next-move-from"}

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
					//js.Global.Call("alert", "from is legal")

					// from is legal, what about others?
					if app.nextMove.Destination == square.NoSquare {
						// to not filled
						//js.Global.Call("alert", "to not filled")

						// mark possible to squares
						for move, _ := range legalFromMoves {
							if nextMoveMarkerClasses[move.Destination] == nil {
								// add next-move mark
								oponentColor := piece.Colors[(int(app.game.Position().ActiveColor)+1)%2]
								if app.game.Position().OnSquare(move.Destination).Color == oponentColor {
									nextMoveMarkerClasses[move.Destination] = []string{"next-move", "next-move-" + color, "next-move-possible-to", "kill"}
								} else {
									nextMoveMarkerClasses[move.Destination] = []string{"next-move", "next-move-" + color, "next-move-possible-to", "nokill"}
								}
							}
						}

						// add handlers to possible next move
						for move, _ := range legalFromMoves {
							app.squaresHandler[move.Destination] = func(sq square.Square, _ *js.Object) {
								app.nextMove.Destination = sq
								//js.Global.Call("alert", "calling update board from square handler possible next move")
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
							//js.Global.Call("alert", "to is also legal")

							// mark from square
							nextMoveMarkerClasses[app.nextMove.Destination] = []string{"next-move", "next-move-" + color, "next-move-to"}

							// to is also legal. but the whole move is illegal. there have to be a promotion behind it
							if app.nextMove.Promote == piece.None {
								// promote not filled, do something about it
								//js.Global.Call("alert", "promote not filled, do something about it")
								if promotionOverlay := js.Global.Get("document").Call("getElementById", "promotion-overlay"); promotionOverlay != nil {

									//TODO? hide unpromotable pieces

									// fill piece figure to promotable pieces elements
									for _, pt := range promotablePiecesType {
										if _, ok := legalFromToMoves[move.Move{
											Source:      app.nextMove.Source,
											Destination: app.nextMove.Destination,
											Promote:     pt,
										}]; ok {
											if promotionPiece := js.Global.Get("document").Call("getElementById", "promote-to-"+pieceTypesToName[pt]); promotionPiece != nil {
												promotionPiece.Set("innerHTML", piece.New(app.game.Position().ActiveColor, pt).Figurine())
											} else {
												//TODO return error
											}
										}
									}
									promotionOverlay.Get("classList").Call("add", "show")
								}
							} else {
								// promote is filled, but is illegal
								nextMoveError = errors.New("next move promotion is illegal! from: " + app.nextMove.Source.String() + ", to: " + app.nextMove.Destination.String() + ", promote: " + app.nextMove.Promote.String())
							}
						} else {
							// to is illegal
							nextMoveError = errors.New("next move to square is illegal! from: " + app.nextMove.Source.String() + ", to: " + app.nextMove.Destination.String())
							//TODO repair & update & test
						}
					}
				} else {
					// from is illegal

					// but if from is active piece and to and promotion are empty,
					// dont throw error, cause this is just piece with no legal moves selection
					if !(app.game.Position().OnSquare(app.nextMove.Source).Color == app.game.Position().ActiveColor &&
						app.nextMove.Destination == square.NoSquare &&
						app.nextMove.Promote == piece.None) {
						nextMoveError = errors.New("next move from square is illegal! from: " + app.nextMove.Source.String())
					}
				}
			}
		}
	}

	//js.Global.Call("alert", "drawing grid")
	for i := int(63); i >= 0; i-- {
		sq := square.Square(i)
		sqElm := js.Global.Get("document").Call("getElementById", sq.String())
		if sqElm == nil {
			return errors.New("Can't find square element: " + sq.String())
		}

		markerClasses := []string{"marker"}
		{ // marker classes fill

			{ // last move
				lm := app.game.Position().LastMove
				if lm != move.Null && (lm.Source == sq || lm.Destination == sq) {
					// last-move from or to marker is on square
					dir := "from"
					if lm.Destination == sq {
						dir = "to"
					}
					oponentColor := strings.ToLower(piece.Colors[(int(app.game.Position().ActiveColor)+1)%2].String())
					markerClasses = append(markerClasses, "last-move", "last-move-"+oponentColor, "last-move-"+dir)
				}
			}

			{ // next move
				// use precomputed classes
				if m, ok := nextMoveMarkerClasses[sq]; ok {
					markerClasses = append(markerClasses, m...)
				}
			}
		}

		pc := drawPosition.OnSquare(sq)
		innerHTML := "<span class=\"" + strings.Join(markerClasses, " ") + "\">" + piecesToString[pc] + "</span>"

		sqElm.Set("innerHTML", innerHTML)
	}

	if nextMoveError != nil {
		return nextMoveError
	}

	//js.Global.Call("alert", len(app.gameGc))

	// fill thrown out pieces
	if gcl := len(app.gameGc); gcl > 0 {
		for _, c := range piece.Colors {
			id := "thrown-outs-" + strings.ToLower(c.String())
			if tosElm := js.Global.Get("document").Call("getElementById", id); tosElm != nil {
				tos := throwOuts{}
				for p, n := range app.gameGc[gcl-1] {
					if p.Color == c {
						tos[p] = n
					}
				}
				tosElmStr := ""

				if len(tos) != 0 {
					for _, pt := range playablePiecesType {
						if n, ok := tos[piece.New(c, pt)]; ok {
							tosElmStr += "<div class=\"piececount\">"
							tosElmStr += piecesToString[piece.New(c, pt)]
							tosElmStr += "<span class=\"count\">" + strconv.Itoa(int(n)) + "</span>"
							tosElmStr += "</div>"
						}
					}
				}

				tosElm.Set("innerHTML", tosElmStr)
			}
		}
	}

	// fill game status
	text := "Moving player"
	player := ""
	if st := app.game.Status(); st != game.InProgress {
		text = st.String()
		if st&game.Draw != 0 {
			player = piecesToString[piece.New(piece.White, piece.King)] + piecesToString[piece.New(piece.Black, piece.King)]
		} else if st&game.WhiteWon != 0 {
			player = piecesToString[piece.New(piece.White, piece.King)]
		} else if st&game.BlackWon != 0 {
			player = piecesToString[piece.New(piece.Black, piece.King)]
		}
	} else {
		player = piecesToString[piece.New(app.game.ActiveColor(), piece.King)]
	}

	if gameStatusText := js.Global.Get("document").Call("getElementById", "game-status-text"); gameStatusText != nil {
		gameStatusText.Set("innerHTML", text)
	}
	if gameStatusPlayer := js.Global.Get("document").Call("getElementById", "game-status-player"); gameStatusPlayer != nil {
		gameStatusPlayer.Set("innerHTML", player)
	}

	//js.Global.Call("alert", "end")

	return nil
}

func (app *app) squareHandler(event *js.Object) {
	elm := event.Get("currentTarget")
	if elm == nil || elm == js.Undefined {
		//js.Global.Call("alert", "no current target element")
		return
	}
	elmId := elm.Call("getAttribute", "id").String()
	if strings.TrimSpace(elmId) == "" {
		//js.Global.Call("alert", "no id attribute in target")
		return
	}

	sq := square.Parse(elmId)
	handler, ok := app.squaresHandler[sq]
	if !ok {
		//js.Global.Call("alert", "no "+elmId+" square handler")
		app.nextMove = move.Null
		//js.Global.Call("alert", "calling update board from square handler no handler")
		if err := app.updateBoard(); err != nil {
			js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
			return
		}
		return
	}

	//js.Global.Call("alert", "handle "+elmId+" square")
	handler(sq, event)
}

type gameThrowOuts []throwOuts
type throwOuts map[piece.Piece]uint8
