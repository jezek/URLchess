package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

const encodePosAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var encodePromotionCharToPiece map[byte]piece.Type = map[byte]piece.Type{
	'@': piece.Rook,
	'$': piece.Knight,
	'^': piece.Bishop,
	'*': piece.Queen,
}

var encodePieceToPromotionChar map[piece.Type]byte = func() map[piece.Type]byte {
	res := map[piece.Type]byte{}
	for b, p := range encodePromotionCharToPiece {
		res[p] = b
	}
	return res
}()

func encodeMove(m move.Move) (string, error) {
	res := ""

	if int(m.Source) < 0 || int(m.Source) >= len(encodePosAlphabet) {
		return "", errors.New("Source square integer out of bounds: " + strconv.Itoa(int(m.Source)))
	}
	res += string(encodePosAlphabet[int(m.Source)])

	if int(m.Destination) < 0 || int(m.Destination) >= len(encodePosAlphabet) {
		return "", errors.New("Destination square integer out of bounds: " + strconv.Itoa(int(m.Destination)))
	}
	res += string(encodePosAlphabet[int(m.Destination)])

	if m.Promote != piece.None {
		b, ok := encodePieceToPromotionChar[m.Promote]
		if !ok {
			return "", errors.New("Invalid promotion piece: " + m.Promote.String())
		}
		res += string(b)
	}

	return res, nil
}

func EncodeGame(g *ChessGame) (string, error) {
	res := ""
	for _, position := range g.game.Positions {
		if position.LastMove != move.Null {
			if m, err := encodeMove(position.LastMove); err != nil {
				return "", err
			} else {
				res += m
			}
		}
	}
	return res, nil
}

func DecodeMoves(moves string) ([]move.Move, error) {
	res := []move.Move{}
	if moves == "" {
		return res, nil
	}

	for moves != "" {
		move := move.Move{}

		fromInt := strings.Index(encodePosAlphabet, string(moves[0]))
		if fromInt == -1 {
			return nil, errors.New("Invalid move from position character: " + string(moves[0]))
		}
		moves = moves[1:]

		fromSquare := square.Square(fromInt)
		if fromSquare < 0 || fromSquare > square.LastSquare {
			return nil, errors.New("Invalid move from square integer: " + strconv.Itoa(fromInt))
		}

		move.Source = fromSquare

		if len(moves) == 0 {
			return nil, errors.New("Missing move to position character")
		}

		toInt := strings.Index(encodePosAlphabet, string(moves[0]))
		if toInt == -1 {
			return nil, errors.New("Invalid move to position character: " + string(moves[0]))
		}
		moves = moves[1:]

		toSquare := square.Square(toInt)
		if toSquare < 0 || toSquare > square.LastSquare {
			return nil, errors.New("Invalid move to square integer: " + strconv.Itoa(toInt))
		}

		move.Destination = toSquare

		if len(moves) > 0 {
			if piece, ok := encodePromotionCharToPiece[moves[0]]; ok {
				// next char is promotion character
				moves = moves[1:]
				move.Promote = piece
			}
		}
		res = append(res, move)
	}
	return res, nil
}
