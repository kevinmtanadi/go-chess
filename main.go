package main

import (
	"fmt"
	"go-chess/chess"
	"sync"
	"time"
)

// func main() {
// 	start := time.Now()
// 	cboard := chess.NewChessboard()
// 	cboard.PrintBoard()

// 	fmt.Println("")

// 	possibleBoard := []chess.Chessboard{*cboard}

// 	for turn := 1; turn < 5; turn++ {
// 		posBoard := []chess.Chessboard{}
// 		for _, p := range possibleBoard {
// 			moves := p.GenerateMoves()

// 			for _, m := range moves {
// 				newBoard := p.Copy()
// 				newBoard.Move(m)
// 				newBoard.SwitchTurn()
// 				posBoard = append(posBoard, *newBoard)
// 			}

// 			possibleBoard = posBoard
// 		}

// 		fmt.Println("Round ", turn, " - Possible moves: ", len(possibleBoard))
// 	}

// 	fmt.Println("Time elapsed to calculate 5 rounds: ", time.Since(start))

// }

func main() {
	start := time.Now()
	cboard := chess.NewChessboard()
	cboard.PrintBoard()

	fmt.Println("")

	possibleBoard := []*chess.Chessboard{cboard}

	for turn := 1; turn < 10; turn++ {
		posBoard := []*chess.Chessboard{}
		moveCh := make(chan *chess.Chessboard)
		var wg sync.WaitGroup

		for _, p := range possibleBoard {
			wg.Add(1)
			go func(board *chess.Chessboard) {
				defer wg.Done()

				moves := board.GenerateMoves()

				for _, m := range moves {
					newBoard := board.Copy()
					newBoard.Point[m.StartPosition].HasMoved = true
					newBoard.Move(m)
					newBoard.SwitchTurn()
					moveCh <- newBoard
				}
			}(p)
		}

		go func() {
			wg.Wait()
			close(moveCh)
		}()

		for newBoard := range moveCh {
			posBoard = append(posBoard, newBoard)
		}

		possibleBoard = posBoard

		fmt.Println("Round ", turn, " - Possible moves: ", len(possibleBoard))
	}

	fmt.Println("Time elapsed to calculate 5 rounds: ", time.Since(start))
}
