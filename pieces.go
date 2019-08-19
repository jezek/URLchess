package main

import (
	"URLchess/shf"
	"strings"

	"github.com/andrewbackes/chess/piece"
)

var playablePiecesType = []piece.Type{piece.Pawn, piece.Rook, piece.Knight, piece.Bishop, piece.Queen, piece.King}
var promotablePiecesType = []piece.Type{piece.Rook, piece.Knight, piece.Bishop, piece.Queen}
var thrownOutPiecesOrderType = []piece.Type{piece.Pawn, piece.Knight, piece.Bishop, piece.Rook, piece.Queen}

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

func pieceElement(p piece.Piece) shf.Element {
	elm := shf.CreateElement("span")
	elm.Get("classList").Call("add", "piece")
	if p.Color != piece.NoColor {
		elm.Get("classList").Call("add", strings.ToLower(p.Color.String()))
	}
	if p.Type != piece.None {
		elm.Get("classList").Call("add", pieceTypesToName[p.Type])
	}
	elm.Set("textContent", p.Figurine())
	return elm
}

func complementColor(c piece.Color) piece.Color {
	if c == piece.White || c == piece.Black {
		return piece.Colors[(int(c)+1)%2]
	}
	return piece.NoColor
}
