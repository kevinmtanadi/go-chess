package chess

import (
	"fmt"
	"go-chess/helper"
	"strings"

	"github.com/mohae/deepcopy"
)

type Chessboard struct {
	Point []ChessPiece
	Turn  int
}

func NewChessboard() *Chessboard {
	cboard := Chessboard{
		Point: make([]ChessPiece, 64),
		Turn:  0,
	}
	cboard.GenerateFen()

	return &cboard
}

func (cb *Chessboard) PrintBoard() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if cb.Point[i*8+j].Type == "" {
				fmt.Print(". ")
			} else {
				fmt.Print(cb.Point[i*8+j].ToPrintable())
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

func (cb *Chessboard) Move(m Move) {
	cb.Point[m.EndPosition] = cb.Point[m.StartPosition]
	cb.Point[m.StartPosition] = ChessPiece{}
}

func (cb *Chessboard) Copy() *Chessboard {
	return deepcopy.Copy(cb).(*Chessboard)
}

func (cb *Chessboard) SwitchTurn() {
	if cb.Turn == 0 {
		cb.Turn = 1
	} else {
		cb.Turn = 0
	}
}

const (
	Type_PAWN   = "pawn"
	Type_QUEEN  = "queen"
	Type_BISHOP = "bishop"
	Type_ROOK   = "rook"
	Type_KNIGHT = "knight"
	Type_KING   = "king"

	FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func (cb *Chessboard) GenerateFen(fenString ...string) {
	var fen string
	if len(fenString) <= 0 {
		fen = FEN
	} else {
		fen = fenString[0]
	}

	parts := strings.Split(fen, " ")
	position := parts[0]

	ranks := strings.Split(position, "/")
	cb.Point = make([]ChessPiece, 64)

	for rankIndex, rank := range ranks {
		index := 0
		for _, char := range rank {
			if index >= 8 {
				break
			}

			if char >= '1' && char <= '8' {
				index += int(char - '0')
			} else {
				var pieceValue ChessPiece
				position := rankIndex*8 + index
				switch char {
				case 'p':
					pieceValue = ChessPiece{Type: Type_PAWN, Value: 1, Color: -1, Position: position} // Black pawn
				case 'b':
					pieceValue = ChessPiece{Type: Type_BISHOP, Value: 3, Color: -1, Position: position}
				case 'r':
					pieceValue = ChessPiece{Type: Type_ROOK, Value: 5, Color: -1, Position: position}
				case 'n':
					pieceValue = ChessPiece{Type: Type_KNIGHT, Value: 4, Color: -1, Position: position}
				case 'q':
					pieceValue = ChessPiece{Type: Type_QUEEN, Value: 9, Color: -1, Position: position}
				case 'k':
					pieceValue = ChessPiece{Type: Type_KING, Value: 100, Color: -1, Position: position}
				case 'P':
					pieceValue = ChessPiece{Type: Type_PAWN, Value: 1, Color: 1, Position: position} // White pawn
				case 'B':
					pieceValue = ChessPiece{Type: Type_BISHOP, Value: 3, Color: 1, Position: position}
				case 'R':
					pieceValue = ChessPiece{Type: Type_ROOK, Value: 5, Color: 1, Position: position}
				case 'N':
					pieceValue = ChessPiece{Type: Type_KNIGHT, Value: 4, Color: 1, Position: position}
				case 'Q':
					pieceValue = ChessPiece{Type: Type_QUEEN, Value: 9, Color: 1, Position: position}
				case 'K':
					pieceValue = ChessPiece{Type: Type_KING, Value: 100, Color: 1, Position: position}
				default:
					panic("Invalid FEN string")
				}

				cb.Point[rankIndex*8+index] = pieceValue
				index++
			}
		}
	}
}

type Move struct {
	StartPosition int
	EndPosition   int
}

type Piece interface {
	GenerateMoves(board Chessboard) []Move
}

type ChessPiece struct {
	Type     string
	Value    int
	Color    int // black = -1 white = 1
	Position int
	HasMoved bool
}

func (cp ChessPiece) ToPrintable() string {
	switch cp.Type {
	case Type_KING:
		return AssertColor("k", cp.Color)
	case Type_QUEEN:
		return AssertColor("q", cp.Color)
	case Type_ROOK:
		return AssertColor("r", cp.Color)
	case Type_KNIGHT:
		return AssertColor("n", cp.Color)
	case Type_BISHOP:
		return AssertColor("b", cp.Color)
	case Type_PAWN:
		return AssertColor("p", cp.Color)
	}

	return ""
}

func AssertColor(piece string, color int) string {
	if color == -1 {
		return strings.ToLower(piece)
	}

	return strings.ToUpper(piece)
}

func (cb Chessboard) GenerateMoves() []Move {
	moves := []Move{}

	for _, p := range cb.Point {
		// white turn
		if cb.Turn == 0 && cb.Point[p.Position].Color != 1 {
			continue
		}

		if cb.Turn == 1 && cb.Point[p.Position].Color != -1 {
			continue
		}

		switch p.Type {
		case Type_KING:
			moves = append(moves, p.GenerateKingMoves(cb)...)
		case Type_QUEEN:
			moves = append(moves, p.GenerateQueenMoves(cb)...)
		case Type_ROOK:
			moves = append(moves, p.GenerateRookMoves(cb)...)
		case Type_KNIGHT:
			moves = append(moves, p.GenerateKnightMoves(cb)...)
		case Type_BISHOP:
			moves = append(moves, p.GenerateBishopMoves(cb)...)
		case Type_PAWN:
			moves = append(moves, p.GeneratePawnMoves(cb)...)
		default:
			continue
		}
	}
	return moves
}

func (king ChessPiece) GenerateKingMoves(board Chessboard) []Move {
	directions := []int{-8, 8, -1, 1, -9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := king.Position + d

		// Igne positionore thess because they are out of bounds
		if newPosition < 0 || newPosition >= 64 {
			continue
		}

		// Target is teammate
		if board.Point[newPosition].Color == king.Color {
			continue
		}

		moves = append(moves, Move{king.Position, newPosition})
	}

	return moves
}

func (queen ChessPiece) GenerateQueenMoves(board Chessboard) []Move {
	directions := []int{-8, 8, -1, 1, -9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := queen.Position

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if board.Point[newPosition].Color != queen.Color {
				moves = append(moves, Move{queen.Position, newPosition})
			}

			// Stop going if queen is blocked
			if &board.Point[newPosition] != nil {
				break
			}

			moves = append(moves, Move{queen.Position, newPosition})
		}
	}

	return moves
}

func (rook ChessPiece) GenerateRookMoves(board Chessboard) []Move {
	directions := []int{-8, 8, -1, 1}
	moves := []Move{}

	for _, d := range directions {
		newPosition := rook.Position

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if board.Point[newPosition].Color != rook.Color {
				moves = append(moves, Move{rook.Position, newPosition})
				break
			}

			// Stop going if rook is blocked
			if &board.Point[newPosition] != nil {
				break
			}

			moves = append(moves, Move{rook.Position, newPosition})
		}
	}

	return moves
}

func (bishop ChessPiece) GenerateBishopMoves(board Chessboard) []Move {
	directions := []int{-9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := bishop.Position

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if board.Point[newPosition].Color != bishop.Color {
				moves = append(moves, Move{bishop.Position, newPosition})
				break
			}

			// Stop going if bishop is blocked
			if &board.Point[newPosition] != nil {
				break
			}

			moves = append(moves, Move{bishop.Position, newPosition})
		}
	}

	return moves
}

func (knight ChessPiece) GenerateKnightMoves(board Chessboard) []Move {
	moves := []Move{}

	possibleMovement := []int{-17, -10, -15, -6, 17, 10, 15, 6}

	for _, m := range possibleMovement {
		endPosition := knight.Position + m

		if endPosition < 0 || endPosition >= 64 {
			continue
		}
		if board.Point[endPosition].Color == knight.Color {
			continue
		}

		if board.Point[endPosition].Type != "" {
			moves = append(moves, Move{knight.Position, endPosition})
		}

		sourceRank, sourceFile := knight.Position/8, knight.Position%8
		destRank, destFile := endPosition/8, endPosition%8

		if helper.Abs(sourceRank-destRank) <= 2 && helper.Abs(sourceFile-destFile) <= 2 {
			moves = append(moves, Move{knight.Position, endPosition})
		}
	}

	return moves
}

func (pawn ChessPiece) GeneratePawnMoves(board Chessboard) []Move {
	moves := []Move{}

	var endPosition int
	if pawn.Color == 1 {
		endPosition = pawn.Position - 8
	} else {
		endPosition = pawn.Position + 8
	}

	if (endPosition >= 0 && endPosition < 64) && board.Point[endPosition].Type == "" {
		moves = append(moves, Move{pawn.Position, endPosition})
	}

	if (pawn.Position < 16 && pawn.Color == -1) || (pawn.Position > 49 && pawn.Color == 1) {
		var endPosition int
		if pawn.Color == 1 {
			endPosition = pawn.Position - 16
		} else {
			endPosition = pawn.Position + 16
		}

		if (endPosition >= 0 && endPosition < 64) && board.Point[endPosition].Type == "" {
			moves = append(moves, Move{pawn.Position, endPosition})
		}
		return moves
	}

	var leftAttackPosition int
	if pawn.Color == -1 {
		leftAttackPosition = pawn.Position - 9
	} else {
		leftAttackPosition = pawn.Position + 7
	}

	if (leftAttackPosition > 0 && leftAttackPosition < 64) && &board.Point[leftAttackPosition] != nil &&
		(leftAttackPosition/8-pawn.Position/8) == 1 {
		if board.Point[leftAttackPosition].Color != pawn.Color {
			moves = append(moves, Move{pawn.Position, leftAttackPosition})
		}
	}

	var rightAttackPosition int
	if pawn.Color == -1 {
		rightAttackPosition = pawn.Position - 7
	} else {
		rightAttackPosition = pawn.Position + 9
	}

	if (rightAttackPosition > 0 && rightAttackPosition < 64) && &board.Point[rightAttackPosition] != nil &&
		(rightAttackPosition/8-pawn.Position/8) == -1 {
		if board.Point[rightAttackPosition].Color != pawn.Color {
			moves = append(moves, Move{pawn.Position, rightAttackPosition})
		}
	}

	return moves
}
