package game

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	Empty     = "e"
	Power1    = "p1"
	Power2    = "p2"
	Power3    = "p3"
	Wall      = "0"
	Block     = "d"
	Flame     = "f"
	Explosion = "ex"

	Bomb = "B"
)

func CreateMap() [][]string {
	rows, cols := 13, 19
	centerRow, centerCol := rows/2, cols/2
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
			if (i == 0 && j == 0) || (i == 0 && j == cols-1) || (i == rows-1 && j == 0) || (i == rows-1 && j == cols-1) {
				continue
			}
			if i%3 == 0 && j%3 == 0 {
				board[i][j] = Wall
			}
		}
	}
	// Clear player spaces
	clearPlayerSpace(board, 0, 0)
	clearPlayerSpace(board, 0, cols-1)
	clearPlayerSpace(board, rows-1, 0)
	clearPlayerSpace(board, rows-1, cols-1)

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

				distance := math.Abs(float64(i-centerRow)) + math.Abs(float64(j-centerCol))

				probability := 0.6 + 0.01*distance

				if rand.Float32() < float32(probability) {
					board[i][j] = Block
				}
			}
		}
	}

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

func IsCollision(gameboard [][]string, x, y int, players []Player, currentPlayerID int) (bool, string) {
	rows := len(gameboard)
	cols := len(gameboard[0])

	if x < 0 || x >= cols || y < 0 || y >= rows {
		return true, "Wall"
	}
	for _, player := range players {
		if player.ID != currentPlayerID && player.X == x && player.Y == y && player.Lives > 0 {
			return true, "Player"
		}
	}
	return gameboard[y][x] == Wall || gameboard[y][x] == Block || gameboard[y][x] == Bomb || gameboard[y][x] == Flame || gameboard[y][x] == Explosion || gameboard[y][x] == Power1 || gameboard[y][x] == Power2 || gameboard[y][x] == Power3, gameboard[y][x]
	// return false, ""
}

func clearPlayerSpace(board [][]string, x, y int) {
	for dx := 0; dx <= 2; dx++ {
		for dy := 0; dy <= 2; dy++ {
			newX, newY := x+dx, y+dy
			if (dx == 2 && dy == 2) || (dx == 2 && dy == 1) || (dx == 1 && dy == 2) {
				continue
			}
			if newX < len(board) && newY < len(board[0]) && board[newX][newY] != Wall {
				board[newX][newY] = Empty
			}
		}
	}
}
