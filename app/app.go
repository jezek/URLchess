package app

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

type Elementer interface {
	Element() *js.Object
}

type Updater interface {
	Update(*AppEvents) (*js.Object, error)
}

type BoardEdging struct {
	elm      *js.Object
	Position string //top, bottom, left, right, top-left, top-right, bottom-left, bottom-right
}

func (e *BoardEdging) Update(_ *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if e.elm == nil {
		newElm = js.Global.Get("document").Call("createElement", "div")
		newElm.Set("id", "edging-"+e.Position)
		newElm.Get("classList").Call("add", "edging")

		e.elm = newElm
	}
	return newElm, nil
}

type EdgingHorizontal struct {
	BoardEdging
}

func (e *EdgingHorizontal) Update(events *AppEvents) (*js.Object, error) {
	newElm, err := e.BoardEdging.Update(events)
	if newElm != nil {
		newElm.Get("classList").Call("add", "horizontal")
		for i := 0; i < 8; i++ {
			letter := js.Global.Get("document").Call("createElement", "div")
			letter.Set("textContent", string(rune('a'+i)))
			newElm.Call("appendChild", letter)
		}
	}
	return newElm, err
}

type EdgingVertical struct {
	BoardEdging
}

func (e *EdgingVertical) Update(events *AppEvents) (*js.Object, error) {
	newElm, err := e.BoardEdging.Update(events)
	if newElm != nil {
		newElm.Get("classList").Call("add", "vertical")
		for i := 8; i > 0; i-- {
			number := js.Global.Get("document").Call("createElement", "div")
			number.Set("textContent", strconv.Itoa(i))
			newElm.Call("appendChild", number)
		}
	}
	return newElm, err
}

type EdgingCorner struct {
	BoardEdging
}

func (e *EdgingCorner) Update(events *AppEvents) (*js.Object, error) {
	newElm, err := e.BoardEdging.Update(events)
	if newElm != nil {
		newElm.Get("classList").Call("add", "corner")
	}
	return newElm, err
}

type EdgingCornerRotating struct {
	EdgingCorner
	Disabled bool
}

func (e *EdgingCornerRotating) Update(events *AppEvents) (*js.Object, error) {
	newElm, err := e.EdgingCorner.Update(events)
	if newElm != nil && !e.Disabled {
		newElm.Set("innerHTML", "↻")
		newElm.Get("classList").Call("add", "enabled")
		events.Click(newElm, e.Click)
	}
	return newElm, err
}

func (_ *EdgingCornerRotating) Click(g *ChessGame, m *HtmlModel) error {
	m.Rotated180deg = !m.Rotated180deg
	return nil
}

type GridSquare struct {
	elm     *js.Object
	Markers struct {
		LastMove struct {
			White, Black struct {
				From, To bool
			}
		}
		NextMove struct {
			White, Black struct {
				From, To, PossibleTo bool
			}
		}
		Check bool
	}
	Piece piece.Piece
}

func (s *GridSquare) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if s.elm == nil {
		// create main board element
		newElm = js.Global.Get("document").Call("createElement", "div")
		s.elm = newElm
	}

	// update square, generate content & replace
	marker := js.Global.Get("document").Call("createElement", "span")

	marker.Get("classList").Call("add", "marker")
	if s.Markers.LastMove.White.From {
		marker.Get("classList").Call("add", "last-move")
		marker.Get("classList").Call("add", "last-move-white")
		marker.Get("classList").Call("add", "last-move-from")
	}
	if s.Markers.LastMove.White.To {
		marker.Get("classList").Call("add", "last-move")
		marker.Get("classList").Call("add", "last-move-white")
		marker.Get("classList").Call("add", "last-move-to")
	}
	if s.Markers.LastMove.Black.From {
		marker.Get("classList").Call("add", "last-move")
		marker.Get("classList").Call("add", "last-move-black")
		marker.Get("classList").Call("add", "last-move-from")
	}
	if s.Markers.LastMove.Black.To {
		marker.Get("classList").Call("add", "last-move")
		marker.Get("classList").Call("add", "last-move-black")
		marker.Get("classList").Call("add", "last-move-to")
	}
	if s.Markers.NextMove.White.From {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-white")
		marker.Get("classList").Call("add", "next-move-from")
	}
	if s.Markers.NextMove.White.To {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-white")
		marker.Get("classList").Call("add", "next-move-to")
	}
	if s.Markers.NextMove.White.PossibleTo {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-white")
		marker.Get("classList").Call("add", "next-move-possible-to")
	}
	if s.Markers.NextMove.Black.From {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-black")
		marker.Get("classList").Call("add", "next-move-from")
	}
	if s.Markers.NextMove.Black.To {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-black")
		marker.Get("classList").Call("add", "next-move-to")
	}
	if s.Markers.NextMove.Black.PossibleTo {
		marker.Get("classList").Call("add", "next-move")
		marker.Get("classList").Call("add", "next-move-black")
		marker.Get("classList").Call("add", "next-move-possible-to")
	}
	if s.Markers.Check {
		marker.Get("classList").Call("add", "check")
	}

	if s.Piece.Type != piece.None {
		marker.Call("appendChild", pieceElement(s.Piece))
	}

	s.elm.Set("innerHTML", "")
	s.elm.Call("appendChild", marker)

	// events
	//if s.Piece.Type == piece.King {
	//	events.Click(s.elm, func(g ChessGame, m HtmlModel) (ChessGame, HtmlModel) {
	//		js.Global.Call("alert", "kingclick")
	//		return g, m
	//	})
	//}

	return newElm, nil
}

type BoardGrid struct {
	elm     *js.Object
	Squares [64]GridSquare
}

var BoardGridSquareTones = []string{"light-square", "dark-square"}

func (g *BoardGrid) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if g.elm == nil {
		// create main board element
		newElm = js.Global.Get("document").Call("createElement", "div")
		newElm.Get("classList").Call("add", "grid")

		g.elm = newElm
	}

	for i := int(63); i >= 0; i-- {
		if created, err := g.Squares[i].Update(events); err != nil {
			return nil, err
		} else if created != nil {
			created.Set("id", square.Square(i).String())
			created.Get("classList").Call("add", BoardGridSquareTones[(i%8+i/8)%2])

			g.elm.Call("appendChild", created)
		}
	}
	return newElm, nil
}

type BoardPromotionOverlay struct {
	elm   *js.Object
	Shown bool
}

func (p *BoardPromotionOverlay) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object

	if p.elm == nil {
		// create
		newElm = js.Global.Get("document").Call("createElement", "div")
		newElm.Set("id", "promotion-overlay")

		for _, pieceType := range promotablePiecesType {
			span := js.Global.Get("document").Call("createElement", "span")
			span.Set("id", "promote-to-"+pieceTypesToName[pieceType])
			span.Get("classList").Call("add", "piece")
			span.Set("piece", pieceTypesToName[pieceType])

			newElm.Call("appendChild", span)
		}
		p.elm = newElm
	}

	if p.Shown {
		p.elm.Get("classList").Call("add", "show")
	} else {
		p.elm.Get("classList").Call("remove", "show")
	}

	return newElm, nil
}

type ModelBoard struct {
	elm     *js.Object
	Edgings struct {
		Top, Bottom          EdgingHorizontal
		Left, Right          EdgingVertical
		TopRight, BottomLeft EdgingCornerRotating
		TopLeft, BottomRight EdgingCorner
	}
	Grid             BoardGrid
	PromotionOverlay BoardPromotionOverlay
}

func (b *ModelBoard) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if b.elm == nil {
		// create main board element
		newElm = js.Global.Get("document").Call("createElement", "div")
		newElm.Set("id", "board")

		b.Edgings.TopLeft.Position = "top-left"
		b.Edgings.Top.Position = "top"
		b.Edgings.TopRight.Position = "top-right"
		b.Edgings.Left.Position = "left"
		b.Edgings.Right.Position = "right"
		b.Edgings.BottomLeft.Position = "bottom-left"
		b.Edgings.Bottom.Position = "bottom"
		b.Edgings.BottomRight.Position = "bottom-right"

		b.elm = newElm
	}
	// update main board element

	updaters := []Updater{
		&b.Edgings.TopLeft, &b.Edgings.Top, &b.Edgings.TopRight,
		&b.Edgings.Left, &b.Grid, &b.Edgings.Right,
		&b.Edgings.BottomLeft, &b.Edgings.Bottom, &b.Edgings.BottomRight,
		&b.PromotionOverlay,
	}

	for _, updater := range updaters {
		if created, err := updater.Update(events); err != nil {
			return nil, err
		} else if created != nil {
			b.elm.Call("appendChild", created)
		}
	}

	return newElm, nil
}

type ThrownOutsContainer struct {
	elm              *js.Object
	Color            piece.Color
	PieceCount       map[piece.Type]int
	LastMoveThrowOut piece.Type
}

func (c *ThrownOutsContainer) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if c.elm == nil {
		c.elm = js.Global.Get("document").Call("createElement", "div")
		c.elm.Set("id", "thrown-outs-"+strings.ToLower(c.Color.String()))
		c.elm.Get("classList").Call("add", "thrown-outs")

		newElm = c.elm
	}

	c.elm.Set("innerHTML", "")
	for _, pieceType := range thrownOutPiecesOrderType {
		div := js.Global.Get("document").Call("createElement", "div")
		div.Get("classList").Call("add", "piececount")
		if c.LastMoveThrowOut == pieceType {
			div.Get("classList").Call("add", "last-move")
		}
		if c.PieceCount[pieceType] == 0 {
			div.Get("classList").Call("add", "hidden")
		}

		div.Call("appendChild", pieceElement(piece.New(c.Color, pieceType)))

		span := js.Global.Get("document").Call("createElement", "span")
		span.Get("classList").Call("add", "count")
		span.Set("textContent", strconv.Itoa(c.PieceCount[pieceType]))
		div.Call("appendChild", span)

		c.elm.Call("appendChild", div)
	}

	return newElm, nil
}

type ModelThrownouts struct {
	elm          *js.Object
	White, Black ThrownOutsContainer
}

func (t *ModelThrownouts) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object
	if t.elm == nil {
		t.elm = js.Global.Get("document").Call("createElement", "div")
		t.elm.Set("id", "thrown-outs-container")

		t.White.Color = piece.White
		t.Black.Color = piece.Black

		newElm = t.elm
	}

	for _, updater := range []Updater{&t.White, &t.Black} {
		if created, err := updater.Update(events); err != nil {
			return newElm, err
		} else if created != nil {
			t.elm.Call("appendChild", created)
		}
	}

	return newElm, nil
}

type ModelNextMove struct {
	elm          *js.Object
	NextMoveHash string
}

func (nm *ModelNextMove) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object

	if nm.elm == nil {

	}

	return newElm, nil
}

type StatusText struct {
	elm  *js.Object
	Text string
}

func (st *StatusText) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object

	if st.elm == nil {
		st.elm = js.Global.Get("document").Call("createElement", "p")
		st.elm.Set("id", "game-status-text")

		newElm = st.elm
	}

	st.elm.Set("textContent", st.Text)

	return newElm, nil
}

type StatusIcon struct {
	elm          *js.Object
	White, Black bool
}

func (si *StatusIcon) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object

	if si.elm == nil {
		si.elm = js.Global.Get("document").Call("createElement", "p")
		si.elm.Set("id", "game-status-player")

		newElm = si.elm
	}

	si.elm.Set("innerHTML", "")
	if si.White {
		si.elm.Call("appendChild", pieceElement(piece.New(piece.White, piece.King)))
	}
	if si.Black {
		si.elm.Call("appendChild", pieceElement(piece.New(piece.Black, piece.King)))
	}

	return newElm, nil
}

type ModelGameStatus struct {
	elm     *js.Object
	Message StatusText
	Icons   StatusIcon
}

func (gs *ModelGameStatus) Update(events *AppEvents) (*js.Object, error) {
	var newElm *js.Object

	if gs.elm == nil {
		gs.elm = js.Global.Get("document").Call("createElement", "div")
		gs.elm.Set("id", "game-status")

		newElm = gs.elm
	}

	for _, updater := range []Updater{&gs.Message, &gs.Icons} {
		if created, err := updater.Update(events); err != nil {
			return newElm, err
		} else if created != nil {
			gs.elm.Call("appendChild", created)
		}
	}
	return newElm, nil
}

type ModelNotification struct {
	elm   *js.Object
	Shown bool
	Text  string
}
type HtmlModel struct {
	Rotated180deg bool

	Board        ModelBoard
	ThrownOuts   ModelThrownouts
	NextMove     ModelNextMove
	GameStatus   ModelGameStatus
	Notification ModelNotification
}

func (m *HtmlModel) Update(events *AppEvents) ([]*js.Object, error) {

	newElms := []*js.Object{}

	updaters := []Updater{
		&m.Board, &m.ThrownOuts,
		//m.NextMove,
		&m.GameStatus,
		//m.Notification,
	}

	for _, updater := range updaters {
		if created, err := updater.Update(events); err != nil {
			return newElms, err
		} else if created != nil {
			newElms = append(newElms, created)
		}
	}

	if m.Rotated180deg {
		m.Board.elm.Get("classList").Call("add", "rotated180deg")
		m.ThrownOuts.elm.Get("classList").Call("add", "rotated180deg")
	} else {
		m.Board.elm.Get("classList").Call("remove", "rotated180deg")
		m.ThrownOuts.elm.Get("classList").Call("remove", "rotated180deg")
	}

	return newElms, nil
}

type ChessGame struct {
	game       game.Game
	gameGc     GameThrownOuts
	currMoveNo int
	nextMove   move.Move
}

// Creates new chess game from moves string.
// The moves string is basicaly move coordinates from & to (0...63) encoded in base64 (with some improvements for promotions, etc...). See encoding.go
func NewGame(movesString string) (ChessGame, error) {
	moves, err := DecodeMoves(movesString) // encoding.go
	if err != nil {
		return ChessGame{}, errors.New("decoding moves error: " + err.Error())
	}

	g := *game.New()
	//TODO move thown out pieces to game
	gtos := make(GameThrownOuts, len(moves))

	{ // apply game moves
		for i, move := range moves {
			if g.Status() != game.InProgress {
				return ChessGame{}, errors.New("Too many moves in url string! " + strconv.Itoa(i+1) + " moves are enough")
			}

			// position before move
			pbm := g.Position()

			_, merr := g.MakeMove(move)
			if merr != nil {
				return ChessGame{}, errors.New("Errorneous move number " + strconv.Itoa(i+1) + ": " + merr.Error())
			}

			// throw outs for this move
			tos := ThrownOuts{}

			// copy previous move throw outs
			if i > 0 {
				for p, c := range gtos[i-1] {
					tos[p] = c
				}
			}

			// was a piece thrown out regulary? = move destination contains some piece
			if p := pbm.OnSquare(move.To()); p.Type != piece.None {
				if _, ok := tos[p]; !ok {
					tos[p] = 0
				}
				tos[p]++
			}

			// was en passant throw out? = moved piece is pawn and move destination is an en passan square in previous move
			if mp := pbm.OnSquare(move.From()); mp.Type == piece.Pawn && move.To() == pbm.EnPassant {
				p := piece.New(piece.Colors[(pbm.ActiveColor+1)%2], piece.Pawn)
				if _, ok := tos[p]; !ok {
					tos[p] = 0
				}
				tos[p]++
			}

			gtos[i] = tos
		}

		// prepend one empty throw outs
		gtos = append(GameThrownOuts{ThrownOuts{}}, gtos...)
	}

	return ChessGame{g, gtos, len(gtos) - 1, move.Null}, nil
}

func (ch ChessGame) UpdateModel(m *HtmlModel) error {

	if ch.currMoveNo < 0 || ch.currMoveNo >= len(ch.game.Positions) {
		return errors.New("curren move number is out of bounds")
	}

	position := ch.game.Positions[ch.currMoveNo]

	{ // update board pieces
		for i := int(63); i >= 0; i-- {
			m.Board.Grid.Squares[i].Piece = position.OnSquare(square.Square(i))
		}
	}

	{ // update status

		m.GameStatus.Icons.White = false
		m.GameStatus.Icons.Black = false
		if st := ch.game.Status(); st != game.InProgress {
			// game ended
			m.GameStatus.Message.Text = st.String()
			if st&game.Draw != 0 {
				// game ended in draw
				m.GameStatus.Icons.White = true
				m.GameStatus.Icons.Black = true
			} else if st&game.WhiteWon != 0 {
				// white wins
				m.GameStatus.Icons.White = true
			} else if st&game.BlackWon != 0 {
				// black wins
				m.GameStatus.Icons.Black = true
			}
		} else {
			// game in progress
			m.GameStatus.Message.Text = "Moving player"
			if position.ActiveColor == piece.White {
				// white moves
				m.GameStatus.Icons.White = true
			} else {
				// black moves
				m.GameStatus.Icons.Black = true
			}
		}
	}

	return nil
}

type EventFunc func(g *ChessGame, m *HtmlModel) error
type jsEventFunc func(*js.Object)

type AppEvents struct {
	app        *HtmlApp
	registered map[string]map[Elementer]jsEventFunc
}

func (e *AppEvents) Click(elm *js.Object, ef EventFunc) error {
	if e.app == nil {
		return errors.New("no app")
	}
	if elm == nil {
		return errors.New("no element provided for click event")
	}

	jsEventName := "click"

	//if e.registered[jsEventName][elm] != nil {
	//	elm.Element().Call("removeEventListener", jsEventName, e.registered[jsEventName][elm], false)
	//}

	jsEventCallback := func(event *js.Object) {
		ef(&e.app.Game, &e.app.Model)
		if err := e.app.UpdateDom(); err != nil {
			js.Global.Call("alert", "after "+jsEventName+" event app dom update error: "+err.Error())
		}
	}

	elm.Call("addEventListener", jsEventName, jsEventCallback, false)
	js.Global.Call("alert", "registered event: "+elm.String()+":"+elm.Get("id").String())
	//if e.registered == nil {
	//	e.registered = map[string]map[Elementer]jsEventFunc{}
	//}
	//if e.registered[jsEventName] == nil {
	//	e.registered[jsEventName] = map[Elementer]jsEventFunc{}
	//}
	//e.registered[jsEventName][elm] = jsEventCallback
	return nil
}

type HtmlApp struct {
	Game  ChessGame
	Model HtmlModel

	initialized       bool
	rootElement       *js.Object
	events            *AppEvents
	rotationSupported bool
}

func (app *HtmlApp) SetRootElement(elm *js.Object) {
	//TODO if nil, delete??
	app.rootElement = elm
}

func (app *HtmlApp) Init() error {
	{ // initialize game
		game, err := NewGame("")
		if err != nil {
			return err
		}
		app.Game = game
	}

	app.events = &AppEvents{app: app}

	// is rotation supported?
	if div := js.Global.Get("document").Call("createElement", "div"); div != js.Undefined {
		if div.Get("style").Get("transform") != js.Undefined {
			app.rotationSupported = true
		}
		div.Call("remove")
	}

	app.initialized = true
	return nil
}

func (app *HtmlApp) UpdateDom() error {
	if !app.initialized {
		if err := app.Init(); err != nil {
			return err
		}
	}

	{ // update html model from game
		err := app.Game.UpdateModel(&app.Model)
		if err != nil {
			return err
		}
	}

	if !app.rotationSupported {
		app.Model.Rotated180deg = false
		app.Model.Board.Edgings.BottomLeft.Disabled = true
		app.Model.Board.Edgings.TopRight.Disabled = true
	}

	{ // update html dom from html model
		created, err := app.Model.Update(app.events)
		if err != nil {
			return err
		} else if len(created) > 0 {
			if app.rootElement == nil {
				return errors.New("no application element")
			}
			for _, ce := range created {
				app.rootElement.Call("appendChild", ce)
			}
		}
	}
	return nil
}

func (app *HtmlApp) RotateBoardForPlayer() {
	if !app.rotationSupported {
		return
	}
	app.Model.Rotated180deg = false
	if app.Game.game.ActiveColor() == piece.Black {
		app.Model.Rotated180deg = true
	}
}

/*
// Draws chess board, game-status, next-move and notification elements to document.
// Also sets event listeners for grid & copy, make, undo next move and notification
func (app *HtmlApp) DrawBoard() {
	document := js.Global.Get("document")


	// is execCommand supported?
	canExec := false
	if exec := js.Global.Get("document").Get("execCommand"); exec != nil && exec != js.Undefined {
		canExec = true
	}



	{ // next move
		copyOrHint := `<span class="hint">copy this link and send it to your oponent</span>`
		if canExec {
			copyOrHint = `<a href="" id="next-move-copy">Copy to clipboard</a>`
		}

		document.Call("write", `<div id="next-move" class="hidden">
	<p class="link">
		Next move URL link:
		<input id="next-move-input" readonly="readonly" value=""/>
		`+copyOrHint+`
	</p>
  <p class="actions">
		<a id="next-move-link" href="">Make move</a>
		<a id="next-move-back" href="">Undo move</a>
	</p>
</div>`)
		{ // listeners
			if canExec {
				if copy := document.Call("getElementById", "next-move-copy"); copy != nil { // next move copy
					copy.Call(
						"addEventListener",
						"click",
						func(event *js.Object) {
							event.Call("preventDefault")
							// select input text & copy to clipboard
							if nextMoveInput := js.Global.Get("document").Call("getElementById", "next-move-input"); nextMoveInput != nil {
								nextMoveInput.Call("focus")
								nextMoveInput.Call("setSelectionRange", 0, nextMoveInput.Get("value").Get("length"))
								js.Global.Get("document").Call("execCommand", "Copy")
								nextMoveInput.Call("blur")

								// notification
								if notificationOverlay := js.Global.Get("document").Call("getElementById", "notification-overlay"); notificationOverlay != nil {
									// change message
									if notificationMessage := js.Global.Get("document").Call("getElementById", "notification-message"); notificationMessage != nil {
										notificationMessage.Set("innerHTML", "Link has been copied to clipboard, send it to your oponent.")
									}
									// show notification
									notificationOverlay.Get("classList").Call("remove", "hidden")
								}
							}
						},
						false,
					)
				}
			}
			if back := document.Call("getElementById", "next-move-back"); back != nil { // next move back
				back.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						event.Call("preventDefault")
						app.nextMove = move.Null
						//js.Global.Call("alert", "calling update board from next-move-back listener")
						if err := app.UpdateBoard(); err != nil {
							document.Call("getElementById", "game-status").Set("innerHTML", err.Error())
							return
						}
					},
					false,
				)
			}

			if link := document.Call("getElementById", "next-move-link"); link != nil { // next move link
				link.Call(
					"addEventListener",
					"click",
					func(event *js.Object) {
						if nm := document.Call("getElementById", "next-move"); nm != nil {
							nm.Get("classList").Call("add", "hidden")
						}

						if !rotationSupported {
							return
						}

						if board := document.Call("getElementById", "board"); board != nil {
							// rotate only if needed
							shouldBeRotatedInNextMove := app.game.Position().ActiveColor != piece.Black
							isRotated := board.Get("classList").Call("contains", "rotated180deg").Bool()

							if shouldBeRotatedInNextMove != isRotated {
								event.Call("preventDefault")

								//rotate, wait, change location
								rotateBoard180degFunc(event)

								url := link.Call("getAttribute", "href")
								js.Global.Call(
									"setTimeout",
									func() {
										js.Global.Get("location").Set("href", url)
									},
									775,
								)
							}
						}
					},
					false,
				)
			}
		}
	}


	{ // notification overlay
		document.Call("write", `<div id="notification-overlay" class="hidden">
		<p id="notification-message" class="message">long live this notification</p>
		<p class="hint">Click or tap anywhere to close</p>
</div>`)

		// listeners
		if notificationOverlay := document.Call("getElementById", "notification-overlay"); notificationOverlay != nil { // notification overlay
			notificationOverlay.Set("hidden", true)
			notificationOverlay.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					notificationOverlay.Get("classList").Call("add", "hidden")

					// reset message
					if notificationMessage := js.Global.Get("document").Call("getElementById", "notification-message"); notificationMessage != nil {
						notificationMessage.Set("innerHTML", ". . .")
					}
				},
				false,
			)
		}
	}

	{ // event listeners

		for i := int(63); i >= 0; i-- { // map click events to grid squares
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

		if promotionOverlay := document.Call("getElementById", "promotion-overlay"); promotionOverlay != nil { // promotion overlay
			promotionOverlay.Call(
				"addEventListener",
				"click",
				func(event *js.Object) {
					event.Call("preventDefault")
					app.nextMove = move.Null
					//js.Global.Call("alert", "calling update board from promotion-overlay listener")
					if err := app.UpdateBoard(); err != nil {
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
									if err := app.UpdateBoard(); err != nil {
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
func (app *HtmlApp) UpdateBoard() error {
	//js.Global.Call("alert", "update: nextMove: "+app.nextMove.String())

	drawPosition := app.game.Position()
	drawThrownOuts := ThrownOuts{}
	if gcl := len(app.gameGc); gcl > 0 {
		for p, n := range app.gameGc[gcl-1] {
			drawThrownOuts[p] = n
		}
	}

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
					if err := app.UpdateBoard(); err != nil {
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

			{ //check for thrown out piece and add to container if true

				// regular throw out
				if tsp := app.game.Position().OnSquare(app.nextMove.To()); tsp.Type != piece.None {
					if _, ok := drawThrownOuts[tsp]; !ok {
						drawThrownOuts[tsp] = 0
					}
					drawThrownOuts[tsp]++
				}

				// en passant throw out
				if mp := app.game.Position().OnSquare(app.nextMove.From()); mp.Type == piece.Pawn && app.nextMove.To() == app.game.EnPassant() {
					tsp := piece.New(piece.Colors[(app.game.ActiveColor()+1)%2], piece.Pawn)
					if _, ok := drawThrownOuts[tsp]; !ok {
						drawThrownOuts[tsp] = 0
					}
					drawThrownOuts[tsp]++
				}
			}

			// update next-move properties
			nextMoveString, err := encodeMove(app.nextMove) // encoding.go
			if err != nil {
				nextMoveError = errors.New("Next move encoding error: " + err.Error())
			} else {
				location := js.Global.Get("location")
				url := location.Get("origin").String() + location.Get("pathname").String() + "#" + app.movesString + nextMoveString

				if nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link"); nextMoveLink != nil {
					nextMoveLink.Set("href", url)
				}
				if nextMoveInput := js.Global.Get("document").Call("getElementById", "next-move-input"); nextMoveInput != nil {
					nextMoveInput.Set("value", url)
				}

				// show next-move layer
				if nextMove := js.Global.Get("document").Call("getElementById", "next-move"); nextMove != nil {
					nextMove.Get("classList").Call("remove", "hidden")
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
						if err := app.UpdateBoard(); err != nil {
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
								if err := app.UpdateBoard(); err != nil {
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

		pc := drawPosition.OnSquare(sq)
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

			{ // king in check
				if pc.Type == piece.King && drawPosition.Check(pc.Color) {
					markerClasses = append(markerClasses, "check")
				}
			}
		}

		innerHTML := "<span class=\"" + strings.Join(markerClasses, " ") + "\">" + piecesToString[pc] + "</span>"

		sqElm.Set("innerHTML", innerHTML)
	}

	if nextMoveError != nil {
		return nextMoveError
	}

	// fill thrown out pieces
	if len(drawThrownOuts) > 0 {
		for _, c := range piece.Colors {
			id := "thrown-outs-" + strings.ToLower(c.String())
			if tosElm := js.Global.Get("document").Call("getElementById", id); tosElm != nil {
				tos := ThrownOuts{}
				// fill thrown outs only for current color "c"
				for p, n := range drawThrownOuts {
					if p.Color == c {
						tos[p] = n
					}
				}
				tosElmStr := ""

				if len(tos) != 0 {
					// there are thrown out pieces for current color

					// fill black in reversed order
					thrownOutPieces := make([]piece.Type, len(thrownOutPiecesOrderType))
					for i, p := range thrownOutPiecesOrderType {
						j := i
						if c == piece.Black {
							j = len(thrownOutPiecesOrderType) - 1 - i
						}
						thrownOutPieces[j] = p
					}

					// fill blanks if black
					if c == piece.Black {
						bc := len(thrownOutPiecesOrderType) - len(tos)
						for i := 0; i < bc; i++ {
							tosElmStr += "<div class=\"piececount\"></div>"
						}
					}

					lastMoveThrownOutPiece := piece.Piece{}
					//TODO refactor duplicate code
					if drawPosition.Equals(app.game.Position()) {
						if gcl := len(app.game.Positions); gcl-2 >= 0 {
							prevPos := app.game.Positions[gcl-2]
							lastMove := app.game.Position().LastMove

							lastMoveThrownOutPiece = prevPos.OnSquare(lastMove.To())
							// last move with pawn to en passant position
							if mp := prevPos.OnSquare(lastMove.From()); mp.Type == piece.Pawn && prevPos.EnPassant == lastMove.To() {
								lastMoveThrownOutPiece = piece.New(piece.Colors[(prevPos.ActiveColor+1)%2], piece.Pawn)
							}
						}
					} else {
						prevPos := app.game.Position()
						lastMove := app.nextMove

						lastMoveThrownOutPiece = prevPos.OnSquare(lastMove.To())
						// last move with pawn to en passant position
						if mp := prevPos.OnSquare(lastMove.From()); mp.Type == piece.Pawn && prevPos.EnPassant == lastMove.To() {
							lastMoveThrownOutPiece = piece.New(piece.Colors[(prevPos.ActiveColor+1)%2], piece.Pawn)
						}
					}

					// fill pieces
					for _, pt := range thrownOutPieces {
						pc := piece.New(c, pt)
						if n, ok := tos[pc]; ok {
							class := []string{"piececount"}
							if lastMoveThrownOutPiece == pc {
								class = append(class, "last-move")
							}
							tosElmStr += "<div class=\"" + strings.Join(class, " ") + "\">"
							tosElmStr += piecesToString[piece.New(c, pt)]
							tosElmStr += "<span class=\"count\">" + strconv.Itoa(int(n)) + "</span>"
							tosElmStr += "</div>"
						}
					}

					// fill blanks if white
					if c == piece.White {
						bc := len(thrownOutPiecesOrderType) - len(tos)
						for i := 0; i < bc; i++ {
							tosElmStr += "<div class=\"piececount\"></div>"
						}
					}
				}

				tosElm.Set("innerHTML", tosElmStr)
			}
		}
	}


	//js.Global.Call("alert", "end")

	// notification if game not in progress (game ended)
	if st := app.game.Status(); st != game.InProgress {
		if notificationOverlay := js.Global.Get("document").Call("getElementById", "notification-overlay"); notificationOverlay != nil {
			// change message
			if notificationMessage := js.Global.Get("document").Call("getElementById", "notification-message"); notificationMessage != nil {
				notificationMessage.Set("innerHTML", text+".<br />"+`<a href="/">New game</a>?`)
			}
			// show notification
			notificationOverlay.Get("classList").Call("remove", "hidden")
		}
	}

	return nil
}

func (app *HtmlApp) squareHandler(event *js.Object) {
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
		if err := app.UpdateBoard(); err != nil {
			js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
			return
		}
		return
	}

	//js.Global.Call("alert", "handle "+elmId+" square")
	handler(sq, event)
}
*/

type GameThrownOuts []ThrownOuts

type ThrownOuts map[piece.Piece]uint8
