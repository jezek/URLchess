// +build js
package main

import (
	"URLchess/shf"
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

type BoardEdging struct {
	*shf.Element
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
	this.Set("id", "edging-"+this.Position)
}

func (this *BoardEdging) Init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Get("classList").Call("add", "edging")
		this.SetPosition(this.Position)
	}
	return nil
}
func (this *BoardEdging) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("BoardEdging is nil")
	}
	return nil
}

type EdgingHorizontal struct {
	*BoardEdging
}

func (this *EdgingHorizontal) Init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}
	}

	if this.Element != nil {
		this.Get("classList").Call("add", "horizontal")
		for i := 0; i < 8; i++ {
			letter := tools.CreateElement("div")
			letter.Set("textContent", string(rune('a'+i)))
			this.Call("appendChild", letter.Object())
		}
	}
	return nil
}
func (this *EdgingHorizontal) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingHorizontal is nil")
	}
	return tools.Update(this.BoardEdging)
}

type EdgingVertical struct {
	*BoardEdging
}

func (this *EdgingVertical) Init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}
	}

	if this.Element != nil {
		this.Get("classList").Call("add", "vertical")
		for i := 8; i > 0; i-- {
			number := tools.CreateElement("div")
			number.Set("textContent", strconv.Itoa(i))
			this.Call("appendChild", number.Object())
		}
	}

	return nil
}
func (this *EdgingVertical) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingVertical is nil")
	}
	return tools.Update(this.BoardEdging)
}

type EdgingCorner struct {
	*BoardEdging
}

func (this *EdgingCorner) Init(tools *shf.Tools) error {
	if this.BoardEdging == nil {
		this.BoardEdging = &BoardEdging{}
		if err := tools.Update(this.BoardEdging); err != nil {
			return err
		}
	}

	if this.Element != nil {
		this.Get("classList").Call("add", "corner")
	}
	return nil
}
func (this *EdgingCorner) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("EdgingCorner is nil")
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
	this.Set("innerHTML", "↻")
	this.Get("classList").Call("add", "enabled")
}

func (this *EdgingCornerRotating) Disable() {
	if this == nil {
		return
	}
	this.Enabled = false

	if this.Element == nil {
		return
	}
	this.Set("innerHTML", "")
	this.Get("classList").Call("remove", "enabled")
}
func (this *EdgingCornerRotating) Init(tools *shf.Tools) error {
	if this.EdgingCorner == nil {
		this.EdgingCorner = &EdgingCorner{}
		if err := tools.Update(this.EdgingCorner); err != nil {
			return err
		}
	}

	if this.Element != nil {
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
	return tools.Update(this.EdgingCorner)
}

type MarkersByColor struct {
	LastMove struct {
		From, To bool
	}
	NextMove struct {
		From, To, PossibleTo bool
	}
}

type SquareMarkers struct {
	ByColor [2]MarkersByColor
	Check   bool
}
type GridSquare struct {
	*shf.Element
	Id      square.Square
	Piece   piece.Piece
	Markers SquareMarkers
}

func (s *GridSquare) Init(tools *shf.Tools) error {

	if s.Element == nil {
		boardGridSquareTones := []string{"light-square", "dark-square"}
		s.Element = tools.CreateElement("div")
		s.Set("id", s.Id.String())
		s.Get("classList").Call("add", boardGridSquareTones[(int(s.Id)%8+int(s.Id)/8)%2])
	}
	return nil
}

func (s *GridSquare) Update(tools *shf.Tools) error {
	if s == nil {
		return errors.New("GridSquare is nil")
	}

	// update square, generate content & replace
	marker := tools.CreateElement("span")

	marker.Get("classList").Call("add", "marker")
	for _, color := range piece.Colors {
		colorString := strings.ToLower(color.String())

		if s.Markers.ByColor[color].LastMove.From {
			marker.Get("classList").Call("add", "last-move")
			marker.Get("classList").Call("add", "last-move-"+colorString)
			marker.Get("classList").Call("add", "last-move-from")
		}
		if s.Markers.ByColor[color].LastMove.To {
			marker.Get("classList").Call("add", "last-move")
			marker.Get("classList").Call("add", "last-move-"+colorString)
			marker.Get("classList").Call("add", "last-move-to")
		}
		if s.Markers.ByColor[color].NextMove.From {
			marker.Get("classList").Call("add", "next-move")
			marker.Get("classList").Call("add", "next-move-"+colorString)
			marker.Get("classList").Call("add", "next-move-from")
		}
		if s.Markers.ByColor[color].NextMove.To {
			marker.Get("classList").Call("add", "next-move")
			marker.Get("classList").Call("add", "next-move-"+colorString)
			marker.Get("classList").Call("add", "next-move-to")
		}
		if s.Markers.ByColor[color].NextMove.PossibleTo {
			marker.Get("classList").Call("add", "next-move")
			marker.Get("classList").Call("add", "next-move-"+colorString)
			marker.Get("classList").Call("add", "next-move-possible-to")
		}
	}
	if s.Markers.Check {
		marker.Get("classList").Call("add", "check")
	}

	if s.Piece.Type != piece.None {
		marker.Call("appendChild", pieceElement(s.Piece))
	}

	s.Set("innerHTML", "")
	s.Call("appendChild", marker.Object())

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
	*shf.Element
	Squares [64]*GridSquare
}

func (g *BoardGrid) Init(tools *shf.Tools) error {
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
		g.Element = tools.CreateElement("div")
		g.Get("classList").Call("add", "grid")
		for i := int(63); i >= 0; i-- {
			if g.Squares[i].Element != nil {
				//TODO append/replace to position
				g.Call("appendChild", g.Squares[i].Element.Object())
			}
		}
	}

	return nil
}
func (g *BoardGrid) Update(tools *shf.Tools) error {
	if g == nil {
		return errors.New("BoardGrid is nil")
	}

	for i := int(63); i >= 0; i-- {
		if err := tools.Update(g.Squares[i]); err != nil {
			return err
		}
	}
	return nil
}

type BoardPromotionOverlay struct {
	*shf.Element
	Shown bool
}

func (p *BoardPromotionOverlay) Init(tools *shf.Tools) error {
	if p.Element == nil {
		p.Element = tools.CreateElement("div")
		p.Set("id", "promotion-overlay")

		for _, pieceType := range promotablePiecesType {
			span := tools.CreateElement("span")
			span.Set("id", "promote-to-"+pieceTypesToName[pieceType])
			span.Get("classList").Call("add", "piece")
			span.Set("piece", pieceTypesToName[pieceType])

			p.Call("appendChild", span.Object())
		}
	}
	return nil
}
func (p *BoardPromotionOverlay) Update(tools *shf.Tools) error {
	if p == nil {
		return errors.New("BoardPromotionOverlay is nil")
	}

	if p.Shown {
		p.Get("classList").Call("add", "show")
	} else {
		p.Get("classList").Call("remove", "show")
	}

	return nil
}

type ModelBoard struct {
	*shf.Element
	Edgings struct {
		Top, Bottom          *EdgingHorizontal
		Left, Right          *EdgingVertical
		TopRight, BottomLeft *EdgingCornerRotating
		TopLeft, BottomRight *EdgingCorner
	}
	Grid             *BoardGrid
	PromotionOverlay *BoardPromotionOverlay
}

func (b *ModelBoard) Init(tools *shf.Tools) error {
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
		b.Element = tools.CreateElement("div")
		b.Set("id", "board")

		b.Call("appendChild", b.Edgings.TopLeft.Element.Object())
		b.Call("appendChild", b.Edgings.Top.Element.Object())
		b.Call("appendChild", b.Edgings.TopRight.Element.Object())

		b.Call("appendChild", b.Edgings.Left.Element.Object())
		b.Call("appendChild", b.Grid.Element.Object())
		b.Call("appendChild", b.Edgings.Right.Element.Object())

		b.Call("appendChild", b.Edgings.BottomLeft.Element.Object())
		b.Call("appendChild", b.Edgings.Bottom.Element.Object())
		b.Call("appendChild", b.Edgings.BottomRight.Element.Object())

		b.Call("appendChild", b.PromotionOverlay.Element.Object())
	}
	return nil
}
func (b *ModelBoard) Update(tools *shf.Tools) error {
	if b == nil {
		return errors.New("ModelBoard is nil")
	}

	return tools.Update(
		b.Edgings.TopLeft, b.Edgings.Top, b.Edgings.TopRight,
		b.Edgings.Left, b.Grid, b.Edgings.Right,
		b.Edgings.BottomLeft, b.Edgings.Bottom, b.Edgings.BottomRight,
		b.PromotionOverlay,
	)
}

type ThrownOutsContainer struct {
	*shf.Element
	Color            piece.Color
	PieceCount       map[piece.Type]int
	LastMoveThrowOut piece.Type
}

func (c *ThrownOutsContainer) Init(tools *shf.Tools) error {
	if c.PieceCount == nil {
		c.PieceCount = map[piece.Type]int{}
	}
	if c.Element == nil {
		c.Element = tools.CreateElement("div")
		c.Set("id", "thrown-outs-"+strings.ToLower(c.Color.String()))
		c.Get("classList").Call("add", "thrown-outs")
	}
	return nil
}
func (c *ThrownOutsContainer) Update(tools *shf.Tools) error {
	if c == nil {
		return errors.New("ThrownOutsContainer is nil")
	}

	c.Set("innerHTML", "")
	for _, pieceType := range thrownOutPiecesOrderType {
		div := tools.CreateElement("div")
		div.Get("classList").Call("add", "piececount")
		if c.LastMoveThrowOut == pieceType {
			div.Get("classList").Call("add", "last-move")
		}
		if c.PieceCount[pieceType] == 0 {
			div.Get("classList").Call("add", "hidden")
		}

		div.Call("appendChild", pieceElement(piece.New(c.Color, pieceType)))

		span := tools.CreateElement("span")
		span.Get("classList").Call("add", "count")
		span.Set("textContent", strconv.Itoa(c.PieceCount[pieceType]))
		div.Call("appendChild", span.Object())

		c.Call("appendChild", div.Object())
	}

	return nil
}

type ModelThrownouts struct {
	*shf.Element
	White, Black *ThrownOutsContainer
}

func (t *ModelThrownouts) Init(tools *shf.Tools) error {
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
		t.Element = tools.CreateElement("div")
		t.Set("id", "thrown-outs-container")
		t.Call("appendChild", t.White.Element.Object())
		t.Call("appendChild", t.Black.Element.Object())
	}
	return nil
}
func (t *ModelThrownouts) Update(tools *shf.Tools) error {
	if t == nil {
		return errors.New("ModelThrownouts is nil")
	}
	return tools.Update(t.White, t.Black)
}

type StatusText struct {
	*shf.Element
	Text string
}

func (st *StatusText) Init(tools *shf.Tools) error {
	if st.Element == nil {
		st.Element = tools.CreateElement("p")
		st.Set("id", "game-status-text")

	}
	return nil
}
func (st *StatusText) Update(tools *shf.Tools) error {
	if st == nil {
		return errors.New("StatusText is nil")
	}

	st.Set("textContent", st.Text)

	return nil
}

type StatusIcon struct {
	*shf.Element
	White, Black bool
}

func (si *StatusIcon) Init(tools *shf.Tools) error {
	if si.Element == nil {
		si.Element = tools.CreateElement("p")
		si.Set("id", "game-status-player")
	}
	return nil
}
func (si *StatusIcon) Update(tools *shf.Tools) error {
	if si == nil {
		return errors.New("StatusIcon is nil")
	}

	si.Set("innerHTML", "")
	if si.White {
		si.Call("appendChild", pieceElement(piece.New(piece.White, piece.King)))
	}
	if si.Black {
		si.Call("appendChild", pieceElement(piece.New(piece.Black, piece.King)))
	}

	return nil
}

type ModelGameStatus struct {
	*shf.Element
	Message *StatusText
	Icons   *StatusIcon
}

func (gs *ModelGameStatus) Init(tools *shf.Tools) error {
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
		gs.Element = tools.CreateElement("div")
		gs.Set("id", "game-status")
		gs.Call("appendChild", gs.Message.Element.Object())
		gs.Call("appendChild", gs.Icons.Element.Object())
	}
	return nil
}
func (gs *ModelGameStatus) Update(tools *shf.Tools) error {
	if gs == nil {
		return errors.New("ModelGameStatus is nil")
	}

	return tools.Update(gs.Message, gs.Icons)
}

type ModelMoveStatus struct {
	*shf.Element
	Shown    bool
	NextMove bool
	MoveHash string
	MoveCopy bool
}

func (ms *ModelMoveStatus) Init(tools *shf.Tools) error {
	if ms.Element == nil {
		ms.Element = tools.CreateElement("div")
		ms.Set("id", "next-move")
	}
	return nil
}

func (ms *ModelMoveStatus) Update(tools *shf.Tools) error {
	if ms == nil {
		return errors.New("ModelMoveStatus is nil")
	}

	if ms.Shown {
		// update only if shown
		ms.Set("innerHTML", "")

		{ // link paragraph
			link := tools.CreateElement("p")
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
				input := tools.CreateElement("input")
				input.Set("id", "next-move-input")
				input.Set("readonly", "readonly")
				input.Set("value", ms.MoveHash)

				link.Call("appendChild", input.Object())
				if ms.MoveCopy { // copy link
					copy := tools.CreateElement("a")
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
					link.Call("appendChild", copy.Object())
				}
			}
			if ms.NextMove { // hint
				hint := tools.CreateElement("span")
				hint.Get("classList").Call("add", "hint")
				hint.Set("textContent", "copy this link and send it to your oponent")

				link.Call("appendChild", hint.Object())
			}

			ms.Call("appendChild", link.Object())
		}

	}

	if ms.Shown {
		ms.Get("classList").Call("remove", "hidden")
	} else {
		ms.Get("classList").Call("add", "hidden")
	}

	return nil
}

type ModelNotification struct {
	*shf.Element
	Shown bool
	Text  string
}

func (n *ModelNotification) Init(tools *shf.Tools) error {
	if n.Element == nil {
		n.Element = tools.CreateElement("div")
		n.Set("id", "notification-overlay")
	}
	return nil
}
func (n *ModelNotification) Update(tools *shf.Tools) error {
	if n == nil {
		return errors.New("ModelNotification is nil")
	}

	if n.Shown {
		n.Get("classList").Call("remove", "hidden")
	} else {
		n.Get("classList").Call("add", "hidden")
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

func (h *HtmlModel) Init(tools *shf.Tools) error {
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

	if h.Rotated180deg {
		h.Board.Get("classList").Call("add", "rotated180deg")
		h.ThrownOuts.Get("classList").Call("add", "rotated180deg")
	} else {
		h.Board.Get("classList").Call("remove", "rotated180deg")
		h.ThrownOuts.Get("classList").Call("remove", "rotated180deg")
	}

	return tools.Update(h.Board, h.ThrownOuts, h.GameStatus, h.MoveStatus, h.Notification)
}

func (h *HtmlModel) RotateBoard() func(*shf.Event) error {
	return func(_ *shf.Event) error {
		h.Rotated180deg = !h.Rotated180deg
		return nil
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
func (ch *ChessGame) UpdateModel(tools *shf.Tools, m *HtmlModel) error {
	if ch == nil {
		return errors.New("ChessGame is nil")
	}

	if ch.currMoveNo < 0 || ch.currMoveNo >= len(ch.game.Positions) {
		return errors.New("curren move number is out of bounds")
	}
	if len(ch.game.Positions) != len(ch.gameGc) {
		return errors.New("count of game moves and throuwn outs does not match")
	}

	position := ch.game.Positions[ch.currMoveNo]
	thrownOutPieces := ch.gameGc[ch.currMoveNo]

	{ // validate next move
		if _, err := nextMoveState(position, ch.nextMove); err != nil {
			// this should not happen
			return err
		}
		// next move is valid (a valid move, or waiting to fill some params)
	}

	{ // if next move is a legal move, update drawing position
		//TODO
	}

	{ // update board pieces
		for i := int(63); i >= 0; i-- {
			m.Board.Grid.Squares[i].Piece = position.OnSquare(square.Square(i))

			// reset all square markers
			m.Board.Grid.Squares[i].Markers = SquareMarkers{}

			if m.Board.Grid.Squares[i].Piece.Type == piece.King {
				// piece is a king ... is he in check?
				if position.Check(m.Board.Grid.Squares[i].Piece.Color) {
					m.Board.Grid.Squares[i].Markers.Check = true
				}
			}
		}

		if position.LastMove != move.Null { // last move marker
			m.Board.Grid.Squares[int(position.LastMove.From())].Markers.ByColor[complementColor(position.ActiveColor)].LastMove.From = true
			m.Board.Grid.Squares[int(position.LastMove.To())].Markers.ByColor[complementColor(position.ActiveColor)].LastMove.To = true
		}

		if ch.nextMove.From() != square.NoSquare { // next move from marker
			m.Board.Grid.Squares[int(ch.nextMove.From())].Markers.ByColor[position.ActiveColor].NextMove.From = true
		}
		if ch.nextMove.To() != square.NoSquare { // next move to marker
			m.Board.Grid.Squares[int(ch.nextMove.To())].Markers.ByColor[position.ActiveColor].NextMove.To = true
		}

		if ch.nextMove.From() != square.NoSquare && ch.nextMove.To() == square.NoSquare {
			// fill possible moves
			// mark possible to squares
			for move, _ := range position.LegalMoves() {
				if move.From() != ch.nextMove.From() {
					continue
				}
				m.Board.Grid.Squares[int(move.To())].Markers.ByColor[position.ActiveColor].NextMove.PossibleTo = true
			}
		}
	}

	{ // update thrown out pieces
		for p, count := range thrownOutPieces {
			container := m.ThrownOuts.White
			if p.Color == piece.Black {
				container = m.ThrownOuts.Black
			}
			container.PieceCount[p.Type] = int(count)
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

	{ // update handlers
		{ // board grid squares
			for _, sq := range m.Board.Grid.Squares {
				if position.ActiveColor == sq.Piece.Color {
					// every moving player figure gets event
					if err := tools.Click(sq.Element, ch.NextMoveFrom(sq.Id)); err != nil {
						return err
					}
				}

				if sq.Markers.ByColor[position.ActiveColor].NextMove.PossibleTo {
					// every moving player possible to move gets event
					if err := tools.Click(sq.Element, ch.NextMoveTo(sq.Id)); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (ch *ChessGame) NextMoveFrom(sq square.Square) func(_ *shf.Event) error {
	return func(_ *shf.Event) error {
		ch.nextMove.Source = sq
		ch.nextMove.Destination = square.NoSquare
		ch.nextMove.Promote = piece.None
		return nil
	}
}

func (ch *ChessGame) NextMoveTo(sq square.Square) func(_ *shf.Event) error {
	return func(_ *shf.Event) error {
		ch.nextMove.Destination = sq
		return nil
	}
}

type Model struct {
	Game *ChessGame
	Html *HtmlModel

	rotationSupported bool
}

func (m *Model) Init(tools *shf.Tools) error {
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

	{ // update html model from game
		err := m.Game.UpdateModel(tools, m.Html)
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
						func(event *js.Element) {
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
					func(event *js.Element) {
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
					func(event *js.Element) {
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
				func(event *js.Element) {
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
				func(event *js.Element) {
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
						func(event *js.Element) {
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

	// precalculate next move markers and stuff
	var nextMoveError error
	nextMoveMarkerClasses := map[square.Square][]string{}
	if app.nextMove == move.Null {
		//js.Global.Call("alert", "no next move")

		// set handlers for moving player pieces
		for i := int(63); i >= 0; i-- {
			sq := square.Square(i)
			pc := drawPosition.OnSquare(sq)
			if pc.Color == app.game.Position().ActiveColor {
				app.squaresHandler[sq] = func(sq square.Square, _ *js.Element) {
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
		//js.Global.Call("alert", "some move, legal or illegal or incomplete")

		color := strings.ToLower(app.game.Position().ActiveColor.String())
		if _, ok := app.game.Position().LegalMoves()[app.nextMove]; ok {
			//js.Global.Call("alert", "legal move")

			drawPosition = app.game.Position().MakeMove(app.nextMove)


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
			//js.Global.Call("alert", "illegal move")

			// set handlers for moving player pieces
			for i := int(63); i >= 0; i-- {
				sq := square.Square(i)
				pc := drawPosition.OnSquare(sq)
				if pc.Color == app.game.Position().ActiveColor {
					app.squaresHandler[sq] = func(sq square.Square, _ *js.Element) {
						app.nextMove.Source = sq
						//js.Global.Call("alert", "calling update board from square handler moving piece illegal move")
						if err := app.UpdateBoard(); err != nil {
							js.Global.Get("document").Call("getElementById", "game-status").Set("innerHTML", err.Error())
							return
						}
					}
				}
			}

			if app.nextMove.Source == square.NoSquare {
			} else {
				//js.Global.Call("alert", "from filled")

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
							app.squaresHandler[move.Destination] = func(sq square.Square, _ *js.Element) {
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
				}
			}
		}
	}


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

func (app *Model) squareHandler(event *js.Element) {
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

func isLegalMove(p *position.Position, m move.Move) bool {
	_, ok := p.LegalMoves()[m]
	return ok
}
func isLegalMoveFrom(p *position.Position, f square.Square) bool {
	for move, _ := range p.LegalMoves() {
		if move.Source == f {
			return true
		}
	}
	return false
}
func isLegalMoveFromTo(p *position.Position, f, t square.Square) bool {
	for move, _ := range p.LegalMoves() {
		if move.Source == f && move.Destination == t {
			return true
		}
	}
	return false
}

const (
	NMError = iota
	NMLegalMove
	NMWaitFrom
	NMWaitTo
	NMWaitPromote
)

func nextMoveState(p *position.Position, m move.Move) (int, error) {
	if m == move.Null {
		// no next move
		return NMWaitFrom, nil
	}
	// some move, legal or illegal or incomplete

	if isLegalMove(p, m) {
		// legal move
		return NMLegalMove, nil
	}
	// illegal od incomplete move

	if m.Source == square.NoSquare {
		// from not filled
		return NMError, errors.New("next move is not null, but has no from square filled")
	}
	// from filled

	if !isLegalMoveFrom(p, m.From()) {
		// from is illegal
		if p.OnSquare(m.From()).Color == p.ActiveColor && m.To() == square.NoSquare && m.Promote == piece.None {
			// but if only from is filled & piece on from square is an ctive piece, so let it be valid
			return NMWaitTo, nil
		}
		return NMError, errors.New("next move from square is illegal! from: " + m.From().String())
	}
	// from is legal

	if m.To() == square.NoSquare {
		// to not filled
		if m.Promote != piece.None {
			// fault move, from filled, to not, but promotion yes
			return NMError, errors.New("next move from and promotion is filled, but to not! from: " + m.From().String() + ", to: " + m.To().String() + ", promote: " + m.Promote.String())
		}
		// incomplete move, needs destination, but valid
		return NMWaitTo, nil
	}
	//to filled

	if !isLegalMoveFromTo(p, m.From(), m.To()) {
		// from, to pair is illegal
		return NMError, errors.New("next move to square is illegal! from: " + m.From().String() + ", to: " + m.To().String())
	}
	// from and to is a legal pair, but whole move not, so promotion is missing or illegal

	if m.Promote != piece.None {
		// from, to & promotion filled, frtom & to are valid, but whole move illegal. So promotion is invalid
		return NMError, errors.New("promotion is invalid! from: " + m.From().String() + ", to: " + m.To().String() + ", promote: " + m.Promote.String())
	}

	// from, to filled & valid, promotion is missing
	return NMWaitPromote, nil
}

/*
type Model struct {
	*shf.Element
	Child   *Child
}
func (this *Model) Init(tools *app.Tools) error {
	if this.Child == nil {
		this.Child = &Child{}
		if err := tools.Update(this.Child); err != nil {
			return err
		}
	}
	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Set("id", "")
		this.Get("classList").Call("add", "class")

		if this.Child.Element != nil {
			this.Call("appendChild", this.Child.Element.Object())
		}
	}
	return nil
}
func (this *Model) Update(tools *app.Tools) error {
	if this == nil {
		return errors.New("Model is nil")
	}
	return tools.Update(this.Child)
}
*/
