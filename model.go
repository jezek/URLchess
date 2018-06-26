// +build js
package main

import (
	"URLchess/shf"
	"errors"
	"strconv"
	"strings"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"github.com/gopherjs/gopherjs/js"
)

type BoardEdging struct {
	Element  *js.Object
	Position string //top, bottom, left, right, top-left, top-right, bottom-left, bottom-right
}

func (this *BoardEdging) SetPosition(s string) {
	if this == nil {
		return
	}
	this.Position = s

	if this.Element == nil {
		return
	}
	this.Element.Set("id", "edging-"+this.Position)
}

func (this *BoardEdging) init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = js.Global.Get("document").Call("createElement", "div")
		this.Element.Get("classList").Call("add", "edging")
		this.SetPosition(this.Position)
	}
	return nil
}
func (this *BoardEdging) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("BoardEdging is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return nil
}

type EdgingHorizontal struct {
	*BoardEdging
}

func (this *EdgingHorizontal) init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}

		this.Element.Get("classList").Call("add", "horizontal")
		for i := 0; i < 8; i++ {
			letter := js.Global.Get("document").Call("createElement", "div")
			letter.Set("textContent", string(rune('a'+i)))
			this.Element.Call("appendChild", letter)
		}
	}
	return nil
}
func (this *EdgingHorizontal) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingHorizontal is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return tools.Update(this.BoardEdging)
}

type EdgingVertical struct {
	*BoardEdging
}

func (this *EdgingVertical) init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}

		this.Element.Get("classList").Call("add", "vertical")
		for i := 8; i > 0; i-- {
			number := js.Global.Get("document").Call("createElement", "div")
			number.Set("textContent", strconv.Itoa(i))
			this.Element.Call("appendChild", number)
		}
	}
	return nil
}
func (this *EdgingVertical) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingVertical is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return tools.Update(this.BoardEdging)
}

type EdgingCorner struct {
	*BoardEdging
}

func (this *EdgingCorner) init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}

		this.Element.Get("classList").Call("add", "corner")
	}
	return nil
}
func (this *EdgingCorner) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingCorner is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return tools.Update(this.BoardEdging)
}

type EdgingCornerRotating struct {
	*EdgingCorner
	Enabled bool
}

func (this *EdgingCornerRotating) Enable() {
	if this == nil {
		return
	}
	this.Enabled = true

	if this.Element == nil {
		return
	}
	this.Element.Set("innerHTML", "↻")
	this.Element.Get("classList").Call("add", "enabled")
}

func (this *EdgingCornerRotating) Disable() {
	if this == nil {
		return
	}
	this.Enabled = false

	if this.Element == nil {
		return
	}
	this.Element.Set("innerHTML", "")
	this.Element.Get("classList").Call("remove", "enabled")
}

func (this *EdgingCornerRotating) init(tools *shf.Tools) error {
	if this.EdgingCorner == nil {
		this.EdgingCorner = &EdgingCorner{}
		if err := tools.Update(this.EdgingCorner); err != nil {
			return err
		}

		if this.Enabled {
			this.Enable()
		} else {
			this.Disable()
		}
	}
	return nil
}
func (this *EdgingCornerRotating) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingCornerRotating is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return tools.Update(this.EdgingCorner)
}

type GridSquare struct {
	Id      square.Square
	Element *js.Object
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

func (s *GridSquare) init(tools *shf.Tools) error {
	if s.Element == nil {
		boardGridSquareTones := []string{"light-square", "dark-square"}
		s.Element = js.Global.Get("document").Call("createElement", "div")
		s.Element.Set("id", s.Id.String())
		s.Element.Get("classList").Call("add", boardGridSquareTones[(int(s.Id)%8+int(s.Id)/8)%2])
	}
	return nil
}
func (s *GridSquare) Update(tools *shf.Tools) error {
	if s == nil {
		return errors.New("GridSquare is nil")
	}
	if err := s.init(tools); err != nil {
		return err
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

	s.Element.Set("innerHTML", "")
	s.Element.Call("appendChild", marker)

	// events
	//if s.Piece.Type == piece.King {
	//	events.Click(s.Element, func(g ChessGame, m HtmlModel) (ChessGame, HtmlModel) {
	//		js.Global.Call("alert", "kingclick")
	//		return g, m
	//	})
	//}

	return nil
}

type BoardGrid struct {
	Element *js.Object
	Squares [64]*GridSquare
}

func (g *BoardGrid) init(tools *shf.Tools) error {
	for i, sq := range g.Squares {
		if sq == nil {
			g.Squares[i] = &GridSquare{
				Id: square.Square(i),
			}
			if err := tools.Update(g.Squares[i]); err != nil {
				return err
			}
		}
	}

	if g.Element == nil {
		g.Element = js.Global.Get("document").Call("createElement", "div")
		g.Element.Get("classList").Call("add", "grid")
		for i := int(63); i >= 0; i-- {
			g.Element.Call("appendChild", g.Squares[i].Element)
		}
	}
	return nil
}
func (g *BoardGrid) Update(tools *shf.Tools) error {
	if g == nil {
		return errors.New("BoardGrid is nil")
	}
	if err := g.init(tools); err != nil {
		return err
	}

	for i := int(63); i >= 0; i-- {
		if err := tools.Update(g.Squares[i]); err != nil {
			return err
		}
	}
	return nil
}

type BoardPromotionOverlay struct {
	Element *js.Object
	Shown   bool
}

func (p *BoardPromotionOverlay) init(tools *shf.Tools) error {
	if p.Element == nil {
		p.Element = js.Global.Get("document").Call("createElement", "div")
		p.Element.Set("id", "promotion-overlay")

		for _, pieceType := range promotablePiecesType {
			span := js.Global.Get("document").Call("createElement", "span")
			span.Set("id", "promote-to-"+pieceTypesToName[pieceType])
			span.Get("classList").Call("add", "piece")
			span.Set("piece", pieceTypesToName[pieceType])

			p.Element.Call("appendChild", span)
		}
	}
	return nil
}
func (p *BoardPromotionOverlay) Update(tools *shf.Tools) error {
	if p == nil {
		return errors.New("BoardPromotionOverlay is nil")
	}
	if err := p.init(tools); err != nil {
		return err
	}

	if p.Shown {
		p.Element.Get("classList").Call("add", "show")
	} else {
		p.Element.Get("classList").Call("remove", "show")
	}

	return nil
}

type ModelBoard struct {
	Element *js.Object
	Edgings struct {
		Top, Bottom          *EdgingHorizontal
		Left, Right          *EdgingVertical
		TopRight, BottomLeft *EdgingCornerRotating
		TopLeft, BottomRight *EdgingCorner
	}
	Grid             *BoardGrid
	PromotionOverlay *BoardPromotionOverlay
}

func (b *ModelBoard) init(tools *shf.Tools) error {
	if b.Edgings.TopLeft == nil {
		b.Edgings.TopLeft = &EdgingCorner{}
		if err := tools.Update(b.Edgings.TopLeft); err != nil {
			return err
		}

		b.Edgings.TopLeft.SetPosition("top-left")
	}
	if b.Edgings.Top == nil {
		b.Edgings.Top = &EdgingHorizontal{}
		if err := tools.Update(b.Edgings.Top); err != nil {
			return err
		}

		b.Edgings.Top.SetPosition("top")

	}
	if b.Edgings.TopRight == nil {
		b.Edgings.TopRight = &EdgingCornerRotating{}
		if err := tools.Update(b.Edgings.TopRight); err != nil {
			return err
		}
		b.Edgings.TopRight.SetPosition("top-right")
	}

	if b.Edgings.Left == nil {
		b.Edgings.Left = &EdgingVertical{}
		if err := tools.Update(b.Edgings.Left); err != nil {
			return err
		}
		b.Edgings.Left.SetPosition("left")
	}

	if b.Grid == nil {
		b.Grid = &BoardGrid{}
		if err := tools.Update(b.Grid); err != nil {
			return err
		}
	}

	if b.Edgings.Right == nil {
		b.Edgings.Right = &EdgingVertical{}
		if err := tools.Update(b.Edgings.Right); err != nil {
			return err
		}
		b.Edgings.Right.SetPosition("right")
	}

	if b.Edgings.BottomLeft == nil {
		b.Edgings.BottomLeft = &EdgingCornerRotating{}
		if err := tools.Update(b.Edgings.BottomLeft); err != nil {
			return err
		}
		b.Edgings.BottomLeft.SetPosition("bottom-left")
	}
	if b.Edgings.Bottom == nil {
		b.Edgings.Bottom = &EdgingHorizontal{}
		if err := tools.Update(b.Edgings.Bottom); err != nil {
			return err
		}
		b.Edgings.Bottom.SetPosition("bottom")
	}
	if b.Edgings.BottomRight == nil {
		b.Edgings.BottomRight = &EdgingCorner{}
		if err := tools.Update(b.Edgings.BottomRight); err != nil {
			return err
		}
		b.Edgings.BottomRight.SetPosition("bottom-right")
	}

	if b.PromotionOverlay == nil {
		b.PromotionOverlay = &BoardPromotionOverlay{}
		if err := tools.Update(b.PromotionOverlay); err != nil {
			return err
		}
	}

	if b.Element == nil {
		// create main board element
		b.Element = js.Global.Get("document").Call("createElement", "div")
		b.Element.Set("id", "board")

		b.Element.Call("appendChild", b.Edgings.TopLeft.Element)
		b.Element.Call("appendChild", b.Edgings.Top.Element)
		b.Element.Call("appendChild", b.Edgings.TopRight.Element)

		b.Element.Call("appendChild", b.Edgings.Left.Element)
		b.Element.Call("appendChild", b.Grid.Element)
		b.Element.Call("appendChild", b.Edgings.Right.Element)

		b.Element.Call("appendChild", b.Edgings.BottomLeft.Element)
		b.Element.Call("appendChild", b.Edgings.Bottom.Element)
		b.Element.Call("appendChild", b.Edgings.BottomRight.Element)

		b.Element.Call("appendChild", b.PromotionOverlay.Element)
	}
	return nil
}
func (b *ModelBoard) Update(tools *shf.Tools) error {
	if b == nil {
		return errors.New("ModelBoard is nil")
	}
	if err := b.init(tools); err != nil {
		return err
	}

	return tools.Update(
		b.Edgings.TopLeft, b.Edgings.Top, b.Edgings.TopRight,
		b.Edgings.Left, b.Grid, b.Edgings.Right,
		b.Edgings.BottomLeft, b.Edgings.Bottom, b.Edgings.BottomRight,
		b.PromotionOverlay,
	)
}

type ThrownOutsContainer struct {
	Element          *js.Object
	Color            piece.Color
	PieceCount       map[piece.Type]int
	LastMoveThrowOut piece.Type
}

func (c *ThrownOutsContainer) init(tools *shf.Tools) error {
	if c.Element == nil {
		c.Element = js.Global.Get("document").Call("createElement", "div")
		c.Element.Set("id", "thrown-outs-"+strings.ToLower(c.Color.String()))
		c.Element.Get("classList").Call("add", "thrown-outs")
	}
	return nil
}
func (c *ThrownOutsContainer) Update(tools *shf.Tools) error {
	if c == nil {
		return errors.New("ThrownOutsContainer is nil")
	}
	if err := c.init(tools); err != nil {
		return err
	}

	c.Element.Set("innerHTML", "")
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

		c.Element.Call("appendChild", div)
	}

	return nil
}

type ModelThrownouts struct {
	Element      *js.Object
	White, Black *ThrownOutsContainer
}

func (t *ModelThrownouts) init(tools *shf.Tools) error {
	if t.White == nil {
		t.White = &ThrownOutsContainer{
			Color: piece.White,
		}
		if err := tools.Update(t.White); err != nil {
			return err
		}
	}
	if t.Black == nil {
		t.Black = &ThrownOutsContainer{
			Color: piece.Black,
		}
		if err := tools.Update(t.Black); err != nil {
			return err
		}
	}
	if t.Element == nil {
		t.Element = js.Global.Get("document").Call("createElement", "div")
		t.Element.Set("id", "thrown-outs-container")
		t.Element.Call("appendChild", t.White.Element)
		t.Element.Call("appendChild", t.Black.Element)
	}
	return nil
}
func (t *ModelThrownouts) Update(tools *shf.Tools) error {
	if t == nil {
		return errors.New("ModelThrownouts is nil")
	}
	if err := t.init(tools); err != nil {
		return err
	}
	return tools.Update(t.White, t.Black)
}

type StatusText struct {
	Element *js.Object
	Text    string
}

func (st *StatusText) init(tools *shf.Tools) error {
	if st.Element == nil {
		st.Element = js.Global.Get("document").Call("createElement", "p")
		st.Element.Set("id", "game-status-text")

	}
	return nil
}
func (st *StatusText) Update(tools *shf.Tools) error {
	if st == nil {
		return errors.New("StatusText is nil")
	}
	if err := st.init(tools); err != nil {
		return err
	}

	st.Element.Set("textContent", st.Text)

	return nil
}

type StatusIcon struct {
	Element      *js.Object
	White, Black bool
}

func (si *StatusIcon) init(tools *shf.Tools) error {
	if si.Element == nil {
		si.Element = js.Global.Get("document").Call("createElement", "p")
		si.Element.Set("id", "game-status-player")
	}
	return nil
}
func (si *StatusIcon) Update(tools *shf.Tools) error {
	if si == nil {
		return errors.New("StatusIcon is nil")
	}
	if err := si.init(tools); err != nil {
		return err
	}

	si.Element.Set("innerHTML", "")
	if si.White {
		si.Element.Call("appendChild", pieceElement(piece.New(piece.White, piece.King)))
	}
	if si.Black {
		si.Element.Call("appendChild", pieceElement(piece.New(piece.Black, piece.King)))
	}

	return nil
}

type ModelGameStatus struct {
	Element *js.Object
	Message *StatusText
	Icons   *StatusIcon
}

func (gs *ModelGameStatus) init(tools *shf.Tools) error {
	if gs.Message == nil {
		gs.Message = &StatusText{}
		if err := tools.Update(gs.Message); err != nil {
			return err
		}
	}
	if gs.Icons == nil {
		gs.Icons = &StatusIcon{}
		if err := tools.Update(gs.Icons); err != nil {
			return err
		}
	}

	if gs.Element == nil {
		gs.Element = js.Global.Get("document").Call("createElement", "div")
		gs.Element.Set("id", "game-status")
		gs.Element.Call("appendChild", gs.Message.Element)
		gs.Element.Call("appendChild", gs.Icons.Element)
	}
	return nil
}
func (gs *ModelGameStatus) Update(tools *shf.Tools) error {
	if gs == nil {
		return errors.New("ModelGameStatus is nil")
	}
	if err := gs.init(tools); err != nil {
		return err
	}

	return tools.Update(gs.Message, gs.Icons)
}

type ModelMoveStatus struct {
	Element  *js.Object
	Shown    bool
	NextMove bool
	MoveHash string
	MoveCopy bool
}

func (ms *ModelMoveStatus) init(tools *shf.Tools) error {
	if ms.Element == nil {
		ms.Element = js.Global.Get("document").Call("createElement", "div")
		ms.Element.Set("id", "next-move")
	}
	return nil
}

func (ms *ModelMoveStatus) Update(tools *shf.Tools) error {
	if ms == nil {
		return errors.New("ModelMoveStatus is nil")
	}
	if err := ms.init(tools); err != nil {
		return err
	}

	if ms.Shown {
		// update only if shown
		ms.Element.Set("innerHTML", "")

		{ // link paragraph
			link := js.Global.Get("document").Call("createElement", "p")
			link.Get("classList").Call("add", "link")

			{ // link text
				text := "Last move "
				if ms.NextMove {
					text = "NextMove "
				}
				text += "URL"
				link.Call("appendChild", js.Global.Get("document").Call("createTextNode", text))
			}
			{ // link input
				input := js.Global.Get("document").Call("createElement", "input")
				input.Set("id", "next-move-input")
				input.Set("readonly", "readonly")
				input.Set("value", ms.MoveHash)

				link.Call("appendChild", input)
				if ms.MoveCopy { // copy link
					copy := js.Global.Get("document").Call("createElement", "a")
					copy.Set("href", "")
					copy.Set("textContent", "Copy to clipboard")

					//events.Click(copy, func(_ *ChessGame, _ *HtmlModel) error {
					//	input.Call("focus")
					//	input.Call("setSelectionRange", 0, input.Get("value").Get("length"))
					//	js.Global.Get("document").Call("execCommand", "Copy")
					//	input.Call("blur")

					//	//TODO notification
					//	return nil
					//})
					link.Call("appendChild", copy)
				}
			}
			if ms.NextMove { // hint
				hint := js.Global.Get("document").Call("createElement", "span")
				hint.Get("classList").Call("add", "hint")
				hint.Set("textContent", "copy this link and send it to your oponent")

				link.Call("appendChild", hint)
			}

			ms.Element.Call("appendChild", link)
		}

	}

	if ms.Shown {
		ms.Element.Get("classList").Call("remove", "hidden")
	} else {
		ms.Element.Get("classList").Call("add", "hidden")
	}

	return nil
}

type ModelNotification struct {
	Element *js.Object
	Shown   bool
	Text    string
}

func (n *ModelNotification) init(tools *shf.Tools) error {
	if n.Element == nil {
		n.Element = js.Global.Get("document").Call("createElement", "div")
		n.Element.Set("id", "notification-overlay")
	}
	return nil
}
func (n *ModelNotification) Update(tools *shf.Tools) error {
	if n == nil {
		return errors.New("ModelNotification is nil")
	}
	if err := n.init(tools); err != nil {
		return err
	}

	if n.Shown {
		n.Element.Get("classList").Call("remove", "hidden")
	} else {
		n.Element.Get("classList").Call("add", "hidden")
	}
	return nil
}

type HtmlModel struct {
	Rotated180deg bool

	Board        *ModelBoard
	ThrownOuts   *ModelThrownouts
	GameStatus   *ModelGameStatus
	MoveStatus   *ModelMoveStatus
	Notification *ModelNotification
}

func (h *HtmlModel) init(tools *shf.Tools) error {
	if h.Board == nil {
		h.Board = &ModelBoard{}
		if err := tools.Update(h.Board); err != nil {
			return err
		}
	}
	if h.ThrownOuts == nil {
		h.ThrownOuts = &ModelThrownouts{}
		if err := tools.Update(h.ThrownOuts); err != nil {
			return err
		}
	}
	if h.GameStatus == nil {
		h.GameStatus = &ModelGameStatus{}
		if err := tools.Update(h.GameStatus); err != nil {
			return err
		}
	}
	if h.MoveStatus == nil {
		h.MoveStatus = &ModelMoveStatus{}
		if err := tools.Update(h.MoveStatus); err != nil {
			return err
		}
	}
	if h.Notification == nil {
		h.Notification = &ModelNotification{}
		if err := tools.Update(h.Notification); err != nil {
			return err
		}
	}
	return nil
}

func (h *HtmlModel) Update(tools *shf.Tools) error {
	if h == nil {
		return errors.New("HtmlModel is nil")
	}
	if err := h.init(tools); err != nil {
		return err
	}

	if h.Rotated180deg {
		h.Board.Element.Get("classList").Call("add", "rotated180deg")
		h.ThrownOuts.Element.Get("classList").Call("add", "rotated180deg")
	} else {
		h.Board.Element.Get("classList").Call("remove", "rotated180deg")
		h.ThrownOuts.Element.Get("classList").Call("remove", "rotated180deg")
	}

	return tools.Update(h.Board, h.ThrownOuts, h.GameStatus, h.MoveStatus, h.Notification)
}

func (h *HtmlModel) RotateBoard() func(*js.Object) {
	return func(_ *js.Object) {
		h.Rotated180deg = !h.Rotated180deg
	}
}

type ThrownOuts map[piece.Piece]uint8
type GameThrownOuts []ThrownOuts
type ChessGame struct {
	game       *game.Game
	gameGc     GameThrownOuts
	currMoveNo int
	nextMove   move.Move
}

// Creates new chess game from moves string.
// The moves string is basicaly move coordinates from & to (0...63) encoded in base64 (with some improvements for promotions, etc...). See encoding.go
func NewGame(movesString string) (*ChessGame, error) {
	moves, err := DecodeMoves(movesString) // encoding.go
	if err != nil {
		return nil, errors.New("decoding moves error: " + err.Error())
	}

	g := game.New()
	//TODO move thown out pieces to game
	gtos := make(GameThrownOuts, len(moves))

	{ // apply game moves
		for i, move := range moves {
			if g.Status() != game.InProgress {
				return nil, errors.New("Too many moves in url string! " + strconv.Itoa(i+1) + " moves are enough")
			}

			// position before move
			pbm := g.Position()

			_, merr := g.MakeMove(move)
			if merr != nil {
				return nil, errors.New("Errorneous move number " + strconv.Itoa(i+1) + ": " + merr.Error())
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

	return &ChessGame{g, gtos, len(gtos) - 1, move.Null}, nil
}
func (ch *ChessGame) UpdateModel(m *HtmlModel) error {
	if ch == nil {
		return errors.New("ChessGame is nil")
	}

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

type Model struct {
	Game *ChessGame
	Html *HtmlModel

	rotationSupported bool
}

func (m *Model) init(tools *shf.Tools) error {
	if m.Game == nil {
		m.Game, _ = NewGame("")
	}

	if m.Html == nil {
		m.Html = &HtmlModel{}
		if err := tools.Update(m.Html); err != nil {
			return err
		}

		if !m.rotationSupported {
			m.Html.Rotated180deg = false
			m.Html.Board.Edgings.BottomLeft.Disable()
			m.Html.Board.Edgings.TopRight.Disable()
		} else {
			if err := tools.Click(m.Html.Board.Edgings.BottomLeft.Element, m.Html.RotateBoard()); err != nil {
				return err
			}
			m.Html.Board.Edgings.BottomLeft.Enable()
			if err := tools.Click(m.Html.Board.Edgings.TopRight.Element, m.Html.RotateBoard()); err != nil {
				return err
			}
			m.Html.Board.Edgings.TopRight.Enable()
		}
	}
	return nil
}

func (m *Model) Update(tools *shf.Tools) error {
	if m == nil {
		return errors.New("Model is nil")
	}
	if err := m.init(tools); err != nil {
		return err
	}

	{ // update html model from game
		err := m.Game.UpdateModel(m.Html)
		if err != nil {
			return err
		}
	}

	return tools.Update(m.Html)
}
func (app *Model) RotateBoardForPlayer() {
	if !app.rotationSupported {
		return
	}
	app.Html.Rotated180deg = false
	if app.Game.game.ActiveColor() == piece.Black {
		app.Html.Rotated180deg = true
	}
}

/*
// Draws chess board, game-status, next-move and notification elements to document.
// Also sets event listeners for grid & copy, make, undo next move and notification
func (app *Model) DrawBoard() {
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
func (app *Model) UpdateBoard() error {
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

func (app *Model) squareHandler(event *js.Object) {
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
/*
type Model struct {
	Element *js.Object
	Child   *Child
}
func (this *Model) init(tools *app.Tools) error {
	if this.Child == nil {
		this.Child = &Child{}
		if err := tools.Update(this.Child); err != nil {
			return err
		}
	}
	if this.Element == nil {
		this.Element = js.Global.Get("document").Call("createElement", "div")
		this.Element.Set("id", "")
		this.Element.Get("classList").Call("add", "class")

		this.Element.Call("appendChild", this.Child.Element)
	}
	return nil
}
func (this *Model) Update(tools *app.Tools) error {
	if this == nil {
		return errors.New("Model is nil")
	}
	if err := this.init(tools); err != nil {
		return err
	}
	return tools.Update(this.Child)
}
*/
