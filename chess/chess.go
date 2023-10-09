package chess

import (
	"fmt"
	"go-chess/helper"
	"strings"
)

var Board [64][]int

type Move struct {
	StartPosition int
	EndPosition   int
}

func PrintBoard(b []int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b[i*8+j] == 0 {
				fmt.Print(". ")
			} else {
				fmt.Print(ToPrintable(b[i*8+j]))
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

const (
	Type_WHITE_KING   = 100
	Type_WHITE_QUEEN  = 9
	Type_WHITE_ROOK   = 5
	Type_WHITE_KNIGHT = 4
	Type_WHITE_BISHOP = 3
	Type_WHITE_PAWN   = 1
	Type_BLACK_KING   = -100
	Type_BLACK_QUEEN  = -9
	Type_BLACK_ROOK   = -5
	Type_BLACK_KNIGHT = -4
	Type_BLACK_BISHOP = -3
	Type_BLACK_PAWN   = -1

	FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func ToPrintable(piece int) string {
	switch piece {
	case Type_WHITE_KING:
		return "K"
	case Type_WHITE_QUEEN:
		return "Q"
	case Type_WHITE_ROOK:
		return "R"
	case Type_WHITE_KNIGHT:
		return "N"
	case Type_WHITE_BISHOP:
		return "B"
	case Type_WHITE_PAWN:
		return "P"
	case Type_BLACK_KING:
		return "k"
	case Type_BLACK_QUEEN:
		return "q"
	case Type_BLACK_ROOK:
		return "r"
	case Type_BLACK_KNIGHT:
		return "n"
	case Type_BLACK_BISHOP:
		return "b"
	case Type_BLACK_PAWN:
		return "p"
	}

	return ""
}

func GenerateStartingBoard() []int {
	parts := strings.Split(FEN, " ")
	position := parts[0]

	ranks := strings.Split(position, "/")
	board := make([]int, 64)

	for rankIndex, rank := range ranks {
		index := 0
		for _, char := range rank {
			if index >= 8 {
				break
			}

			if char >= '1' && char <= '8' {
				index += int(char - '0')
			} else {
				var pieceValue int
				switch char {
				case 'p':
					pieceValue = Type_BLACK_PAWN
				case 'b':
					pieceValue = Type_BLACK_BISHOP
				case 'r':
					pieceValue = Type_BLACK_ROOK
				case 'n':
					pieceValue = Type_BLACK_KNIGHT
				case 'q':
					pieceValue = Type_BLACK_QUEEN
				case 'k':
					pieceValue = Type_BLACK_KING
				case 'P':
					pieceValue = Type_WHITE_PAWN
				case 'B':
					pieceValue = Type_WHITE_BISHOP
				case 'R':
					pieceValue = Type_WHITE_ROOK
				case 'N':
					pieceValue = Type_WHITE_KNIGHT
				case 'Q':
					pieceValue = Type_WHITE_QUEEN
				case 'K':
					pieceValue = Type_WHITE_KING
				default:
					panic("Invalid FEN string")
				}

				board[rankIndex*8+index] = pieceValue
				index++
			}
		}
	}

	return board
}

func GenerateMoves(board []int, turn int) []Move {
	var moves []Move
	// White turn
	if turn%2 == 1 {
		for i := 0; i < 64; i++ {
			piece := board[i]

			switch piece {
			case Type_WHITE_KING:
				moves = append(moves, GenerateKingMoves(board, i)...)
			case Type_WHITE_QUEEN:
				moves = append(moves, GenerateQueenMoves(board, i)...)
			case Type_WHITE_ROOK:
				moves = append(moves, GenerateRookMoves(board, i)...)
			case Type_WHITE_KNIGHT:
				moves = append(moves, GenerateKnightMoves(board, i)...)
			case Type_WHITE_BISHOP:
				moves = append(moves, GenerateBishopMoves(board, i)...)
			case Type_WHITE_PAWN:
				moves = append(moves, GeneratePawnMoves(board, i, 1)...)
			}
		}
	} else {
		for i := 0; i < 64; i++ {
			piece := board[i]

			switch piece {
			case Type_BLACK_KING:
				moves = append(moves, GenerateKingMoves(board, i)...)
			case Type_BLACK_QUEEN:
				moves = append(moves, GenerateQueenMoves(board, i)...)
			case Type_BLACK_ROOK:
				moves = append(moves, GenerateRookMoves(board, i)...)
			case Type_BLACK_KNIGHT:
				moves = append(moves, GenerateKnightMoves(board, i)...)
			case Type_BLACK_BISHOP:
				moves = append(moves, GenerateBishopMoves(board, i)...)
			case Type_BLACK_PAWN:
				moves = append(moves, GeneratePawnMoves(board, i, -1)...)
			}
		}
	}

	return moves
}

func GenerateKingMoves(board []int, king int) []Move {
	directions := []int{-8, 8, -1, 1, -9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := king + d

		// Igne positionore thess because they are out of bounds
		if newPosition < 0 || newPosition >= 64 {
			continue
		}

		// Target is teammate
		if IsTeammate(board[king], board[newPosition]) {
			continue
		}

		moves = append(moves, Move{king, newPosition})
	}

	return moves
}

func GenerateQueenMoves(board []int, queen int) []Move {
	directions := []int{-8, 8, -1, 1, -9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := queen

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if !IsTeammate(board[queen], board[newPosition]) {
				moves = append(moves, Move{queen, newPosition})
			}

			// Stop going if queen is blocked
			if board[newPosition] != 0 {
				break
			}

			moves = append(moves, Move{queen, newPosition})
		}
	}

	return moves
}

func GenerateRookMoves(board []int, rook int) []Move {
	directions := []int{-8, 8, -1, 1}
	moves := []Move{}

	for _, d := range directions {
		newPosition := rook

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if !IsTeammate(board[rook], board[newPosition]) {
				moves = append(moves, Move{rook, newPosition})
			}

			// Stop going if rook is blocked
			if board[newPosition] != 0 {
				break
			}

			moves = append(moves, Move{rook, newPosition})
		}
	}

	return moves
}

func GenerateBishopMoves(board []int, bishop int) []Move {
	directions := []int{-9, 9, -7, 7}
	moves := []Move{}

	for _, d := range directions {
		newPosition := bishop

		for {
			newPosition += d
			if newPosition < 0 || newPosition >= 64 {
				break
			}

			// target is enemy
			if !IsTeammate(board[bishop], board[newPosition]) {
				moves = append(moves, Move{bishop, newPosition})
			}

			// Stop going if bishop is blocked
			if board[newPosition] != 0 {
				break
			}

			moves = append(moves, Move{bishop, newPosition})
		}
	}

	return moves
}

func GenerateKnightMoves(board []int, knight int) []Move {
	moves := []Move{}

	possibleMovement := []int{-17, -10, -15, -6, 17, 10, 15, 6}

	for _, m := range possibleMovement {
		endPosition := knight + m

		if endPosition < 0 || endPosition >= 64 {
			continue
		}

		if IsTeammate(board[knight], board[endPosition]) {
			continue
		}

		if board[endPosition] != 0 {
			moves = append(moves, Move{knight, endPosition})
		}

		sourceRank, sourceFile := knight/8, knight%8
		destRank, destFile := endPosition/8, endPosition%8

		if helper.Abs(sourceRank-destRank) <= 2 && helper.Abs(sourceFile-destFile) <= 2 {
			moves = append(moves, Move{knight, endPosition})
		}
	}

	return moves
}

func GeneratePawnMoves(board []int, pawn int, color int) []Move {
	moves := []Move{}

	frontEmpty := false
	var endPosition int
	if color == 1 {
		endPosition = pawn - 8
	} else {
		endPosition = pawn + 8
	}

	if (endPosition >= 0 && endPosition < 64) && board[endPosition] == 0 {
		frontEmpty = true
		moves = append(moves, Move{pawn, endPosition})
	}

	if ((pawn < 16 && color == -1) || (pawn > 47 && color == 1)) && frontEmpty {
		var endPosition int
		if color == 1 {
			endPosition = pawn - 16
		} else {
			endPosition = pawn + 16
		}

		if (endPosition >= 0 && endPosition < 64) && board[endPosition] == 0 {
			moves = append(moves, Move{pawn, endPosition})
		}
		return moves
	}

	var leftAttackPosition int
	if color == -1 {
		leftAttackPosition = pawn - 9
	} else {
		leftAttackPosition = pawn + 7
	}

	if (leftAttackPosition > 0 && leftAttackPosition < 64) && !IsTeammate(board[pawn], board[leftAttackPosition]) &&
		(leftAttackPosition/8-pawn/8) == 1 {
		moves = append(moves, Move{pawn, leftAttackPosition})
	}

	var rightAttackPosition int
	if color == -1 {
		rightAttackPosition = pawn - 7
	} else {
		rightAttackPosition = pawn + 9
	}

	if (rightAttackPosition > 0 && rightAttackPosition < 64) && !IsTeammate(board[pawn], board[rightAttackPosition]) &&
		(rightAttackPosition/8-pawn/8) == -1 {
		moves = append(moves, Move{pawn, rightAttackPosition})
	}

	return moves
}

func IsTeammate(piece1 int, piece2 int) bool {
	if piece1 == 0 || piece2 == 0 {
		return false
	}

	if piece1 > 0 {
		return piece2 > 0
	} else {
		return piece2 < 0
	}
}

func MovePiece(board []int, move Move) {
	board[move.EndPosition] = board[move.StartPosition]
	board[move.StartPosition] = 0
}
