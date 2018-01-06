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

const encodePosAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var encodePromotionCharToPiece map[byte]piece.Type = map[byte]piece.Type{
	'$': piece.Knight,
	'@': piece.Bishop,
	'#': piece.Rook,
	'*': piece.Queen,
}

var piecesString map[piece.Color]map[piece.Type]string = map[piece.Color]map[piece.Type]string{
	piece.White: map[piece.Type]string{
		piece.King:   "♔",
		piece.Queen:  "♕",
		piece.Bishop: "♗",
		piece.Knight: "♘",
		piece.Rook:   "♖",
		piece.Pawn:   "♙",
	},
	piece.Black: map[piece.Type]string{
		piece.King:   "♚",
		piece.Queen:  "♛",
		piece.Bishop: "♝",
		piece.Knight: "♞",
		piece.Rook:   "♜",
		piece.Pawn:   "♟",
	},
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

func decodeMoves(moves string) ([]move.Move, error) {
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

func main() {
	document := js.Global.Get("document")
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
	//document.Call("write", "moves raw string: "+movesString+"<br/>")

	moves, err := decodeMoves(movesString)
	if err != nil {
		document.Call("write", "Error decoding moves: "+err.Error())
		return
	}

	g := game.New()

	{
		var err error
		for i, move := range moves {
			if g.Status() != game.InProgress {
				err = errors.New("Too many moves in url string! " + strconv.Itoa(i+1) + " moves are enough")
				break
			}

			_, merr := g.MakeMove(move)
			if merr != nil {
				err = errors.New("Errorneous move number " + strconv.Itoa(i+1) + ": " + merr.Error())
				break
			}
		}

		flipBoard := false // g.ActiveColor() == piece.Black

		document.Call("write", "<div id=\"board\">")

		document.Call("write", "<div id=\"edging-top\">")
		for i := 0; i < 8; i++ {
			n := i
			if flipBoard {
				n = 7 - i
			}
			document.Call("write", "<div>"+string(rune('a'+n))+"</div>")

		}
		document.Call("write", "</div>") //edging-top

		document.Call("write", "<div id=\"edging-left\">")
		for i := 8; i > 0; i-- {
			n := i
			if flipBoard {
				n = 9 - i
			}
			document.Call("write", "<div>"+strconv.Itoa(n)+"</div>")

		}
		document.Call("write", "</div>") //edging-left

		{
			document.Call("write", "<div id=\"grid\">")
			squareTones := []string{"light-square", "dark-square"}
			beg, end, inc := int(63), int(-1), int(-1)
			if flipBoard {
				beg, end, inc = int(0), int(64), int(1)
			}
			for i := beg; i != end; i += inc {
				sq := square.Square(i)
				class := []string{squareTones[(i%8+i/8)%2]}
				lm := g.Position().LastMove
				if lm != move.Null {
					if lm.Source == sq || lm.Destination == sq {
						class = append(class, "last-move")
					}
				}
				pc := g.Position().OnSquare(sq)
				document.Call("write", "<div id=\""+sq.String()+"\" class=\""+strings.Join(class, " ")+"\">"+piecesString[pc.Color][pc.Type]+"</div>")
			}
			document.Call("write", "</div>") //grid
		}

		document.Call("write", "<div id=\"edging-right\">")
		for i := 8; i > 0; i-- {
			n := i
			if flipBoard {
				n = 9 - i
			}
			document.Call("write", "<div>"+strconv.Itoa(n)+"</div>")

		}
		document.Call("write", "</div>") //edging-right

		document.Call("write", "<div id=\"edging-bottom\">")
		for i := 0; i < 8; i++ {
			n := i
			if flipBoard {
				n = 7 - i
			}
			document.Call("write", "<div>"+string(rune('a'+n))+"</div>")

		}
		document.Call("write", "</div>") //edging-bottom

		document.Call("write", "</div>") //board
		//document.Call("write", "<div style=\"margin-bottom:1em;\">black: prnbqk<pre>"+g.String()+"</pre>white: PRNBQK</div>")

		if err != nil {
			document.Call("write", "<div>"+err.Error()+"</div>")
			return
		}
	}

	if g.Status() != game.InProgress {
		document.Call("write", "<div style=\"margin-top: 1em\">Game has ended. "+g.Status().String()+"</div>")
		return
	}
	document.Call("write", "<div id=\"game-status\">")
	document.Call("write", "<p>Player on the move: <span style=\"font-size: 1cm;\">"+piecesString[g.ActiveColor()][piece.King]+"</span></p>")
	document.Call("write", "</div>") //game-status

	document.Call("write", `<div id="next-move">
<p>Your next move: <input id="move-input"/> eg. e2e4 (or e7e8q for promotion to queen) and press [ENTER]</p>
<a id="next-move-link" href=""></a><span id="next-move-error"></span>
</div>`)

	moveInput := js.Global.Get("document").Call("getElementById", "move-input")
	if moveInput == nil {
		document.Call("write", "Next move input element not found")
	} else {
		moveInput.Call(
			"addEventListener",
			"keyup",
			func(event *js.Object) {
				if keycode := event.Get("keyCode").Int(); keycode == 13 {
					nextMoveError := js.Global.Get("document").Call("getElementById", "next-move-error")
					if nextMoveError == nil {
						document.Call("write", "Next move error element not found")
						return
					}
					nextMoveLink := js.Global.Get("document").Call("getElementById", "next-move-link")
					if nextMoveLink == nil {
						nextMoveError.Set("innerHTML", "Next move link element not found")
						return
					}

					nextMoveError.Set("innerHTML", "")
					nextMoveLink.Set("innerHTML", "")
					nextMoveLink.Set("href", "")

					nextMovePCM := strings.TrimSpace(moveInput.Get("value").String())
					if nextMovePCM == "" {
						return
					}

					nextMove := move.Parse(nextMovePCM)
					if nextMove == move.Null {
						nextMoveError.Set("innerHTML", "Next move is not in PCN format")
						return
					}

					if _, ok := g.LegalMoves()[nextMove]; ok == false {
						nextMoveError.Set("innerHTML", "Next move is not a legal move")
						return
					}

					nextMoveString, err := encodeMove(nextMove)
					if err != nil {
						nextMoveError.Set("innerHTML", "Next move encoding error: "+err.Error())
						return
					}

					url := location.Get("origin").String() + location.Get("pathname").String() + "?" + movesString + nextMoveString
					nextMoveLink.Set("innerHTML", url)
					nextMoveLink.Set("href", url)
					nextMoveError.Set("innerHTML", " <- copy this link and send to your oponent")
				}
			},
			false,
		)
		moveInput.Call("focus")
	}
}
