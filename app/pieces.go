package app

import (
	"strings"

	"github.com/andrewbackes/chess/piece"
)

var playablePiecesType = []piece.Type{piece.Pawn, piece.Rook, piece.Knight, piece.Bishop, piece.Queen, piece.King}
var promotablePiecesType = []piece.Type{piece.Rook, piece.Knight, piece.Bishop, piece.Queen}
var thrownOutPiecesOrderType = []piece.Type{piece.Pawn, piece.Rook, piece.Knight, piece.Bishop, piece.Queen}

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
