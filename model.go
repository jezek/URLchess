package main

import (
	"URLchess/shf"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"github.com/gopherjs/gopherjs/js"
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

		f.Call("appendChild", linkURLchess.Object())
		f.Call("appendChild", tools.CreateTextNode(" v"+Version+" by jEzEk. Source on "))
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
		s.piece = pieceElement(piece.New(piece.NoColor, piece.None))
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

type PromotionPiece struct {
	shf.Element
	Piece piece.Piece
}

func (p *PromotionPiece) RedrawElement() {
	if p.Element == nil {
		return
	}

	p.Element.Set("innerHTML", "")
	p.Element.Call("appendChild", pieceElement(p.Piece).Object())
}
func (p *PromotionPiece) Init(tools *shf.Tools) error {
	if p.Element == nil {
		p.Element = tools.CreateElement("span")
		p.Element.Set("id", "promote-to-"+pieceTypesToName[p.Piece.Type])

		p.RedrawElement()
	}
	return nil
}
func (p *PromotionPiece) Update(tools *shf.Tools) error {
	if p == nil {
		return errors.New("PromotionPiece is nil")
	}
	p.RedrawElement()
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
			if err := tools.Update(pp); err != nil {
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
			div.Call("appendChild", pieceElement(piece.New(c.Color, pieceType)).Object())

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
		sm.Call("appendChild", pieceElement(piece.New(piece.White, piece.King)).Object())
	}
	if sm.Black {
		sm.Call("appendChild", pieceElement(piece.New(piece.Black, piece.King)).Object())
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
		if err := tools.Update(gs.Icons); err != nil {
			return err
		}
	}
	if gs.Message == nil {
		gs.Message = &StatusText{}
		if err := tools.Update(gs.Message); err != nil {
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

type StatusMoves struct {
	shf.Element
	GameMovesCount, GameMoveNo uint
}

func (si *StatusMoves) Init(tools *shf.Tools) error {
	if si.Element == nil {
		si.Element = tools.CreateElement("div")
		si.Set("id", "game-status-moves")
	}
	return nil
}
func (si *StatusMoves) Update(tools *shf.Tools) error {
	if si == nil {
		return errors.New("StatusMoves is nil")
	}

	si.Set("innerHTML", "")
	{ // game moves count
		gmc := tools.CreateElement("div")
		gmc.Get("classList").Call("add", "moves-count")
		gmc.Set("textContent", "game moves")

		{ // inner span
			span := tools.CreateElement("span")
			span.Get("classList").Call("add", "moves-count")
			span.Set("textContent", strconv.Itoa(int(si.GameMovesCount)))
			gmc.Call("appendChild", span.Object())
		}
		si.Call("appendChild", gmc.Object())
	}

	return nil
}

type ModelGameStatus struct {
	shf.Element
	Header *StatusHeader
	Moves  *StatusMoves
}

func (gs *ModelGameStatus) Init(tools *shf.Tools) error {
	if gs.Header == nil {
		gs.Header = &StatusHeader{}
		if err := tools.Update(gs.Header); err != nil {
			return err
		}
	}
	if gs.Moves == nil {
		gs.Moves = &StatusMoves{}
		if err := tools.Update(gs.Moves); err != nil {
			return err
		}
	}

	if gs.Element == nil {
		gs.Element = tools.CreateElement("div")
		gs.Set("id", "game-status")
		gs.Call("appendChild", gs.Header.Element.Object())
		gs.Call("appendChild", gs.Moves.Element.Object())
	}
	return nil
}
func (gs *ModelGameStatus) Update(tools *shf.Tools) error {
	if gs == nil {
		return errors.New("ModelGameStatus is nil")
	}

	return tools.Update(gs.Header, gs.Moves)
}

type LinkCopy struct {
	shf.Element
	Shown bool
}

func (this *LinkCopy) Init(tools *shf.Tools) error {
	if this.Element == nil {
		this.Element = tools.CreateElement("button")
		this.Set("textContent", "Copy to clipboard")

	}
	return nil
}
func (this *LinkCopy) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("LinkCopy is nil")
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

	input shf.Element
	Copy  *LinkCopy
}

func (this *MoveStatusLink) GetURL() string {
	hash := "#" + strings.TrimPrefix(this.MoveHash, "#")
	return strings.Split(js.Global.Get("location").String(), "#")[0] + hash
}
func (this *MoveStatusLink) Init(tools *shf.Tools) error {
	if this.input == nil {
		this.input = tools.CreateElement("input")
		this.input.Set("type", "text")
		this.input.Call("setAttribute", "readonly", "readonly")
	}

	if this.Copy == nil {
		this.Copy = &LinkCopy{}
		if err := tools.Update(this.Copy); err != nil {
			return err
		}
	}

	if this.Element == nil {
		this.Element = tools.CreateElement("div")
		this.Get("classList").Call("add", "link")

		this.Call("appendChild", this.input.Object())
		this.Call("appendChild", this.Copy.Object())
		this.Call("appendChild", tools.CreateTextNode("This URL link represents the state of current chess game. You can copy it and store it or send it."))

	}
	return nil
}
func (this *MoveStatusLink) Update(tools *shf.Tools) error {
	if this == nil {
		return errors.New("MoveStatusLink is nil")
	}

	this.input.Set("value", this.GetURL())

	return tools.Update(this.Copy)
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
		if err := tools.Update(this.Link); err != nil {
			return err
		}
	}

	if this.Undo == nil {
		this.Undo = tools.CreateElement("button")
		this.Undo.Set("textContent", "Back")
		// Model sets Event
		// Model.ChessGame.UpdateModel sets visibility
	}

	if this.Close == nil {
		this.Close = tools.CreateElement("button")
		this.Close.Set("textContent", "Close")
		if err := tools.Click(this.Close, func(_ shf.Event) error {
			this.Shown = false
			return nil
		}); err != nil {
			return err
		}
	}

	if this.tip == nil {
		this.tip = tools.CreateElement("div")
		this.tip.Set("className", "tip")
		if err := tools.Click(this.tip, func(_ shf.Event) error {
			return nil
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

type ModelCover struct {
	shf.Element
	GameStatus *ModelGameStatus
	MoveStatus *ModelMoveStatus
}

func (this *ModelCover) Init(tools *shf.Tools) error {
	if this.GameStatus == nil {
		this.GameStatus = &ModelGameStatus{}
		if err := tools.Update(this.GameStatus); err != nil {
			return err
		}
	}
	if this.MoveStatus == nil {
		this.MoveStatus = &ModelMoveStatus{}
		if err := tools.Update(this.MoveStatus); err != nil {
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

type ModelNotification struct {
	shf.Element
	Shown     bool
	timeoutId int
}

func (n *ModelNotification) cancelTimer() {
	if n.timeoutId == 0 {
		return
	}
	js.Global.Call("clearTimeout", n.timeoutId)
	//js.Global.Call("alert", "timer "+strconv.Itoa(n.timeoutId)+" canceled")
	n.timeoutId = 0
}
func (n *ModelNotification) Message(text, hint string, elements ...shf.Element) {
	n.cancelTimer()

	n.Set("innerHTML", "")

	notification := shf.CreateElement("div")
	notification.Get("classList").Call("add", "notification")
	{ // message
		msg := shf.CreateElement("p")
		msg.Get("classList").Call("add", "message")
		msg.Set("textContent", text)
		notification.Call("appendChild", msg.Object())
	}

	for _, e := range elements {
		if e == nil || e.Object() == nil {
			continue
		}
		notification.Call("appendChild", e.Object())
	}
	if hint != "" { // hint
		hintElm := shf.CreateElement("p")
		hintElm.Get("classList").Call("add", "hint")
		hintElm.Set("textContent", hint)
		notification.Call("appendChild", hintElm.Object())
	}
	n.Call("appendChild", notification.Object())
	n.Shown = true
}
func (n *ModelNotification) TimedMessage(tools *shf.Tools, duration time.Duration, text, hint string, elements ...shf.Element) {
	n.cancelTimer()

	notification := shf.CreateElement("div")
	notification.Get("classList").Call("add", "notification")
	{ // message
		msg := shf.CreateElement("p")
		msg.Get("classList").Call("add", "message")
		msg.Set("textContent", text)
		notification.Call("appendChild", msg.Object())
	}

	for _, e := range elements {
		if e == nil || e.Object() == nil {
			continue
		}
		notification.Call("appendChild", e.Object())
	}
	if hint != "" { // hint
		hintElm := shf.CreateElement("p")
		hintElm.Get("classList").Call("add", "hint")
		hintElm.Set("textContent", hint)
		notification.Call("appendChild", hintElm.Object())
	}

	n.Set("innerHTML", "")
	n.Call("appendChild", notification.Object())

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
		if err := tools.Click(n.Element, func(_ shf.Event) error {
			n.cancelTimer()
			n.Shown = false
			return nil
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
	Notification *ModelNotification
	Footer       *ModelFooter
}

func (h *HtmlModel) Init(tools *shf.Tools) error {
	if h.Header == nil {
		h.Header = &ModelHeader{}
		if err := tools.Update(h.Header); err != nil {
			return err
		}
	}
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
	if h.Cover == nil {
		h.Cover = &ModelCover{}
		if err := tools.Update(h.Cover); err != nil {
			return err
		}
	}
	if h.Notification == nil {
		h.Notification = &ModelNotification{}
		if err := tools.Update(h.Notification); err != nil {
			return err
		}
	}
	if h.Footer == nil {
		h.Footer = &ModelFooter{}
		if err := tools.Update(h.Footer); err != nil {
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

	return tools.Update(h.Header, h.Board, h.ThrownOuts, h.Cover, h.Notification, h.Footer)
}

func (h *HtmlModel) RotateBoard() func(shf.Event) error {
	return func(_ shf.Event) error {
		h.Rotated180deg = !h.Rotated180deg
		return nil
	}
}
func (h *HtmlModel) CopyGameURLToClipboard() error {
	positionX := js.Global.Get("pageXOffset")
	positionY := js.Global.Get("pageYOffset")

	temporaryShow := false
	if !h.Cover.MoveStatus.Shown {
		temporaryShow = true
	}
	if temporaryShow {
		h.Cover.MoveStatus.Element.Get("classList").Call("remove", "hidden")
	}
	h.Cover.MoveStatus.Link.input.Call("focus")
	h.Cover.MoveStatus.Link.input.Call("setSelectionRange", 0, h.Cover.MoveStatus.Link.input.Get("value").Get("length"))
	js.Global.Get("document").Call("execCommand", "Copy")
	h.Cover.MoveStatus.Link.input.Call("blur")
	if temporaryShow {
		h.Cover.MoveStatus.Element.Get("classList").Call("add", "hidden")
	}

	js.Global.Call("scrollTo", positionX, positionY)
	return nil
}

type ThrownOuts map[piece.Piece]uint8
type GameThrownOuts []ThrownOuts
type ChessGame struct {
	gameHash   string
	game       *game.Game
	gameGc     GameThrownOuts
	currMoveNo int
	nextMove   move.Move
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
func NewGame(hash string) (*ChessGame, error) {
	// trim movesString from leading #
	movesString := strings.TrimPrefix(hash, "#")

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

			if _, top := pMakeMove(pbm, move); top.Type != piece.None {
				tos[top] = tos[top] + 1
			}

			gtos[i] = tos
		}

		// prepend one empty throw outs
		gtos = append(GameThrownOuts{ThrownOuts{}}, gtos...)
	}

	return &ChessGame{
		gameHash:   movesString,
		game:       g,
		gameGc:     gtos,
		currMoveNo: len(gtos) - 1,
		nextMove:   move.Null,
	}, nil
}

func (ch *ChessGame) Validate() error {
	if ch == nil {
		return errors.New("ChessGame is nil")
	}
	if ch.currMoveNo < 0 || ch.currMoveNo >= len(ch.game.Positions) {
		return errors.New("current move number is out of bounds")
	}
	if len(ch.game.Positions) != len(ch.gameGc) {
		return errors.New("count of game moves and throuwn outs does not match")
	}
	if _, err := getNextMoveState(ch.game.Positions[ch.currMoveNo], ch.nextMove); err != nil {
		return err
	}
	return nil
}

func (ch *ChessGame) MakeNextMove() error {
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

	// update location hash
	js.Global.Get("location").Set("hash", ch.gameHash)

	return nil
}
func (ch *ChessGame) BackToPreviousMove() error {
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

	js.Global.Get("location").Set("hash", previousGameMoves)

	return nil
}

func (ch *ChessGame) UpdateModel(tools *shf.Tools, m *HtmlModel, execSupported bool) error {
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

				// remove square event for sure
				if err := tools.Click(sq.Element, nil); err != nil {
					return err
				}

				// every square gets a double click callback for zen mode toggle
				if err := tools.DblClick(sq.Element, func(_ shf.Event) error {
					// toggle zen mode by adding/removing "zen-mode" class to body
					js.Global.Get("document").Get("body").Get("classList").Call("toggle", "zen-mode")
					return nil
				}); err != nil {
					return err
				}

				// every empty square or every oponent piece resets next move
				if sq.Piece.Type == piece.None || sq.Piece.Color == complementColor(position.ActiveColor) {
					if err := tools.Click(sq.Element, func(_ shf.Event) error {
						ch.nextMove = move.Null
						m.Cover.MoveStatus.Shown = false
						return nil
					}); err != nil {
						return err
					}
				}

				// every moving player figure gets unique event
				if position.ActiveColor == sq.Piece.Color {
					// bt only if game is in progress
					if st := ch.game.Status(); st == game.InProgress {
						if err := tools.Click(sq.Element, func(event shf.Event) error {
							// set next move from
							ch.nextMove.Source = sq.Id
							ch.nextMove.Destination = square.NoSquare
							ch.nextMove.Promote = piece.None

							// hide move status
							m.Cover.MoveStatus.Shown = false
							return nil
						}); err != nil {
							return err
						}
					}
				}

				// next move from square resets next move
				if ch.nextMove.Source == sq.Id && ch.nextMove.Destination == square.NoSquare && ch.nextMove.Promote == piece.None {
					if err := tools.Click(sq.Element, func(_ shf.Event) error {
						ch.nextMove = move.Null
						m.Cover.MoveStatus.Shown = false
						return nil
					}); err != nil {
						return err
					}
				}

				// square marked as possible to gets unique event
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
						return nil
					}); err != nil {
						return err
					}
				}

			}

			if nextMoveState == NMWaitFrom {

				// every empty square toggles move status
				for _, sq := range m.Board.Grid.Squares {
					if sq.Piece.Type == piece.None {
						if err := tools.Click(sq.Element, func(_ shf.Event) error {
							m.Cover.MoveStatus.Shown = !m.Cover.MoveStatus.Shown
							return nil
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
							return nil
						}); err != nil {
							return err
						}
					} else {
						// copy is not supported, just show move status
						if err := tools.Click(m.Board.Grid.Squares[int(position.LastMove.To())].Element, func(_ shf.Event) error {
							m.Cover.MoveStatus.Shown = true
							return nil
						}); err != nil {
							return err
						}
					}

					// last move from square gets back one move
					if err := tools.Click(m.Board.Grid.Squares[int(position.LastMove.From())].Element, func(_ shf.Event) error {
						if err := ch.BackToPreviousMove(); err != nil {
							return err
						}
						return nil
					}); err != nil {
						return err
					}
				} else {
					//TODO where to put copy to clipboard?
				}
			}

		}
	}

	return nil
}

type Model struct {
	Game *ChessGame
	Html *HtmlModel

	rotationSupported bool
	execSupported     bool
}

func (m *Model) Init(tools *shf.Tools) error {
	if m.Game == nil {
		m.Game, _ = NewGame("")
	}

	if err := tools.HashChange(func(e shf.HashChangeEvent) error {
		// get game hash
		gameHash := "#" + m.Game.gameHash
		// get location hash
		locationHash := js.Global.Get("location").Get("hash").String()
		if gameHash == locationHash {
			// equal, do nothing
			return nil
		}
		// not equal game & location hash
		//js.Global.Call("alert", "not equal game & location hash: "+gameHash+" != "+locationHash)

		// create game grom location hash
		newGame, err := NewGame(locationHash)
		if err != nil {
			// location hash is bad, revert document location hash to game hash
			//TODO proper error showing throug notification
			js.Global.Call("alert", err.Error())
			js.Global.Get("location").Set("hash", m.Game.gameHash)
			return nil
		}

		// set game to location hash game
		m.Game = newGame
		m.Html.Cover.MoveStatus.Shown = false

		return nil
	}); err != nil {
		return err
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
			if err := tools.Click(m.Html.Cover.GameStatus.Header.Element, func(_ shf.Event) error {
				m.RotateBoardForPlayer()
				return nil
			}); err != nil {
				return err
			}
			if err := tools.Click(m.Html.Board.Edgings.BottomLeft.Element, m.Html.RotateBoard()); err != nil {
				return err
			}
			m.Html.Board.Edgings.BottomLeft.Enable()
			if err := tools.Click(m.Html.Board.Edgings.TopRight.Element, m.Html.RotateBoard()); err != nil {
				return err
			}
			m.Html.Board.Edgings.TopRight.Enable()
		}

		if !m.execSupported {
			m.Html.Cover.MoveStatus.Link.Copy.Shown = false
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
				return nil
			}); err != nil {
				return err
			}
		}
	}

	{ // add click events for html header & footer

		newGameButton := tools.CreateElement("button")
		newGameButton.Set("textContent", "new game")
		if err := tools.Click(newGameButton, func(_ shf.Event) error {
			game, err := NewGame("")
			if err != nil {
				return err
			}
			m.Game = game
			m.Html.Notification.Shown = false
			js.Global.Get("location").Set("hash", "")
			m.RotateBoardForPlayer()
			return nil
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
				return nil
			}); err != nil {
				// if there is an error creating event for button, simply do not show it
				copyLinkButton = nil
			}
		}

		zenModeButton := tools.CreateElement("button")
		zenModeButton.Set("textContent", "toggle zen mode")
		if err := tools.Click(zenModeButton, func(_ shf.Event) error {
			m.Html.Notification.Shown = false
			js.Global.Get("document").Get("body").Get("classList").Call("toggle", "zen-mode")
			return nil
		}); err != nil {
			// if there is an error creating event for button, simply do not show it
			zenModeButton = nil
		}

		if err := tools.Click(m.Html.Header.Element, func(_ shf.Event) error {
			m.Html.Notification.Message(
				"Quick actions",
				"tip: double click on empty square to toggle zen mode",
				newGameButton,
				copyLinkButton,
				zenModeButton,
			)
			return nil
		}); err != nil {
			return err
		}
	}
	{ // add promotion events to promotion overlay
		if err := tools.Click(m.Html.Board.PromotionOverlay.Element, func(_ shf.Event) error {
			m.Game.nextMove.Promote = piece.None
			m.Game.nextMove.Destination = square.NoSquare
			m.Html.Board.PromotionOverlay.Shown = false
			return nil
		}); err != nil {
			return err
		}

		for _, p := range m.Html.Board.PromotionOverlay.Pieces {
			promotionPiece := p
			if err := tools.Click(p.Element, func(_ shf.Event) error {
				m.Game.nextMove.Promote = promotionPiece.Piece.Type
				m.Html.Board.PromotionOverlay.Shown = false
				m.Html.Cover.MoveStatus.Shown = true
				return nil
			}); err != nil {
				return err
			}
		}
	}

	{ // add back event for move-status
		if err := tools.Click(m.Html.Cover.MoveStatus.Undo, func(_ shf.Event) error {
			if err := m.Game.BackToPreviousMove(); err != nil {
				return err
			}

			m.Html.Cover.MoveStatus.Shown = false
			return nil
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

	{ // update html model from game
		err := m.Game.UpdateModel(tools, m.Html, m.execSupported)
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
