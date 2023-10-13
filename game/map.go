package game

import (
	"fmt"
	"math/rand"
)

const (
	Empty = "e"
	Cell  = "c"
	Wall  = "0"
	Block = "d"
)

const PixelBlockSize = 64

func CreateMap() [][]string {
	rows, cols := 13, 19
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
		for j := range board[i] {
			board[i][j] = Empty
		}
	}

	// Create walls
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i%3 == 0 && j%3 == 0 {
				board[i][j] = Wall
			}
		}
	}

	// Random blocks
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == Empty {
				// Skip corners where players start
				if
				// Player 1 (Top-left corner)
				(i == 0 && (j == 0 || j == 1)) ||
					(i == 1 && j == 0) ||

					// Player 2 (Top-right corner)
					(i == 0 && (j == cols-1 || j == cols-2)) ||
					(i == 1 && j == cols-1) ||

					// Player 3 (Bottom-left corner)
					(i == rows-1 && (j == 0 || j == 1)) ||
					(i == rows-2 && j == 0) ||

					// Player 4 (Bottom-right corner)
					(i == rows-1 && (j == cols-1 || j == cols-2)) ||
					(i == rows-2 && j == cols-1) {
					continue
				}
				if rand.Float32() < 0.8 {
					board[i][j] = Block
				}
			}
		}
	}

	// Set player starting positions
	board[0][0] = "1"
	board[0][cols-1] = "2"
	board[rows-1][0] = "3"
	board[rows-1][cols-1] = "4"

	// Print the board for demonstration
	fmt.Println("Map Created:")
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%2s", cell)
		}
		fmt.Println()
	}

	return board
}

func IsCollision(gameboard [][]string, x, y int) bool {
	rows := len(gameboard)
	if rows == 0 {
		return false
	}
	cols := len(gameboard[0])

	topLeftRow, topLeftCol := y/PixelBlockSize, x/PixelBlockSize
	topRightRow, topRightCol := y/PixelBlockSize, (x+63)/PixelBlockSize
	bottomLeftRow, bottomLeftCol := (y+63)/PixelBlockSize, x/PixelBlockSize
	bottomRightRow, bottomRightCol := (y+63)/PixelBlockSize, (x+63)/PixelBlockSize

	if topLeftRow < 0 || topLeftRow >= rows || topLeftCol < 0 || topLeftCol >= cols ||
		topRightRow < 0 || topRightRow >= rows || topRightCol < 0 || topRightCol >= cols ||
		bottomLeftRow < 0 || bottomLeftRow >= rows || bottomLeftCol < 0 || bottomLeftCol >= cols ||
		bottomRightRow < 0 || bottomRightRow >= rows || bottomRightCol < 0 || bottomRightCol >= cols {
		return true // Out of bounds
	}

	return gameboard[topLeftRow][topLeftCol] == Wall ||
		gameboard[topRightRow][topRightCol] == Wall ||
		gameboard[bottomLeftRow][bottomLeftCol] == Wall ||
		gameboard[bottomRightRow][bottomRightCol] == Wall
}

func clearPlayerSpace(board [][]string, x, y int) {
	for dx := 0; dx <= 2; dx++ {
		for dy := 0; dy <= 2; dy++ {
			newX, newY := x+dx, y+dy
			if (dx == 2 && dy == 2) || (dx == 2 && dy == 1) || (dx == 1 && dy == 2) {
				continue // Skip diagonal and far diagonal spaces
			}
			if newX < len(board) && newY < len(board[0]) && board[newX][newY] != Wall {
				board[newX][newY] = Empty
			}
		}
	}
}
