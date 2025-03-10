package main

import (
	"URLchess/shf"
	"URLchess/shf/js"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/pgn"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

type ModelHeader struct {
	shf.Element
}

func (h *ModelHeader) Init(tools *shf.Tools) error {
	if h.Element == nil {
		h.Element = tools.CreateElement("div")
		h.Set("id", "header")
		h.Set("innerHTML", "URL&#8203;chess")
	}
	return nil
}
func (h *ModelHeader) Update(tools *shf.Tools) error {
	if h == nil {
		return errors.New("ModelHeader is nil")
	}
	return nil
}

type ModelFooter struct {
	shf.Element
}

func (f *ModelFooter) Init(tools *shf.Tools) error {
	if f.Element == nil {
		f.Element = tools.CreateElement("div")
		f.Set("id", "footer")

		linkURLchess := tools.CreateElement("a")
		linkURLchess.Set("href", "https://jezek.github.io/URLchess")
		linkURLchess.Call("appendChild", tools.CreateTextNode("URLchess"))
		if err := tools.Click(linkURLchess, func(_ shf.Event) error {
			return nil
		}); err != nil {
			return err
		}

		linkGit := tools.CreateElement("a")
		linkGit.Set("href", "https://github.com/jezek/URLchess")
		linkGit.Call("appendChild", tools.CreateTextNode("github"))

		engine := " ()"
		if js.WRAPS == "syscall/js" {
			engine = " (wasm)"
		} else if js.WRAPS == "github.com/gopherjs/gopherjs/js" {
			engine = " (js)"
		}

		f.Call("appendChild", linkURLchess.Object())
		f.Call("appendChild", tools.CreateTextNode(" v"+Version+engine+" by jEzEk. Source on "))
		f.Call("appendChild", linkGit.Object())
	}
	return nil
}
func (f *ModelFooter) Update(tools *shf.Tools) error {
	if f == nil {
		return errors.New("ModelFooter is nil")
	}
	return nil
}

type BoardEdging struct {
	shf.Element
	Position string //top, bottom, left, right, top-left, top-right, bottom-left, bottom-right
}

func (this *BoardEdging) SetPosition(s string) {
	if this == nil {
		return
	}
	this.Position = strings.TrimSpace(s)

	if this.Element == nil {
		return
	}

	if len(this.Position) == 0 {
		this.Delete("id")
	} else {
		this.Set("id", "edging-"+this.Position)
	}

	this.Get("classList").Call("add", "edging")

	if len(this.Position) > 0 {
		for _, posClass := range strings.Split(this.Position, "-") {
			this.Get("classList").Call("add", posClass)
		}
	}
}

func (this *BoardEdging) Init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = tools.CreateElement("div")
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
		if err := tools.Initialize(this.BoardEdging); err != nil {
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
		if err := tools.Initialize(this.BoardEdging); err != nil {
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
		if err := tools.Initialize(this.BoardEdging); err != nil {
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
		if err := tools.Initialize(this.EdgingCorner); err != nil {
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
	Mate    bool
}
type GridSquare struct {
	shf.Element
	Id      square.Square
	Piece   piece.Piece
	Markers SquareMarkers
	piece   shf.Element
	marker  shf.Element
}

func (s *GridSquare) Init(tools *shf.Tools) error {
	if s.piece == nil {
		s.piece = pieceElement(tools, piece.New(piece.NoColor, piece.None))
	}
	if s.marker == nil {
		s.marker = tools.CreateElement("span")
		s.marker.Call("appendChild", s.piece.Object())
	}

	if s.Element == nil {
		boardGridSquareTones := []string{"light-square", "dark-square"}
		s.Element = tools.CreateElement("div")
		s.Set("id", s.Id.String())
		s.Get("classList").Call("add", boardGridSquareTones[(int(s.Id)%8+int(s.Id)/8)%2])

		s.Call("appendChild", s.marker.Object())
	}
	return nil
}

func (s *GridSquare) Update(tools *shf.Tools) error {
	if s == nil {
		return errors.New("GridSquare is nil")
	}

	// update square, generate content & replace
	s.marker.Set("className", "marker")

	for _, color := range piece.Colors {
		colorString := strings.ToLower(color.String())

		if s.Markers.ByColor[color].LastMove.From {
			s.marker.Get("classList").Call("add", "last-move")
			s.marker.Get("classList").Call("add", "last-move-"+colorString)
			s.marker.Get("classList").Call("add", "last-move-from")
		}
		if s.Markers.ByColor[color].LastMove.To {
			s.marker.Get("classList").Call("add", "last-move")
			s.marker.Get("classList").Call("add", "last-move-"+colorString)
			s.marker.Get("classList").Call("add", "last-move-to")
		}
		if s.Markers.ByColor[color].NextMove.From {
			s.marker.Get("classList").Call("add", "next-move")
			s.marker.Get("classList").Call("add", "next-move-"+colorString)
			s.marker.Get("classList").Call("add", "next-move-from")
		}
		if s.Markers.ByColor[color].NextMove.To {
			s.marker.Get("classList").Call("add", "next-move")
			s.marker.Get("classList").Call("add", "next-move-"+colorString)
			s.marker.Get("classList").Call("add", "next-move-to")
		}
		if s.Markers.ByColor[color].NextMove.PossibleTo {
			s.marker.Get("classList").Call("add", "next-move")
			s.marker.Get("classList").Call("add", "next-move-"+colorString)
			s.marker.Get("classList").Call("add", "next-move-possible-to")
		}
	}
	if s.Markers.Check {
		if s.Markers.Mate {
			s.marker.Get("classList").Call("add", "check-mate")
		} else {
			s.marker.Get("classList").Call("add", "check")
		}
	}

	s.piece.Set("className", "piece")
	if s.Piece.Color != piece.NoColor {
		s.piece.Get("classList").Call("add", strings.ToLower(s.Piece.Color.String()))
	}
	if s.Piece.Type != piece.None {
		s.piece.Get("classList").Call("add", pieceTypesToName[s.Piece.Type])
	}
	s.piece.Set("textContent", s.Piece.Figurine())

	return nil
}

type BoardGrid struct {
	shf.Element
	Squares [64]*GridSquare
}

func (g *BoardGrid) Init(tools *shf.Tools) error {
	for i, sq := range g.Squares {
		if sq == nil {
			g.Squares[i] = &GridSquare{
				Id: square.Square(i),
			}
			if err := tools.Initialize(g.Squares[i]); err != nil {
				return err
			}
		}
	}

	if g.Element == nil {
		g.Element = tools.CreateElement("div")
		g.Get("classList").Call("add", "grid")
		for i := int(63); i >= 0; i-- {
			if g.Squares[i].Element != nil {
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

type PromotionPiece struct {
	shf.Element
	Piece        piece.Piece
	PieceElement shf.Element
}

func (p *PromotionPiece) RedrawElement(tools *shf.Tools) {
	if p.Element == nil {
		return
	}

	tools.Destroy(p.PieceElement)
	//TODO jezek - Look for other places where to destroy element (innerHTML, ...).
	p.Element.Set("innerHTML", "")

	p.PieceElement = pieceElement(tools, p.Piece)
	p.Element.Call("appendChild", p.PieceElement.Object())
}
func (p *PromotionPiece) Init(tools *shf.Tools) error {
	if p.Element == nil {
		p.Element = tools.CreateElement("span")
		p.Element.Set("id", "promote-to-"+pieceTypesToName[p.Piece.Type])

		//TODO jezek - Test if this can be deleted.
		//p.RedrawElement(tools)
	}
	return nil
}
func (p *PromotionPiece) Update(tools *shf.Tools) error {
	if p == nil {
		return errors.New("PromotionPiece is nil")
	}
	p.RedrawElement(tools)
	return nil
}

type BoardPromotionOverlay struct {
	shf.Element
	Shown  bool
	Color  piece.Color
	Pieces []*PromotionPiece
}

func (p *BoardPromotionOverlay) Init(tools *shf.Tools) error {
	p.Color = piece.NoColor

	if p.Pieces == nil {
		p.Pieces = []*PromotionPiece{}
		for _, pieceType := range promotablePiecesType {
			pp := &PromotionPiece{Piece: piece.New(p.Color, pieceType)}
			if err := tools.Initialize(pp); err != nil {
				return err
			}
			p.Pieces = append(p.Pieces, pp)
		}
	}

	if p.Element == nil {
		p.Element = tools.CreateElement("div")
		p.Set("id", "promotion-overlay")

		for _, piece := range p.Pieces {
			p.Element.Call("appendChild", piece.Object())
		}
	}
	return nil
}
func (p *BoardPromotionOverlay) Update(tools *shf.Tools) error {
	if p == nil {
		return errors.New("BoardPromotionOverlay is nil")
	}

	if p.Shown {
		for _, piece := range p.Pieces {
			piece.Piece.Color = p.Color

			if err := tools.Update(piece); err != nil {
				return err
			}
		}
		p.Get("classList").Call("add", "show")
	} else {
		p.Get("classList").Call("remove", "show")
	}

	return nil
}

type ModelBoard struct {
	shf.Element
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
		if err := tools.Initialize(b.Edgings.TopLeft); err != nil {
			return err
		}

		b.Edgings.TopLeft.SetPosition("top-left")
	}
	if b.Edgings.Top == nil {
		b.Edgings.Top = &EdgingHorizontal{}
		if err := tools.Initialize(b.Edgings.Top); err != nil {
			return err
		}

		b.Edgings.Top.SetPosition("top")

	}
	if b.Edgings.TopRight == nil {
		b.Edgings.TopRight = &EdgingCornerRotating{}
		if err := tools.Initialize(b.Edgings.TopRight); err != nil {
			return err
		}
		b.Edgings.TopRight.SetPosition("top-right")
	}

	if b.Edgings.Left == nil {
		b.Edgings.Left = &EdgingVertical{}
		if err := tools.Initialize(b.Edgings.Left); err != nil {
			return err
		}
		b.Edgings.Left.SetPosition("left")
	}

	if b.Grid == nil {
		b.Grid = &BoardGrid{}
		if err := tools.Initialize(b.Grid); err != nil {
			return err
		}
	}

	if b.Edgings.Right == nil {
		b.Edgings.Right = &EdgingVertical{}
		if err := tools.Initialize(b.Edgings.Right); err != nil {
			return err
		}
		b.Edgings.Right.SetPosition("right")
	}

	if b.Edgings.BottomLeft == nil {
		b.Edgings.BottomLeft = &EdgingCornerRotating{}
		if err := tools.Initialize(b.Edgings.BottomLeft); err != nil {
			return err
		}
		b.Edgings.BottomLeft.SetPosition("bottom-left")
	}
	if b.Edgings.Bottom == nil {
		b.Edgings.Bottom = &EdgingHorizontal{}
		if err := tools.Initialize(b.Edgings.Bottom); err != nil {
			return err
		}
		b.Edgings.Bottom.SetPosition("bottom")
	}
	if b.Edgings.BottomRight == nil {
		b.Edgings.BottomRight = &EdgingCorner{}
		if err := tools.Initialize(b.Edgings.BottomRight); err != nil {
			return err
		}
		b.Edgings.BottomRight.SetPosition("bottom-right")
	}

	if b.PromotionOverlay == nil {
		b.PromotionOverlay = &BoardPromotionOverlay{}
		if err := tools.Initialize(b.PromotionOverlay); err != nil {
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
	shf.Element
	Color            piece.Color
	PieceCount       map[piece.Type]int
	LastMoveThrowOut piece.Type
	pieces           map[piece.Type][2]shf.Element // 0: outer, 1:inner
}

func (c *ThrownOutsContainer) Init(tools *shf.Tools) error {
	if c.PieceCount == nil {
		c.PieceCount = map[piece.Type]int{}
	}
	if c.pieces == nil {
		c.pieces = map[piece.Type][2]shf.Element{}
		for _, pieceType := range thrownOutPiecesOrderType {
			div := tools.CreateElement("div")
			div.Get("classList").Call("add", "piececount")
			div.Call("appendChild", pieceElement(tools, piece.New(c.Color, pieceType)).Object())

			span := tools.CreateElement("span")
			span.Get("classList").Call("add", "count")
			div.Call("appendChild", span.Object())

			c.pieces[pieceType] = [2]shf.Element{div, span}
		}
	}
	if c.Element == nil {
		c.Element = tools.CreateElement("div")
		c.Set("id", "thrown-outs-"+strings.ToLower(c.Color.String()))
		c.Get("classList").Call("add", "thrown-outs")
		for _, pieceType := range thrownOutPiecesOrderType {
			c.Call("appendChild", c.pieces[pieceType][0].Object())
		}
	}
	return nil
}
func (c *ThrownOutsContainer) Update(tools *shf.Tools) error {
	if c == nil {
		return errors.New("ThrownOutsContainer is nil")
	}

	for pieceType, elms := range c.pieces {
		if c.LastMoveThrowOut == pieceType {
			elms[0].Get("classList").Call("add", "last-move")
		} else {
			elms[0].Get("classList").Call("remove", "last-move")
		}
		if c.PieceCount[pieceType] == 0 {
			elms[0].Get("classList").Call("add", "hidden")
		} else {
			elms[0].Get("classList").Call("remove", "hidden")
		}

		elms[1].Set("textContent", strconv.Itoa(c.PieceCount[pieceType]))

	}

	return nil
}

type ModelThrownouts struct {
	shf.Element
	White, Black *ThrownOutsContainer
}

func (t *ModelThrownouts) Init(tools *shf.Tools) error {
	if t.White == nil {
		t.White = &ThrownOutsContainer{
			Color: piece.White,
		}
		if err := tools.Initialize(t.White); err != nil {
			return err
		}
	}
	if t.Black == nil {
		t.Black = &ThrownOutsContainer{
			Color: piece.Black,
		}
		if err := tools.Initialize(t.Black); err != nil {
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

type StatusIcons struct {
	shf.Element
	White, Black bool
}

func (sm *StatusIcons) Init(tools *shf.Tools) error {
	if sm.Element == nil {
		sm.Element = tools.CreateElement("p")
		sm.Set("id", "game-status-icon")
	}
	return nil
}
func (sm *StatusIcons) Update(tools *shf.Tools) error {
	if sm == nil {
		return errors.New("StatusIcons is nil")
	}

	sm.Set("innerHTML", "")
	if sm.White {
		sm.Call("appendChild", pieceElement(tools, piece.New(piece.White, piece.King)).Object())
	}
	if sm.Black {
		sm.Call("appendChild", pieceElement(tools, piece.New(piece.Black, piece.King)).Object())
	}

	return nil
}

type StatusText struct {
	shf.Element
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

type StatusHeader struct {
	shf.Element
	Icons   *StatusIcons
	Message *StatusText
}

func (gs *StatusHeader) Init(tools *shf.Tools) error {
	if gs.Icons == nil {
		gs.Icons = &StatusIcons{}
		if err := tools.Initialize(gs.Icons); err != nil {
			return err
		}
	}
	if gs.Message == nil {
		gs.Message = &StatusText{}
		if err := tools.Initialize(gs.Message); err != nil {
			return err
		}
	}

	if gs.Element == nil {
		gs.Element = tools.CreateElement("div")
		gs.Set("id", "game-status-header")
		gs.Call("appendChild", gs.Icons.Element.Object())
		gs.Call("appendChild", gs.Message.Element.Object())
	}
	return nil
}
func (gs *StatusHeader) Update(tools *shf.Tools) error {
	if gs == nil {
		return errors.New("StatusHeader is nil")
	}

	return tools.Update(gs.Message, gs.Icons)
}

type ControlButton struct {
	shf.Element
	Text     string
	Hash     string
	Disabled bool
}

func (cb *ControlButton) Init(tools *shf.Tools) error {

	if cb.Element == nil {
		cb.Element = tools.CreateElement("button")

		if err := tools.Click(cb.Element, func(e shf.Event) error {
			if !cb.Disabled {
				js.Global().Get("location").Set("hash", cb.Hash)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
func (cb *ControlButton) Update(tools *shf.Tools) error {
	if cb == nil {
		return errors.New("ControlButton is nil")
	}
	cb.Set("textContent", cb.Text)
	if cb.Disabled {
		cb.Call("setAttribute", "disabled", "disabled")
	} else {
		cb.Call("removeAttribute", "disabled")
	}

	return nil
}

type StatusControl struct {
	shf.Element
	Start, Previous, Next, Initial *ControlButton

	refGame *ChessGameModel
}

func (sc *StatusControl) Init(tools *shf.Tools) error {
	if sc.Start == nil {
		sc.Start = &ControlButton{Text: ""}
		if err := tools.Initialize(sc.Start); err != nil {
			return err
		}
		sc.Start.Get("classList").Call("add", "start")
	}
	if sc.Previous == nil {
		sc.Previous = &ControlButton{Text: ""}
		if err := tools.Initialize(sc.Previous); err != nil {
			return err
		}
		sc.Previous.Get("classList").Call("add", "previous")
	}

	if sc.Next == nil {
		sc.Next = &ControlButton{Text: ""}
		if err := tools.Initialize(sc.Next); err != nil {
			return err
		}
		sc.Next.Get("classList").Call("add", "next")
	}
	if sc.Initial == nil {
		sc.Initial = &ControlButton{Text: ""}
		if err := tools.Initialize(sc.Initial); err != nil {
			return err
		}
		sc.Initial.Get("classList").Call("add", "end")
	}

	if sc.Element == nil {
		sc.Element = tools.CreateElement("div")
		sc.Set("id", "game-status-control")
		sc.Call("appendChild", sc.Start.Element.Object())
		sc.Call("appendChild", sc.Previous.Element.Object())
		sc.Call("appendChild", sc.Next.Element.Object())
		sc.Call("appendChild", sc.Initial.Element.Object())
	}
	return nil
}

func (sc *StatusControl) Update(tools *shf.Tools) error {
	if sc == nil {
		return errors.New("StatusControl is nil")
	}

	return tools.Update(sc.Start, sc.Previous, sc.Next, sc.Initial)
}

func (sc *StatusControl) rebuild(tools *shf.Tools) error {
	if sc.refGame == nil {
		return errors.New("StatusControl.Previous click callback: refGame is nil")
	}
	err := error(nil)

	lenInitialMoves, lenCurrentMoves := len(sc.refGame.initialPgn.Moves), len(sc.refGame.pgn.Moves)
	split := false
	for i := range sc.refGame.pgn.Moves {
		if i >= lenInitialMoves || sc.refGame.pgn.Moves[i] != sc.refGame.initialPgn.Moves[i] {
			split = true
			break
		}
	}

	if lenCurrentMoves == 0 {
		sc.Start.Disabled = true
		sc.Previous.Disabled = true
	} else {
		sc.Start.Disabled = false
		sc.Previous.Disabled = false

		sc.Previous.Hash, err = sc.refGame.HashForHalfMove(lenCurrentMoves - 1)
		if err != nil {
			return err
		}
	}

	sc.Next.Disabled = true
	if !split && lenInitialMoves > lenCurrentMoves {
		sc.Next.Hash, err = sc.refGame.HashForInitialHalfMove(lenCurrentMoves + 1)
		if err != nil {
			return err
		}
		sc.Next.Disabled = false
	}

	sc.Initial.Disabled = true
	if split || lenCurrentMoves != lenInitialMoves {
		sc.Initial.Hash, err = sc.refGame.HashForInitialHalfMove(lenInitialMoves)
		if err != nil {
			return err
		}
		sc.Initial.Disabled = false
		if split {
			sc.Initial.Get("classList").Call("add", "initial")
			sc.Initial.Get("classList").Call("remove", "end")
		} else {
			sc.Initial.Get("classList").Call("add", "end")
			sc.Initial.Get("classList").Call("remove", "initial")
		}
	} else {
		sc.Initial.Get("classList").Call("add", "end")
		sc.Initial.Get("classList").Call("remove", "initial")
	}

	return nil
}

type StatusMove struct {
	shf.Element
	Href             string
	Color            piece.Color
	Text             string
	Initial, Current bool
	Future, Splited  bool
}

func (sm *StatusMove) Init(tools *shf.Tools) error {
	if sm.Element == nil {
		sm.Element = tools.CreateElement("a")
	}
	return nil
}
func (sm *StatusMove) Update(tools *shf.Tools) error {
	//println("StatusMove.Update: text:", sm.Text)
	if sm == nil {
		return errors.New("StatusMove is nil")
	}
	classColor := ""
	if sm.Color == piece.White || sm.Color == piece.Black {
		classColor = " " + strings.ToLower(sm.Color.String())
	}
	classClickable := ""
	if sm.Href != "" && !sm.Current {
		sm.Element.Set("href", sm.Href)
		classClickable = " clickable"
	}
	classInitial := ""
	if sm.Initial {
		classInitial = " initial"
	}
	classCurrent := ""
	if sm.Current {
		classCurrent = " current"
	}
	classFuture := ""
	if sm.Future {
		classFuture = " future"
	}
	classSplited := ""
	if sm.Splited {
		classSplited = " splited"
	}

	sm.Get("classList").Set("value", "move"+classColor+classClickable+classInitial+classCurrent+classFuture+classSplited)
	sm.Set("textContent", sm.Text)
	return nil
}

type StatusMoves struct {
	shf.Element
	MoveZero            *StatusMove
	Moves               []*StatusMove
	SplitLastMove       *StatusMove
	ScrollToCurrentMove bool

	refGame  *ChessGameModel
	refModel *HtmlModel
}

func (sb *StatusMoves) Init(tools *shf.Tools) error {
	//println("StatusMoves.Init")
	if sb.MoveZero == nil {
		mz, err := sb.createHalfMoveNo(tools, StatusMove{nil, "#", piece.NoColor, "New game position", false, false, false, false})
		if err != nil {
			return err
		}
		sb.MoveZero = mz
	}

	if sb.Element == nil {
		sb.Element = tools.CreateElement("div")
		sb.Set("id", "game-status-moves")
	}
	return nil
}

func (sb *StatusMoves) createHalfMoveNo(tools *shf.Tools, sm StatusMove) (*StatusMove, error) {
	if err := tools.Initialize(&sm); err != nil {
		return nil, err
	}

	if sm.Href != "" && !sm.Current {
		if err := tools.Click(sm.Element, func(e shf.Event) error {
			//println("Clicked:", sm.Text)
			e.Call("stopPropagation")
			js.Global().Get("location").Set("hash", sm.Href)
			return nil
		}); err != nil {
			return nil, err
		}
	} else if sm.Current {
		if err := tools.Click(sm.Element, func(_ shf.Event) error {
			if err := sb.refModel.CopyGameURLToClipboard(); err != nil {
				return err
			}
			sb.refModel.Notification.TimedMessage(
				tools,
				5*time.Second,
				"Game URL was copied to clipboard",
				"",
			)
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return nil, err
		}
	}
	return &sm, nil
}

func (sb *StatusMoves) rebuild(tools *shf.Tools) error {
	if sb.refGame == nil {
		return errors.New("StatusMoves.rebuild: refGame is nil")
	}
	if sb.refModel == nil {
		return errors.New("StatusMoves.rebuild: refModel is nil")
	}
	if sb.MoveZero == nil {
		return errors.New("StatusMoves.rebuild: Movezero is nil")
	}

	for _, mv := range sb.Moves {
		tools.ClickRemove(mv)
		tools.Destroy(mv)
	}
	sb.Moves = nil
	if sb.SplitLastMove != nil {
		tools.ClickRemove(sb.SplitLastMove)
		tools.Destroy(sb.SplitLastMove)
	}
	sb.Set("innerHTML", "")
	sb.ScrollToCurrentMove = true

	lenInitialMoves, lenCurrentMoves := len(sb.refGame.initialPgn.Moves), len(sb.refGame.pgn.Moves)
	//println("lenInitialMoves:", lenInitialMoves, "lenCurrentMoves:", lenCurrentMoves)

	{ // Move zero
		sb.MoveZero.Initial = lenInitialMoves == 0
		sb.MoveZero.Current = lenCurrentMoves == 0

		moveNo := tools.CreateElement("span")
		moveNo.Get("classList").Call("add", "move-no")
		moveNo.Set("textContent", "0")

		p := tools.CreateElement("p")
		p.Get("classList").Call("add", "move-0")
		p.Call("appendChild", moveNo.Object())
		p.Call("appendChild", sb.MoveZero.Object())
		sb.Call("appendChild", p.Object())
	}

	wasSplit := false

	maxMovesLen := lenCurrentMoves
	if lenInitialMoves > maxMovesLen {
		maxMovesLen = lenInitialMoves
	}

	for i := 0; i < maxMovesLen; i += 1 {
		if wasSplit && i >= lenCurrentMoves {
			//println("there was a split and i:", i, " >= lenCurrentMoves:", lenCurrentMoves)
			break
		}
		if i%2 == 1 { // Skip every odd half move.
			continue
		}

		if !wasSplit {
			if i < lenCurrentMoves && i+1 < lenInitialMoves && sb.refGame.pgn.Moves[i] != sb.refGame.initialPgn.Moves[i] {
				//println("split at i:", i)

				moveNo := tools.CreateElement("span")
				moveNo.Get("classList").Call("add", "move-no")
				moveNo.Set("textContent", strconv.Itoa(lenInitialMoves/2))

				hash, err := sb.refGame.HashForInitialHalfMove(lenInitialMoves)
				if err != nil {
					return err
				}
				color := piece.Colors[(lenInitialMoves-1)%2]
				text := "... " + sb.refGame.initialPgn.Moves[lenInitialMoves-1]
				sb.SplitLastMove, err = sb.createHalfMoveNo(tools, StatusMove{nil, "#" + hash, color, text, true, false, false, false})
				if err != nil {
					return err
				}

				split := tools.CreateElement("p")
				split.Get("classList").Call("add", "split")
				split.Call("appendChild", moveNo.Object())
				split.Call("appendChild", sb.SplitLastMove.Object())
				sb.Call("appendChild", split.Object())

				wasSplit = true
			}
		}

		no := (i / 2) + 1 // Current move number (not half-move).
		hno := i + 1      // Half-move number.
		//println("i:", i, "no:", no, "hno:", hno)
		future := !wasSplit && i >= lenCurrentMoves
		//println("future:", future)

		moveNo := tools.CreateElement("span")
		moveNo.Get("classList").Call("add", "move-no")
		if future {
			moveNo.Get("classList").Call("add", "future")
		}
		moveNo.Set("textContent", strconv.Itoa(no))

		text, hash := "", ""
		var initial, current bool
		var err error
		if future {
			text = sb.refGame.initialPgn.Moves[i]
			hash, err = sb.refGame.HashForInitialHalfMove(hno)
			if err != nil {
				return err
			}
			initial = i+1 == lenInitialMoves
			current = false
		} else {
			text = sb.refGame.pgn.Moves[i]
			hash, err = sb.refGame.HashForHalfMove(hno)
			if err != nil {
				return err
			}
			initial = !wasSplit && i+1 == lenInitialMoves && sb.refGame.pgn.Moves[i] == sb.refGame.initialPgn.Moves[i]
			current = i+1 == lenCurrentMoves
		}
		moveWhite, err := sb.createHalfMoveNo(tools, StatusMove{nil, "#" + hash, piece.White, text, initial, current, future, wasSplit || i >= lenInitialMoves})
		if err != nil {
			return err
		}
		sb.Moves = append(sb.Moves, moveWhite)

		p := tools.CreateElement("p")
		p.Get("classList").Call("add", "move-"+strconv.Itoa(no))
		p.Call("appendChild", moveNo.Object())
		p.Call("appendChild", moveWhite.Object())

		j := i + 1
		//println("j:", j)
		if wasSplit && j >= lenCurrentMoves {
			//println("there was a split and j:", j, " >= lenCurrentMoves:", lenCurrentMoves)
		} else if j < maxMovesLen {

			if !wasSplit {
				if j < lenCurrentMoves && j < lenInitialMoves && sb.refGame.pgn.Moves[j] != sb.refGame.initialPgn.Moves[j] {

					// Add empty black move and append row.
					moveBlack, err := sb.createHalfMoveNo(tools, StatusMove{nil, "", piece.Black, "", false, false, false, false})
					if err != nil {
						return err
					}
					moveBlack.Update(tools)
					p.Call("appendChild", moveBlack.Object())
					sb.Call("appendChild", p.Object())

					// Add split.

					moveNo := tools.CreateElement("span")
					moveNo.Get("classList").Call("add", "move-no")
					moveNo.Set("textContent", strconv.Itoa(lenInitialMoves/2))

					hash, err := sb.refGame.HashForInitialHalfMove(lenInitialMoves)
					if err != nil {
						return err
					}
					color := piece.Colors[(lenInitialMoves-1)%2]
					text := sb.refGame.initialPgn.Moves[lenInitialMoves-1]
					sb.SplitLastMove, err = sb.createHalfMoveNo(tools, StatusMove{nil, "#" + hash, color, text, true, false, false, false})
					if err != nil {
						return err
					}

					split := tools.CreateElement("p")
					split.Get("classList").Call("add", "split")
					split.Call("appendChild", moveNo.Object())
					split.Call("appendChild", sb.SplitLastMove.Object())
					sb.Call("appendChild", split.Object())

					// Create new row with the same move number and blank white move.
					moveNo = tools.CreateElement("span")
					moveNo.Get("classList").Call("add", "move-no")
					if future {
						moveNo.Get("classList").Call("add", "future")
					}
					moveNo.Set("textContent", strconv.Itoa(no))

					moveWhite, err := sb.createHalfMoveNo(tools, StatusMove{nil, "", piece.White, "", false, false, false, true})
					if err != nil {
						return err
					}
					moveWhite.Update(tools)
					p = tools.CreateElement("p")
					p.Get("classList").Call("add", "move-"+strconv.Itoa(no))
					p.Call("appendChild", moveNo.Object())
					p.Call("appendChild", moveWhite.Object())

					wasSplit = true
				}
			}

			future := !wasSplit && j >= lenCurrentMoves
			//println("future:", future)
			text, hash := "", ""
			var initial, current bool
			var err error
			if future {
				text = sb.refGame.initialPgn.Moves[j]
				hash, err = sb.refGame.HashForInitialHalfMove(hno + 1)
				if err != nil {
					return err
				}
				initial = j+1 == lenInitialMoves
				current = false
			} else {
				text = sb.refGame.pgn.Moves[j]
				hash, err = sb.refGame.HashForHalfMove(hno + 1)
				if err != nil {
					return err
				}
				initial = !wasSplit && j+1 == lenInitialMoves && sb.refGame.pgn.Moves[j] == sb.refGame.initialPgn.Moves[j]
				current = j+1 == lenCurrentMoves
			}
			moveBlack, err := sb.createHalfMoveNo(tools, StatusMove{nil, "#" + hash, piece.Black, text, initial, current, future, wasSplit || j >= lenInitialMoves})
			if err != nil {
				return err
			}
			sb.Moves = append(sb.Moves, moveBlack)
			p.Call("appendChild", moveBlack.Object())
		}

		sb.Call("appendChild", p.Object())

	}

	return nil
}
func (sb *StatusMoves) Update(tools *shf.Tools) error {
	//println("StatusMoves.Update")
	if sb == nil {
		return errors.New("StatusMoves is nil")
	}

	if err := sb.MoveZero.Update(tools); err != nil {
		return err
	}
	for _, sm := range sb.Moves {
		if err := sm.Update(tools); err != nil {
			return err
		}
	}
	if sb.SplitLastMove != nil {
		if err := sb.SplitLastMove.Update(tools); err != nil {
			return err
		}
	}

	if sb.ScrollToCurrentMove {
		parentElement := sb.Element
		targetElement := sb.MoveZero.Element
		lenCurrentMoves := len(sb.refGame.pgn.Moves)
		if lenCurrentMoves > 0 {
			targetElement = sb.Moves[lenCurrentMoves-1].Element
		}

		targetOffsetHeight := targetElement.Get("offsetHeight").Int()
		targetTop := targetElement.Get("offsetTop").Int() - sb.MoveZero.Element.Get("offsetTop").Int()
		parentClientHeight := parentElement.Get("clientHeight").Int()
		parentVisibleTop := parentElement.Get("scrollTop").Int()

		if targetTop < parentVisibleTop || targetTop+targetOffsetHeight > parentVisibleTop+parentClientHeight {
			scrollTo := targetTop - parentClientHeight/2 + targetOffsetHeight/2

			parentElement.Call("scrollTo", map[string]interface{}{
				"top":      scrollTo,
				"behavior": "smooth",
			})
		}

		sb.ScrollToCurrentMove = false
	}
	return nil
}

type ModelGameStatus struct {
	shf.Element
	Header  *StatusHeader
	Control *StatusControl
	Moves   *StatusMoves
}

func (gs *ModelGameStatus) Init(tools *shf.Tools) error {
	if gs.Header == nil {
		gs.Header = &StatusHeader{}
		if err := tools.Initialize(gs.Header); err != nil {
			return err
		}
	}
	if gs.Control == nil {
		gs.Control = &StatusControl{}
		if err := tools.Initialize(gs.Control); err != nil {
			return err
		}
	}
	if gs.Moves == nil {
		gs.Moves = &StatusMoves{}
		if err := tools.Initialize(gs.Moves); err != nil {
			return err
		}
	}

	if gs.Element == nil {
		gs.Element = tools.CreateElement("div")
		gs.Set("id", "game-status")
		gs.Call("appendChild", gs.Header.Element.Object())
		gs.Call("appendChild", gs.Control.Element.Object())
		gs.Call("appendChild", gs.Moves.Element.Object())
	}
	return nil
}
func (gs *ModelGameStatus) Update(tools *shf.Tools) error {
	if gs == nil {
		return errors.New("ModelGameStatus is nil")
	}

	return tools.Update(gs.Header, gs.Control, gs.Moves)
}

func (gs *ModelGameStatus) rebuild(tools *shf.Tools) error {
	if err := gs.Control.rebuild(tools); err != nil {
		return err
	}
	if err := gs.Moves.rebuild(tools); err != nil {
		return err
	}

	return nil
}

type CopyButton struct {
	shf.Element
	Shown bool
}

func (this *CopyButton) Init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = tools.CreateElement("button")
		this.Set("textContent", "Copy to clipboard")

	}
	return nil
}
func (this *CopyButton) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("CopyButton is nil")
	}
	if this.Shown {
		this.Get("classList").Call("remove", "hidden")
	} else {
		this.Get("classList").Call("add", "hidden")
	}
	return nil
}

type MoveStatusLink struct {
	shf.Element
	MoveHash string

	Input shf.Element
	Copy  *CopyButton
}

func (this *MoveStatusLink) GetURL() string {
	hash := "#" + strings.TrimPrefix(this.MoveHash, "#")
	loc := js.Global().Get("location")
	return loc.Get("origin").String() + loc.Get("pathname").String() + hash
}
func (this *MoveStatusLink) Init(tools *shf.Tools) error {
	if this.Input == nil {
		this.Input = tools.CreateElement("input")
		this.Input.Set("type", "text")
		this.Input.Call("setAttribute", "readonly", "readonly")
	}

	if this.Copy == nil {
		this.Copy = &CopyButton{}
		if err := tools.Initialize(this.Copy); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Get("classList").Call("add", "link")

		this.Call("appendChild", this.Input.Object())
		this.Call("appendChild", this.Copy.Object())
		this.Call("appendChild", tools.CreateTextNode("This URL link represents the state of current chess game. You can copy it and store it or send it."))

	}
	return nil
}
func (this *MoveStatusLink) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("MoveStatusLink is nil")
	}

	this.Input.Set("value", this.GetURL())

	return tools.Update(this.Copy)
}

type ModelCover struct {
	shf.Element
	GameStatus *ModelGameStatus
	MoveStatus *ModelMoveStatus
}

type ModelMoveStatus struct {
	shf.Element
	Shown bool

	Link *MoveStatusLink
	//Navigation *MoveStatusNavigation
	Undo  shf.Element
	Close shf.Element
	tip   shf.Element
	tipNo int
}

func (this *ModelMoveStatus) Init(tools *shf.Tools) error {
	if this.Link == nil {
		this.Link = &MoveStatusLink{}
		if err := tools.Initialize(this.Link); err != nil {
			return err
		}
	}

	if this.Undo == nil {
		this.Undo = tools.CreateElement("button")
		this.Undo.Set("textContent", "back")
		// Model sets Event
		// Model.ChessGame.UpdateModel sets visibility
	}

	if this.Close == nil {
		this.Close = tools.CreateElement("button")
		this.Close.Set("textContent", "close")
		if err := tools.Click(this.Close, func(_ shf.Event) error {
			this.Shown = false
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}

	if this.tip == nil {
		this.tip = tools.CreateElement("div")
		this.tip.Set("className", "tip")
		if err := tools.Click(this.tip, func(_ shf.Event) error {
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Set("id", "move-status")

		this.Call("appendChild", this.Link.Object())

		{
			div := tools.CreateElement("div")
			div.Call("appendChild", this.Undo.Object())

			this.Call("appendChild", div.Object())
		}
		{
			div := tools.CreateElement("div")
			div.Call("appendChild", this.Close.Object())

			this.Call("appendChild", div.Object())
		}
		this.Call("appendChild", this.tip.Object())
	}
	return nil
}

var tips = []string{
	"click on this tip, to see next tip",
	"to close notifications, click anywhere except buttons in notification",
	"to quick copy game to clipboard, click on last move piece",
	"to go back a move, click on last move FROM square",
	"to rotate board (if supported), click on corners with arrows (bottom-left, top-right)",
	"to rotate board for current moving player, click on game status icon, or text",
	"to toggle this game URL dialog, click on any empty square on board",
	"to toggle zen mode, try double click on empty chess square",
}

func (this *ModelMoveStatus) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelMoveStatus is nil")
	}

	if this.Shown {
		if len(tips) > 0 {
			this.tip.Set("textContent", "tip: "+tips[this.tipNo])
			this.tipNo = (this.tipNo + 1) % len(tips)
		}
		this.Get("classList").Call("remove", "hidden")
	} else {
		this.Get("classList").Call("add", "hidden")
	}

	return tools.Update(this.Link)
}

func (this *ModelCover) Init(tools *shf.Tools) error {
	if this.GameStatus == nil {
		this.GameStatus = &ModelGameStatus{}
		if err := tools.Initialize(this.GameStatus); err != nil {
			return err
		}
	}
	if this.MoveStatus == nil {
		this.MoveStatus = &ModelMoveStatus{}
		if err := tools.Initialize(this.MoveStatus); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Set("id", "cover")

		this.Call("appendChild", this.GameStatus.Element.Object())
		this.Call("appendChild", this.MoveStatus.Element.Object())
	}
	return nil
}
func (this *ModelCover) Update(tools *shf.Tools) error {
	return tools.Update(this.GameStatus, this.MoveStatus)
}

type ModelExportTagInput struct {
	shf.Element
	Name, Label string
	Input       shf.Element
}

func (this *ModelExportTagInput) Init(tools *shf.Tools) error {
	if this.Input == nil {
		this.Input = tools.CreateElement("input")
	}

	if this.Element == nil {
		label := tools.CreateElement("span")
		label.Set("textContent", this.Label)

		this.Element = tools.CreateElement("p")
		this.Set("id", "export-tag-"+strings.ToLower(this.Name))
		this.Element.Get("classList").Call("add", "tag")

		this.Call("appendChild", label.Object())
		this.Call("appendChild", this.Input.Object())
	}

	return nil
}
func (this *ModelExportTagInput) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelExportTagInput is nil")
	}

	return nil
}

type ModelExportTagSelect struct {
	shf.Element
	Select      shf.Element
	Options     [][2]string
	Selected    string
	Name, Label string
	Disabled    bool
}

func (this *ModelExportTagSelect) Init(tools *shf.Tools) error {
	if this.Select == nil {
		this.Select = tools.CreateElement("select")
	}

	if this.Element == nil {
		label := tools.CreateElement("span")
		label.Set("textContent", this.Label)

		this.Element = tools.CreateElement("p")
		this.Set("id", "export-tag-"+strings.ToLower(this.Name))
		this.Element.Get("classList").Call("add", "tag")

		this.Call("appendChild", label.Object())
		this.Call("appendChild", this.Select.Object())
	}

	return nil
}
func (this *ModelExportTagSelect) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelExportTagSelect is nil")
	}

	if this.Select != nil {
		this.Select.Set("innerHTML", "")

		for _, o := range this.Options {

			option := tools.CreateElement("option")
			option.Set("textContent", o[1])
			option.Set("value", o[0])
			if this.Selected == o[0] {
				option.Call("setAttribute", "selected", "selected")
			}

			this.Select.Call("appendChild", option.Object())
		}

		if this.Disabled {
			this.Select.Call("setAttribute", "disabled", "disabled")
		} else {
			this.Select.Call("removeAttribute", "disabled")
		}
	}

	return nil
}

type ModelExportInput struct {
	shf.Element
	White, Black *ModelExportTagInput
	Round, Date  *ModelExportTagInput
	Result       *ModelExportTagSelect
}

func (this *ModelExportInput) Init(tools *shf.Tools) error {
	// PGN Tags known by lichess (from: https://github.com/lichess-org/lila/blob/master/modules/study/src/main/PgnTags.scala#L30)
	// White,
	// WhiteElo,
	// WhiteTitle,
	// WhiteTeam,
	// WhiteFideId,
	// Black,
	// BlackElo,
	// BlackTitle,
	// BlackTeam,
	// BlackFideId,
	// TimeControl,
	// Date,
	// Result,
	// Termination,
	// Site,
	// Event,
	// Round,
	// Board,
	// Annotator

	if this.White == nil {
		this.White = &ModelExportTagInput{Name: "White", Label: "White name"}
		if err := tools.Initialize(this.White); err != nil {
			return err
		}
	}
	if this.Black == nil {
		this.Black = &ModelExportTagInput{Name: "Black", Label: "Black name"}
		if err := tools.Initialize(this.Black); err != nil {
			return err
		}
	}
	if this.Round == nil {
		this.Round = &ModelExportTagInput{Name: "Round", Label: "Round no."}
		if err := tools.Initialize(this.Round); err != nil {
			return err
		}
	}
	if this.Date == nil {
		this.Date = &ModelExportTagInput{Name: "Date", Label: "Start date"}
		if err := tools.Initialize(this.Date); err != nil {
			return err
		}
	}
	if this.Result == nil {
		this.Result = &ModelExportTagSelect{Name: "Result", Label: "Result",
			Options: [][2]string{
				{"*", game.InProgress.String()},
				{"1-0", game.WhiteWon.String()},
				{"1/2-1/2", game.Draw.String()},
				{"0-1", game.BlackWon.String()},
			}}
		if err := tools.Initialize(this.Result); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Element.Get("classList").Call("add", "tags")

		this.Element.Call("appendChild", this.White.Object())
		this.Element.Call("appendChild", this.Black.Object())
		this.Element.Call("appendChild", this.Round.Object())
		this.Element.Call("appendChild", this.Date.Object())
		this.Element.Call("appendChild", this.Result.Object())
	}
	return nil
}
func (this *ModelExportInput) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelExportInput is nil")
	}

	if err := tools.Update(this.White, this.Black, this.Round, this.Date, this.Result); err != nil {
		return err
	}
	return nil
}

type CloseButton struct {
	shf.Element
}

func (this *CloseButton) Init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = tools.CreateElement("button")
		this.Set("textContent", "close")
	}
	return nil
}
func (this *CloseButton) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("CloseButton is nil")
	}
	return nil
}

type ModelExportOutput struct {
	shf.Element
	PGN *pgn.PGN

	TextArea shf.Element
	Copy     *CopyButton
	Close    *CloseButton
}

func (this *ModelExportOutput) Init(tools *shf.Tools) error {
	if this.PGN == nil {
		this.PGN = &pgn.PGN{}
	}

	if this.Copy == nil {
		this.Copy = &CopyButton{Shown: true}
		if err := tools.Initialize(this.Copy); err != nil {
			return err
		}
	}

	if this.Close == nil {
		this.Close = &CloseButton{}
		if err := tools.Initialize(this.Close); err != nil {
			return err
		}
	}

	if this.TextArea == nil {
		this.TextArea = tools.CreateElement("textarea")
		this.TextArea.Set("id", "export-output")
		this.TextArea.Call("setAttribute", "readonly", "readonly")
	}

	if this.Element == nil {

		buttons := tools.CreateElement("p")
		buttons.Get("classList").Call("add", "buttons")
		buttons.Call("appendChild", this.Copy.Object())
		buttons.Call("appendChild", this.Close.Object())

		this.Element = tools.CreateElement("div")
		this.Get("classList").Call("add", "output")
		this.Call("appendChild", this.TextArea.Object())
		this.Call("appendChild", buttons.Object())
	}

	return nil
}
func (this *ModelExportOutput) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelExportOutput is nil")
	}

	if this.PGN != nil {
		this.TextArea.Set("value", this.PGN.String())
	}

	return tools.Update(this.Copy, this.Close)
}

type ModelExport struct {
	shf.Element
	Shown bool

	Input  *ModelExportInput
	Output *ModelExportOutput
}

// EscapePGNString escapes special characters in a string for PGN tag values.
func EscapePGNString(s string) string {
	// Escape backslashes first to avoid escaping already escaped characters
	s = strings.ReplaceAll(s, "\\", "\\\\")
	// Escape double quotes
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func (this *ModelExport) applyTag(key, value string) error {
	if this.Output == nil || this.Output.PGN == nil || this.Output.PGN.Tags == nil {
		return errors.New("output PGN tags are not initialized")
	}
	if value != "" {
		//TODO - Move escaping to chess engine.
		this.Output.PGN.Tags[key] = EscapePGNString(value)
	} else {
		delete(this.Output.PGN.Tags, key)
	}
	return nil
}
func (this *ModelExport) Init(tools *shf.Tools) error {
	if this.Output == nil {
		this.Output = &ModelExportOutput{}
		if err := tools.Initialize(this.Output); err != nil {
			return err
		}

		if err := tools.Click(this.Output.Close.Element, func(_ shf.Event) error {
			this.Shown = false
			this.Output.PGN = nil
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}

	}
	if this.Input == nil {
		this.Input = &ModelExportInput{}
		if err := tools.Initialize(this.Input); err != nil {
			return err
		}

		// Bind tags input value change to update output PGN tags.
		if err := tools.Input(this.Input.White.Input, func(_ shf.Event) error {
			if err := this.applyTag(this.Input.White.Name, this.Input.White.Input.Get("value").String()); err != nil {
				return err
			}

			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
		if err := tools.Input(this.Input.Black.Input, func(_ shf.Event) error {
			if err := this.applyTag(this.Input.Black.Name, this.Input.Black.Input.Get("value").String()); err != nil {
				return err

			}

			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
		if err := tools.Input(this.Input.Round.Input, func(_ shf.Event) error {
			if err := this.applyTag(this.Input.Round.Name, this.Input.Round.Input.Get("value").String()); err != nil {
				return err

			}

			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
		if err := tools.Input(this.Input.Date.Input, func(_ shf.Event) error {
			if err := this.applyTag(this.Input.Date.Name, this.Input.Date.Input.Get("value").String()); err != nil {
				return err

			}

			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
		if err := tools.Input(this.Input.Result.Select, func(_ shf.Event) error {
			this.Input.Result.Selected = this.Input.Result.Select.Get("value").String()
			if err := this.applyTag(this.Input.Result.Name, this.Input.Result.Selected); err != nil {
				return err

			}

			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Set("id", "export-overlay")
		if err := tools.Click(this.Element, func(e shf.Event) error {
			if e.Get("target").Get("id").String() == "export-overlay" {
				this.Shown = false
				this.Output.PGN = nil
			}
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}

		export := tools.CreateElement("div")
		export.Get("classList").Call("add", "export")

		export.Call("appendChild", this.Input.Element.Object())
		export.Call("appendChild", this.Output.Element.Object())

		this.Call("appendChild", export.Object())
	}
	return nil
}
func (this *ModelExport) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("ModelExport is nil")
	}

	if err := tools.Update(this.Input); err != nil {
		return err
	}

	if this.Input != nil && this.Output != nil && this.Output.PGN != nil && this.Output.PGN.Tags != nil {
		// Populate input tags to output PGN.
		if this.Input.White != nil && this.Input.White.Input != nil {
			this.applyTag(this.Input.White.Name, this.Input.White.Input.Get("value").String())
		}
		if this.Input.Black != nil && this.Input.Black.Input != nil {
			this.applyTag(this.Input.Black.Name, this.Input.Black.Input.Get("value").String())
		}
		if this.Input.Round != nil && this.Input.Round.Input != nil {
			this.applyTag(this.Input.Round.Name, this.Input.Round.Input.Get("value").String())
		}
		if this.Input.Date != nil && this.Input.Date.Input != nil {
			this.applyTag(this.Input.Date.Name, this.Input.Date.Input.Get("value").String())
		}
		if this.Input.Result != nil {
			this.applyTag(this.Input.Result.Name, this.Input.Result.Selected)
		}
	}

	if err := tools.Update(this.Output); err != nil {
		return err
	}

	if this.Shown {
		this.Get("classList").Call("remove", "invisible")
	} else {
		this.Get("classList").Call("add", "invisible")
	}
	return nil
}

type ModelNotification struct {
	shf.Element
	Shown     bool
	timeoutId int
}

func (n *ModelNotification) cancelTimer() {
	if n.timeoutId == 0 {
		return
	}
	js.Global().Call("clearTimeout", n.timeoutId)
	//js.Global().Call("alert", "timer "+strconv.Itoa(n.timeoutId)+" canceled")
	n.timeoutId = 0
}
func (n *ModelNotification) Message(text, hint string, elements ...shf.Element) {
	n.cancelTimer()

	n.Set("innerHTML", "")

	//TODO jezek - Create properly using tools, save and destroy when changing.
	notification := shf.CreateElementObject("div")
	notification.Get("classList").Call("add", "notification")
	{ // message
		msg := shf.CreateElementObject("p")
		msg.Get("classList").Call("add", "message")
		msg.Set("textContent", text)
		notification.Call("appendChild", msg)
	}

	for _, e := range elements {
		notification.Call("appendChild", e.Object())
	}
	if hint != "" { // hint
		hintElm := shf.CreateElementObject("p")
		hintElm.Get("classList").Call("add", "hint")
		hintElm.Set("textContent", hint)
		notification.Call("appendChild", hintElm)
	}
	n.Call("appendChild", notification)
	n.Shown = true
}
func (n *ModelNotification) TimedMessage(tools *shf.Tools, duration time.Duration, text, hint string, elements ...shf.Element) {
	n.cancelTimer()

	//TODO jezek - Create properly using tools, save and destroy when changing.
	notification := shf.CreateElementObject("div")
	notification.Get("classList").Call("add", "notification")
	{ // message
		msg := shf.CreateElementObject("p")
		msg.Get("classList").Call("add", "message")
		msg.Set("textContent", text)
		notification.Call("appendChild", msg)
	}

	for _, e := range elements {
		notification.Call("appendChild", e.Object())
	}
	if hint != "" { // hint
		hintElm := shf.CreateElementObject("p")
		hintElm.Get("classList").Call("add", "hint")
		hintElm.Set("textContent", hint)
		notification.Call("appendChild", hintElm)
	}

	n.Set("innerHTML", "")
	n.Call("appendChild", notification)

	n.Shown = true
	n.timeoutId = tools.Timer(duration, func() {
		n.timeoutId = 0
		n.Shown = false
	})
}

func (n *ModelNotification) Init(tools *shf.Tools) error {

	if n.Element == nil {
		n.Element = tools.CreateElement("div")
		n.Set("id", "notification-overlay")
		if err := tools.Click(n.Element, func(e shf.Event) error {
			if e.Get("target").Get("id").String() == "notification-overlay" {
				n.cancelTimer()
				n.Shown = false
			}
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}
	return nil
}
func (n *ModelNotification) Update(tools *shf.Tools) error {
	if n == nil {
		return errors.New("ModelNotification is nil")
	}

	if n.Shown {
		n.Get("classList").Call("remove", "invisible")
	} else {
		n.Get("classList").Call("add", "invisible")
	}
	return nil
}

type HtmlModel struct {
	Rotated180deg bool

	Header       *ModelHeader
	Board        *ModelBoard
	ThrownOuts   *ModelThrownouts
	Cover        *ModelCover
	Export       *ModelExport
	Notification *ModelNotification
	Footer       *ModelFooter
}

func (h *HtmlModel) Init(tools *shf.Tools) error {
	if h.Header == nil {
		h.Header = &ModelHeader{}
		if err := tools.Initialize(h.Header); err != nil {
			return err
		}
	}
	if h.Board == nil {
		h.Board = &ModelBoard{}
		if err := tools.Initialize(h.Board); err != nil {
			return err
		}
	}
	if h.ThrownOuts == nil {
		h.ThrownOuts = &ModelThrownouts{}
		if err := tools.Initialize(h.ThrownOuts); err != nil {
			return err
		}
	}
	if h.Cover == nil {
		h.Cover = &ModelCover{}
		if err := tools.Initialize(h.Cover); err != nil {
			return err
		}
	}
	if h.Export == nil {
		h.Export = &ModelExport{}
		if err := tools.Initialize(h.Export); err != nil {
			return err
		}
	}
	if h.Notification == nil {
		h.Notification = &ModelNotification{}
		if err := tools.Initialize(h.Notification); err != nil {
			return err
		}
	}
	if h.Footer == nil {
		h.Footer = &ModelFooter{}
		if err := tools.Initialize(h.Footer); err != nil {
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

	return tools.Update(h.Header, h.Board, h.ThrownOuts, h.Cover, h.Export, h.Notification, h.Footer)
}

func (m *Model) RotateBoard() {
	if !m.rotationSupported {
		return
	}
	m.Html.Rotated180deg = !m.Html.Rotated180deg
}
func (h *HtmlModel) CopyGameURLToClipboard() error {
	positionX := js.Global().Get("pageXOffset")
	positionY := js.Global().Get("pageYOffset")

	temporaryShow := false
	if !h.Cover.MoveStatus.Shown {
		temporaryShow = true
	}
	if temporaryShow {
		h.Cover.MoveStatus.Element.Get("classList").Call("remove", "hidden")
	}
	h.Cover.MoveStatus.Link.Input.Call("focus")
	h.Cover.MoveStatus.Link.Input.Call("setSelectionRange", 0, len(h.Cover.MoveStatus.Link.Input.Get("value").String()))
	js.Global().Get("document").Call("execCommand", "Copy")
	h.Cover.MoveStatus.Link.Input.Call("blur")
	if temporaryShow {
		h.Cover.MoveStatus.Element.Get("classList").Call("add", "hidden")
	}

	js.Global().Call("scrollTo", positionX, positionY)
	return nil
}
func (h *HtmlModel) CopyExportOutputToClipboard() error {
	positionX := js.Global().Get("pageXOffset")
	positionY := js.Global().Get("pageYOffset")

	temporaryShow := false
	if !h.Export.Shown {
		temporaryShow = true
	}
	if temporaryShow {
		h.Export.Element.Get("classList").Call("remove", "invisible")
	}
	h.Export.Output.TextArea.Call("focus")
	h.Export.Output.TextArea.Call("setSelectionRange", 0, len(h.Export.Output.TextArea.Get("value").String()))
	js.Global().Get("document").Call("execCommand", "Copy")
	h.Export.Output.TextArea.Call("blur")
	if temporaryShow {
		h.Export.Element.Get("classList").Call("add", "invisible")
	}

	js.Global().Call("scrollTo", positionX, positionY)
	return nil
}

type ThrownOuts map[piece.Piece]uint8
type GameThrownOuts []ThrownOuts
type ChessGameModel struct {
	gameHash   string
	game       *game.Game
	gameGc     GameThrownOuts
	currMoveNo int
	nextMove   move.Move
	pgn        *pgn.PGN

	initialGame *game.Game
	initialPgn  *pgn.PGN
}

func addedThrownOuts(prev, next ThrownOuts) ThrownOuts {
	added := ThrownOuts{}
	for p, c := range prev {
		if next[p] > c {
			added[p] = added[p] + (next[p] - c)
		}
	}
	for p, c := range next {
		if _, ok := prev[p]; !ok {
			added[p] = added[p] + c
		}
	}
	return added
}

func pMakeMove(p *position.Position, m move.Move) (*position.Position, piece.Piece) {
	newPos := p.MakeMove(m)

	// was a piece thrown out regulary? = move destination contains some piece
	if pce := p.OnSquare(m.To()); pce.Type != piece.None {
		return newPos, pce
	}

	// was en passant throw out? = moved piece is pawn and move destination is an en passan square in previous move
	if mp := p.OnSquare(m.From()); mp.Type == piece.Pawn && m.To() == p.EnPassant {
		return newPos, piece.New(piece.Colors[(p.ActiveColor+1)%2], piece.Pawn)
	}

	return newPos, piece.New(piece.NoColor, piece.None)
}

// Creates new chess game from moves string.
// The moves string is basicaly move coordinates from & to (0...63) encoded in base64 (with some improvements for promotions, etc...). See encoding.go
func NewGame(hash string) (*ChessGameModel, error) {
	//println("NewGame(hash: \"" + hash + "\")")
	chgm := &ChessGameModel{}

	if err := chgm.UpdateToHash(hash); err != nil {
		return nil, err
	}

	chgm.initialGame = chgm.game
	chgm.initialPgn = chgm.pgn

	return chgm, nil
}

// Updates chess game to match moves from the moves hash string.
// The moves hash string are basically pair of move coordinates (from, to <0, 63>) encoded in base64 (with some improvements for promotions, etc...). See encoding.go
func (ch *ChessGameModel) UpdateToHash(hash string) error {
	//println("UpdateToHash(" + hash + ")")
	// Trim movesString from leading "#" character.
	movesString := strings.TrimPrefix(hash, "#")

	// Decode moves from hash string.
	moves, err := DecodeMoves(movesString)
	if err != nil {
		return errors.New("decoding moves error: " + err.Error())
	}

	// Create new game and thrown outs structures.
	g := game.New()
	gtos := make(GameThrownOuts, len(moves))

	// Apply decode game moves to new game.
	for i, move := range moves {
		if g.Status() != game.InProgress {
			return errors.New("Too many moves in url string! " + strconv.Itoa(i+1) + " moves are enough")
		}

		// Store position before the move.
		pbm := g.Position()

		// Make the move and check validity.
		_, merr := g.MakeMove(move)
		if merr != nil {
			return errors.New("Erroneous move number " + strconv.Itoa(i+1) + ": " + merr.Error())
		}

		// Create throw outs list for this move and copy thrown outs from previous move.
		tos := ThrownOuts{}
		if i > 0 {
			for p, c := range gtos[i-1] {
				tos[p] = c
			}
		}

		// If there is a thrown out in this move, add it to this move thrown out list.
		if _, top := pMakeMove(pbm, move); top.Type != piece.None {
			tos[top] = tos[top] + 1
		}

		gtos[i] = tos
	}

	// Prepend one empty throw outs structure to the thrown outs list.
	gtos = append(GameThrownOuts{ThrownOuts{}}, gtos...)

	// Update the ChessGameModel structure.
	ch.gameHash = movesString
	ch.game = g
	ch.gameGc = gtos
	ch.currMoveNo = len(gtos) - 1
	ch.nextMove = move.Null
	ch.pgn = pgn.EncodeSAN(g)

	return nil
}

func (ch *ChessGameModel) Validate() error {
	if ch == nil {
		return errors.New("ChessGame is nil")
	}
	if ch.currMoveNo < 0 || ch.currMoveNo >= len(ch.game.Positions) {
		return errors.New("current move number is out of bounds")
	}
	if len(ch.game.Positions) != len(ch.gameGc) {
		return errors.New("count of game moves and thrown outs does not match")
	}
	if _, err := getNextMoveState(ch.game.Positions[ch.currMoveNo], ch.nextMove); err != nil {
		return err
	}
	return nil
}

func (ch *ChessGameModel) MakeNextMove() error {
	if err := ch.Validate(); err != nil {
		return err
	}

	if !isLegalMove(ch.game.Positions[ch.currMoveNo], ch.nextMove) {
		return errors.New("can not make next move, next move is not a legal move ")
	}

	nextMoveHash, err := encodeMove(ch.nextMove)
	if err != nil {
		return err
	}

	{ // update game
		if _, err := ch.game.MakeMove(ch.nextMove); err != nil {
			return err
		}
		// update game hash
		ch.gameHash = ch.gameHash + nextMoveHash

	}
	{ // update throw outs
		// copy previous move throw outs
		nt := ThrownOuts{}
		for p, c := range ch.gameGc[ch.currMoveNo] {
			nt[p] = c
		}
		// add next move thrown out piece, if any
		if _, top := pMakeMove(ch.game.Positions[ch.currMoveNo], ch.nextMove); top.Type != piece.None {
			nt[top] = nt[top] + 1
		}
		ch.gameGc = append(ch.gameGc, nt)
	}

	// advance move number
	ch.currMoveNo = ch.currMoveNo + 1

	// reset next move
	ch.nextMove = move.Null

	ch.pgn = pgn.EncodeSAN(ch.game)

	// update location hash
	js.Global().Get("location").Set("hash", ch.gameHash)

	return nil
}
func (ch *ChessGameModel) BackToPreviousMove() error {
	if err := ch.Validate(); err != nil {
		return err
	}
	if ch.game.Positions[ch.currMoveNo].LastMove == move.Null {
		// no previous move, just return
		return nil
	}

	lastMove, err := encodeMove(ch.game.Positions[ch.currMoveNo].LastMove)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(ch.gameHash, lastMove) {
		return errors.New("last move is not suffix of game moves")
	}

	previousGameMoves := strings.TrimSuffix(ch.gameHash, lastMove)

	js.Global().Get("location").Set("hash", previousGameMoves)

	return nil
}

func gameHashForHalfMove(g *game.Game, n int) (string, error) {
	if n < 0 && n >= len(g.Positions) {
		return "", errors.New("move no " + strconv.Itoa(n) + " is out of bounds <0, " + strconv.Itoa(len(g.Positions)-1) + ">")
	}

	hash := ""
	for i := 1; i <= n; i++ {
		em, err := encodeMove(g.Positions[i].LastMove)
		if err != nil {
			return "", err
		}
		hash += em
	}

	return hash, nil
}
func (ch *ChessGameModel) HashForInitialHalfMove(n int) (string, error) {
	if err := ch.Validate(); err != nil {
		return "", err
	}
	return gameHashForHalfMove(ch.initialGame, n)
}
func (ch *ChessGameModel) HashForHalfMove(n int) (string, error) {
	if err := ch.Validate(); err != nil {
		return "", err
	}
	return gameHashForHalfMove(ch.game, n)
}

func (ch *ChessGameModel) UpdateModel(tools *shf.Tools, m *HtmlModel, execSupported bool) error {
	if err := ch.Validate(); err != nil {
		return err
	}

	position := ch.game.Positions[ch.currMoveNo]
	nextMoveState := NMError
	{ // set next move state & update game if next move is legal

		// validate next move
		nms, err := getNextMoveState(position, ch.nextMove)
		if err != nil {
			// this should not happen
			return err
		}

		// next move is valid (a valid move, or waiting to fill some params)
		nextMoveState = nms

		if nextMoveState == NMLegalMove {
			// next move is a legal move, do it
			if err := ch.MakeNextMove(); err != nil {
				return err
			}
			position = ch.game.Positions[ch.currMoveNo]
			nextMoveState = NMWaitFrom
			m.Cover.GameStatus.rebuild(tools)

		}
	}
	// from now on, nextMoveState != NMLegalMove

	if nextMoveState == NMWaitPromote {
		m.Board.PromotionOverlay.Shown = true
	} else {
		m.Board.PromotionOverlay.Shown = false
	}

	thrownOutPieces := ch.gameGc[ch.currMoveNo]
	lastMoveThrowOutPiece := piece.New(piece.NoColor, piece.None)

	{ // update last move thorwn out piece
		if ch.currMoveNo > 0 {
			// derive last move throw out piece from curren and previous move throwouts
			if added := addedThrownOuts(ch.gameGc[ch.currMoveNo-1], thrownOutPieces); len(added) > 1 {
				// should not happen
				return errors.New("more thrownouts added between previous and curren move")
			} else if len(added) == 1 {
				for p, c := range added {
					if c > 1 {
						// should not happen
						return errors.New("more thrownouts added between previous and curren move")
					}
					lastMoveThrowOutPiece = p
				}
			}
		}
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

					if ch.game.Status()&(game.WhiteWon|game.BlackWon) > 0 {
						// game ended with whith someone winning, has to be check mate
						m.Board.Grid.Squares[i].Markers.Mate = true
					}
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

		// update color in promotion overlay
		m.Board.PromotionOverlay.Color = position.ActiveColor
	}

	{ // update thrown out pieces
		// first clear all previous html thrownouts
		containers := map[piece.Color]*ThrownOutsContainer{piece.White: m.ThrownOuts.White, piece.Black: m.ThrownOuts.Black}
		for color, container := range containers {
			for _, pt := range thrownOutPiecesOrderType {
				pce := piece.New(color, pt)

				// update html thrown out pieces count
				container.PieceCount[pt] = int(thrownOutPieces[pce])
			}
			container.LastMoveThrowOut = piece.None
			if lastMoveThrowOutPiece.Color == color {
				container.LastMoveThrowOut = lastMoveThrowOutPiece.Type
			}
		}
	}

	{ // update status & notification
		m.Cover.GameStatus.Header.Icons.White = false
		m.Cover.GameStatus.Header.Icons.Black = false
		if st := ch.game.Status(); st != game.InProgress { // game ended

			m.Cover.GameStatus.Header.Message.Text = st.String()
			if st&game.Draw != 0 {
				// game ended in draw
				m.Cover.GameStatus.Header.Icons.White = true
				m.Cover.GameStatus.Header.Icons.Black = true
			} else if st&game.WhiteWon != 0 {
				// white wins
				m.Cover.GameStatus.Header.Icons.White = true
			} else if st&game.BlackWon != 0 {
				// black wins
				m.Cover.GameStatus.Header.Icons.Black = true
			}
		} else {
			// game in progress
			m.Cover.GameStatus.Header.Message.Text = position.ActiveColor.String() + " player is on the move"
			if position.ActiveColor == piece.White {
				// white moves
				m.Cover.GameStatus.Header.Icons.White = true
			} else {
				// black moves
				m.Cover.GameStatus.Header.Icons.Black = true
			}
		}
	}

	{ // update move status
		m.Cover.MoveStatus.Link.MoveHash = ch.gameHash

		if position.LastMove != move.Null {
			m.Cover.MoveStatus.Undo.Get("classList").Call("remove", "hidden")
		} else {
			m.Cover.MoveStatus.Undo.Get("classList").Call("add", "hidden")
		}
	}

	{ // update handlers
		{ // board grid squares
			for _, sq := range m.Board.Grid.Squares {
				// square hack to be able to use 'sq' in anonymous functions
				sq := sq

				// remove square events
				if err := tools.ClickRemove(sq.Element); err != nil {
					return err
				}
				if err := tools.DblClickRemove(sq.Element); err != nil {
					return err
				}

				// every empty square gets a double click callback for zen mode toggle
				if sq.Piece.Type == piece.None {
					if err := tools.DblClick(sq.Element, func(_ shf.Event) error {
						// toggle zen mode by adding/removing "zen-mode" class to body
						js.Global().Get("document").Get("body").Get("classList").Call("toggle", "zen-mode")
						//TODO - Do only needed updates.
						return tools.AppUpdate()
					}); err != nil {
						return err
					}
				}

				// square marked as "possible to" gets unique event
				if sq.Markers.ByColor[position.ActiveColor].NextMove.PossibleTo {
					// inspect next move state
					squareNextMove := ch.nextMove
					squareNextMove.Destination = sq.Id
					squareNextMoveState, _ := getNextMoveState(position, squareNextMove)
					if squareNextMoveState != NMLegalMove && squareNextMoveState != NMWaitPromote {
						// should not happen
						return errors.New("square " + sq.Id.String() + " is marked as possible to, but the next move here is not legal move or waiting to promoion")
					}
					// every moving player possible to move gets event
					if err := tools.Click(sq.Element, func(_ shf.Event) error {
						// set next move to
						ch.nextMove.Destination = sq.Id
						if squareNextMoveState == NMLegalMove {
							// if next move is a legal move, show move status
							m.Cover.MoveStatus.Shown = true
						} else if squareNextMoveState == NMWaitPromote {
							// if next move is a legal move, show move status
							m.Board.PromotionOverlay.Shown = true
						}
						//TODO - Do only needed updates.
						return tools.AppUpdate()
					}); err != nil {
						return err
					}

					// a square can have only one click event, continue to next square
					continue
				}

				// next move from square resets next move
				if ch.nextMove.Source == sq.Id && ch.nextMove.Destination == square.NoSquare && ch.nextMove.Promote == piece.None {
					if err := tools.Click(sq.Element, func(_ shf.Event) error {
						ch.nextMove = move.Null
						m.Cover.MoveStatus.Shown = false
						//TODO - Do only needed updates.
						return tools.AppUpdate()
					}); err != nil {
						return err
					}

					// a square can have only one click event, continue to next square
					continue
				}

				// every moving player figure gets unique event
				if position.ActiveColor == sq.Piece.Color {
					// but only if game is in progress
					if st := ch.game.Status(); st == game.InProgress {
						if err := tools.Click(sq.Element, func(_ shf.Event) error {
							// set next move from
							ch.nextMove.Source = sq.Id
							ch.nextMove.Destination = square.NoSquare
							ch.nextMove.Promote = piece.None

							// hide move status
							m.Cover.MoveStatus.Shown = false
							//TODO - Do only needed updates.
							return tools.AppUpdate()
						}); err != nil {
							return err
						}

						// a square can have only one click event, continue to next square
						continue
					}
				}

				// every empty square or every oponent piece resets next move
				if sq.Piece.Type == piece.None || sq.Piece.Color == complementColor(position.ActiveColor) {
					if err := tools.Click(sq.Element, func(_ shf.Event) error {
						ch.nextMove = move.Null
						m.Cover.MoveStatus.Shown = false
						//TODO - Do only needed updates.
						return tools.AppUpdate()
					}); err != nil {
						return err
					}

					// a square can have only one click event, continue to next square
					continue
				}

			}

			if nextMoveState == NMWaitFrom {

				// every empty square toggles move status
				for _, sq := range m.Board.Grid.Squares {
					if sq.Piece.Type == piece.None {
						if err := tools.Click(sq.Element, func(_ shf.Event) error {
							m.Cover.MoveStatus.Shown = !m.Cover.MoveStatus.Shown
							//TODO - Do only needed updates.
							return tools.AppUpdate()
						}); err != nil {
							return err
						}
					}
				}

				// last move gets some events
				if position.LastMove != move.Null {
					if execSupported {
						// last move to square gets copy to clipboard
						if err := tools.Click(m.Board.Grid.Squares[int(position.LastMove.To())].Element, func(_ shf.Event) error {
							if err := m.CopyGameURLToClipboard(); err != nil {
								return err
							}
							m.Notification.TimedMessage(
								tools,
								5*time.Second,
								"Game URL was copied to clipboard",
								"",
							)
							m.Cover.MoveStatus.Shown = false
							//TODO - Do only needed updates.
							return tools.AppUpdate()
						}); err != nil {
							return err
						}
					} else {
						// copy is not supported, just show move status
						if err := tools.Click(m.Board.Grid.Squares[int(position.LastMove.To())].Element, func(_ shf.Event) error {
							m.Cover.MoveStatus.Shown = true
							//TODO - Do only needed updates.
							return tools.AppUpdate()
						}); err != nil {
							return err
						}
					}

					// last move from square gets back one move
					if err := tools.Click(m.Board.Grid.Squares[int(position.LastMove.From())].Element, func(_ shf.Event) error {
						if err := ch.BackToPreviousMove(); err != nil {
							return err
						}
						m.Cover.GameStatus.rebuild(tools)
						//TODO - Do only needed updates.
						return tools.AppUpdate()
					}); err != nil {
						return err
					}
				} else {
					// The copy to clipboard can be found in menu.
				}
			}

		}
	}

	return nil
}

type Model struct {
	ChessGame *ChessGameModel
	Html      *HtmlModel

	rotationSupported bool
	execSupported     bool
}

func (m *Model) showEndGameNotification(tools *shf.Tools) error {
	newGameButton := tools.CreateElement("button")
	newGameButton.Set("textContent", "new game")
	if err := tools.Click(newGameButton, func(_ shf.Event) error {
		if err := m.ChessGame.UpdateToHash(""); err != nil {
			return err
		}
		m.ChessGame.initialGame = m.ChessGame.game
		m.ChessGame.initialPgn = m.ChessGame.pgn
		m.Html.Notification.Shown = false
		js.Global().Get("location").Set("hash", "")
		m.RotateBoardForPlayer()
		m.Html.Cover.GameStatus.rebuild(tools)
		//TODO - Do only needed updates.
		return tools.AppUpdate()
	}); err != nil {
		// if there is an error creating event for button, simply do not show it
		newGameButton = nil
	}
	exportButton := tools.CreateElement("button")
	exportButton.Set("textContent", "export")
	if err := tools.Click(exportButton, func(_ shf.Event) error {
		m.refreshExportOutputData()
		m.Html.Notification.Shown = false
		m.Html.Export.Shown = true
		//TODO - Do only needed updates.
		return tools.AppUpdate()
	}); err != nil {
		// if there is an error creating event for button, simply do not show it
		exportButton = nil
	}
	closeButton := tools.CreateElement("button")
	closeButton.Set("textContent", "close")
	if err := tools.Click(closeButton, func(_ shf.Event) error {
		m.Html.Notification.Shown = false
		//TODO - Do only needed updates.
		return tools.AppUpdate()
	}); err != nil {
		// if there is an error creating event for button, simply do not show it
		closeButton = nil
	}
	m.Html.Notification.Message(
		m.ChessGame.game.Status().String(),
		"tip: also click anywhere outside to close this notification",
		newGameButton, exportButton, closeButton,
	)
	return nil
}

func (m *Model) refreshExportOutputData() {
	gs := m.ChessGame.game.Status()
	if gs == game.InProgress {
		m.Html.Export.Input.Result.Selected = "*"
		m.Html.Export.Input.Result.Disabled = false
	}
	if gs&game.WhiteWon != 0 {
		m.Html.Export.Input.Result.Selected = "1-0"
		m.Html.Export.Input.Result.Disabled = true
	}
	if gs&game.BlackWon != 0 {
		m.Html.Export.Input.Result.Selected = "0-1"
		m.Html.Export.Input.Result.Disabled = true
	}
	if gs&game.Draw != 0 {
		m.Html.Export.Input.Result.Selected = "1/2-1/2"
		m.Html.Export.Input.Result.Disabled = true
	}
	// http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm
	// 8.1.1.3: The Date tag
	//
	// The Date tag value gives the starting date for the game. (Note: this is not necessarily the same as the starting date for the event.) The date is given with respect to the local time of the site given in the Event tag. The Date tag value field always uses a standard ten character format: "YYYY.MM.DD". The first four characters are digits that give the year, the next character is a period, the next two characters are digits that give the month, the next character is a period, and the final two characters are digits that give the day of the month. If the any of the digit fields are not known, then question marks are used in place of the digits.
	//
	// Examples:
	//
	// [Date "1992.08.31"]
	// [Date "1993.??.??"]
	// [Date "2001.01.01"]
	// TODO - Add hints to fields?
	// Fill current date if empty.
	if m.Html.Export.Input.Date.Input.Get("value").String() == "" {
		m.Html.Export.Input.Date.Input.Set("value", time.Now().Format("2006.01.02"))
	}
	//TODO - Add event tag? - http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c8.1.1
	//TODO - Add termination tag? - http://www.saremba.de/chessgml/standards/pgn/pgn-complete.htm#c9.8.1
	m.Html.Export.Output.PGN = m.ChessGame.pgn
}

func (m *Model) Init(tools *shf.Tools) error {
	if m.ChessGame == nil {
		m.ChessGame, _ = NewGame(js.Global().Get("location").Get("hash").String())
	}

	if err := tools.HashChange(func(e shf.HashChangeEvent) error {
		//println("Mode.Init: HashChange")
		// Get game hash.
		gameHash := "#" + m.ChessGame.gameHash
		// Get location hash.
		locationHash := js.Global().Get("location").Get("hash").String()
		//println("Mode.Init: HashChange: gameHash    :", gameHash)
		//println("Mode.Init: HashChange: locationHash:", gameHash)
		if gameHash == locationHash {
			// Equal, do nothing.
			return nil
		}

		// Not equal game & location hashes.
		//js.Global().Call("alert", "not equal game & location hash: "+gameHash+" != "+locationHash)

		// Update game to the location hash.
		if err := m.ChessGame.UpdateToHash(locationHash); err != nil {
			// Location hash is bad, revert document location hash to game hash.
			//TODO - Proper error showing through notification.
			js.Global().Call("alert", err.Error())
			js.Global().Get("location").Set("hash", m.ChessGame.gameHash)
			return nil
		}

		m.Html.Cover.GameStatus.rebuild(tools)
		// Close move status after game is updated.
		m.Html.Cover.MoveStatus.Shown = false

		return tools.AppUpdate()
	}); err != nil {
		return err
	}

	if m.Html == nil {
		m.Html = &HtmlModel{}

		// Initialize the html model.
		if err := tools.Initialize(m.Html); err != nil {
			return err
		}

		// Add references between elements, where needed.
		m.Html.Cover.GameStatus.Control.refGame = m.ChessGame
		m.Html.Cover.GameStatus.Moves.refGame = m.ChessGame
		m.Html.Cover.GameStatus.Moves.refModel = m.Html

		if !m.rotationSupported {
			m.Html.Rotated180deg = false
			m.Html.Board.Edgings.BottomLeft.Disable()
			m.Html.Board.Edgings.TopRight.Disable()
		} else {
			if err := tools.Click(m.Html.Cover.GameStatus.Header.Element, func(_ shf.Event) error {
				// If game ended, notify the player.
				if st := m.ChessGame.game.Status(); st != game.InProgress {
					if err := m.showEndGameNotification(tools); err != nil {
						return err
					}
				}

				m.RotateBoardForPlayer()
				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}
			if err := tools.Click(m.Html.Board.Edgings.BottomLeft.Element, func(_ shf.Event) error {
				m.RotateBoard()

				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}
			m.Html.Board.Edgings.BottomLeft.Enable()
			if err := tools.Click(m.Html.Board.Edgings.TopRight.Element, func(_ shf.Event) error {
				m.RotateBoard()

				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}
			m.Html.Board.Edgings.TopRight.Enable()
		}

		if !m.execSupported {
			m.Html.Cover.MoveStatus.Link.Copy.Shown = false
			m.Html.Export.Output.Copy.Shown = false
		} else {
			m.Html.Cover.MoveStatus.Link.Copy.Shown = true
			if err := tools.Click(m.Html.Cover.MoveStatus.Link.Copy.Element, func(_ shf.Event) error {
				if err := m.Html.CopyGameURLToClipboard(); err != nil {
					return err
				}
				m.Html.Notification.Message(
					"game URL was copied to clipboard",
					"tip: click on last move piece to copy",
				)
				m.Html.Cover.MoveStatus.Shown = false
				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}

			if err := tools.Click(m.Html.Export.Output.Copy, func(_ shf.Event) error {
				if err := m.Html.CopyExportOutputToClipboard(); err != nil {
					return err
				}
				m.Html.Notification.TimedMessage(
					tools,
					5*time.Second,
					"Game PGN was copied to clipboard",
					"",
				)
				//m.Html.Export.Shown = false
				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}
		}
	}

	{ // add click events for html header & footer

		newGameButton := tools.CreateElement("button")
		newGameButton.Set("textContent", "new game")
		if err := tools.Click(newGameButton, func(_ shf.Event) error {
			if err := m.ChessGame.UpdateToHash(""); err != nil {
				return err
			}
			m.ChessGame.initialGame = m.ChessGame.game
			m.ChessGame.initialPgn = m.ChessGame.pgn
			m.Html.Notification.Shown = false
			js.Global().Get("location").Set("hash", "")
			m.RotateBoardForPlayer()
			m.Html.Cover.GameStatus.rebuild(tools)
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			// if there is an error creating event for button, simply do not show it
			newGameButton = nil
		}

		copyLinkButton := shf.Element(nil)
		if m.execSupported {
			copyLinkButton = tools.CreateElement("button")
			copyLinkButton.Set("textContent", "copy link")
			if err := tools.Click(copyLinkButton, func(e shf.Event) error {
				e.Call("stopPropagation")

				if err := m.Html.CopyGameURLToClipboard(); err != nil {
					return err
				}
				m.Html.Cover.MoveStatus.Shown = false
				m.Html.Notification.TimedMessage(
					tools,
					5*time.Second,
					"game URL was copied to clipboard",
					"tip: click on last move piece to copy",
				)
				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				// if there is an error creating event for button, simply do not show it
				copyLinkButton = nil
			}
		}

		zenModeButton := tools.CreateElement("button")
		zenModeButton.Set("textContent", "toggle zen mode")
		if err := tools.Click(zenModeButton, func(_ shf.Event) error {
			m.Html.Notification.Shown = false
			js.Global().Get("document").Get("body").Get("classList").Call("toggle", "zen-mode")
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			// if there is an error creating event for button, simply do not show it
			zenModeButton = nil
		}

		exportButton := tools.CreateElement("button")
		exportButton.Set("textContent", "export game")
		if err := tools.Click(exportButton, func(e shf.Event) error {
			e.Call("stopPropagation")

			m.refreshExportOutputData()
			m.Html.Notification.Shown = false
			m.Html.Export.Shown = true
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			// if there is an error creating event for button, simply do not show it
			exportButton = nil
		}

		if err := tools.Click(m.Html.Header.Element, func(_ shf.Event) error {
			m.Html.Notification.Message(
				"Quick actions",
				"tip: double click on empty square to toggle zen mode",
				newGameButton,
				copyLinkButton,
				zenModeButton,
				exportButton,
			)
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}

	{ // add promotion events to promotion overlay
		if err := tools.Click(m.Html.Board.PromotionOverlay.Element, func(_ shf.Event) error {
			m.ChessGame.nextMove.Promote = piece.None
			m.ChessGame.nextMove.Destination = square.NoSquare
			m.Html.Board.PromotionOverlay.Shown = false
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}

		for _, p := range m.Html.Board.PromotionOverlay.Pieces {
			promotionPiece := p
			if err := tools.Click(p.Element, func(_ shf.Event) error {
				m.ChessGame.nextMove.Promote = promotionPiece.Piece.Type
				m.Html.Board.PromotionOverlay.Shown = false
				m.Html.Cover.MoveStatus.Shown = true
				//TODO - Do only needed updates.
				return tools.AppUpdate()
			}); err != nil {
				return err
			}
		}
	}

	{ // add back event for move-status
		if err := tools.Click(m.Html.Cover.MoveStatus.Undo, func(_ shf.Event) error {
			if err := m.ChessGame.BackToPreviousMove(); err != nil {
				return err
			}

			m.Html.Cover.MoveStatus.Shown = false
			m.Html.Cover.GameStatus.rebuild(tools)
			//TODO - Do only needed updates.
			return tools.AppUpdate()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) Update(tools *shf.Tools) error {

	if m == nil {
		return errors.New("Model is nil")
	}

	{ // Update html model from chess game.
		err := m.ChessGame.UpdateModel(tools, m.Html, m.execSupported)
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
	if app.ChessGame.game.ActiveColor() == piece.Black {
		app.Html.Rotated180deg = true
	}
}

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

func getNextMoveState(p *position.Position, m move.Move) (int, error) {
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
	shf.Element
	Child   *Child
}
func (this *Model) Init(tools *app.Tools) error {
	if this.Child == nil {
		this.Child = &Child{}
		if err := tools.Initialize(this.Child); err != nil {
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
