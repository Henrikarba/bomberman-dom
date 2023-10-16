package game

import (
	"math/rand"
	"time"
)

type GameState struct {
	Type string `json:"type,omitempty"`

	Map         [][]string              `json:"map,omitempty"`
	BlockUpdate []BlockUpdate           `json:"block_updates,omitempty"`
	Players     []Player                `json:"players,omitempty"`
	KeysPressed map[int]map[string]bool `json:"-"`
}

type BlockUpdate struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Block string `json:"block"`
}

func AddPoweup() bool {
	randomNumber := rand.Intn(2) + 1
	return randomNumber == 1
}

func PlantBomb(x int, y int, fireDistance int, gameboard [][]string, mapUpdateChannel chan<- []BlockUpdate) {
	blockUpdate := []BlockUpdate{
		{
			X:     x,
			Y:     y,
			Block: "B",
		},
	}

	mapUpdateChannel <- blockUpdate

	go time.AfterFunc(3*time.Second, func() {
		handleExplosion(x, y, fireDistance, gameboard, mapUpdateChannel)
	})
}

func handleExplosion(x int, y int, fireDistance int, gameboard [][]string, mapUpdateChannel chan<- []BlockUpdate) {
	blockUpdate := []BlockUpdate{
		{
			X:     x,
			Y:     y,
			Block: Explosion,
		},
	}

	// Directions: Up, Down, Left, Right
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for _, dir := range directions {
		dx, dy := dir[0], dir[1]
		for i := 1; i <= fireDistance; i++ {
			newX, newY := x+(i*dx), y+(i*dy)

			// Check boundaries
			if newX < 0 || newX >= len(gameboard[0]) || newY < 0 || newY >= len(gameboard) {
				break
			}

			// Check block type
			blockType := gameboard[newY][newX]
			if blockType == Wall || blockType == Flame {
				break
			} else if blockType == Block {
				blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Flame})
				gameboard[newY][newX] = Flame // Update the gameboard immediately
				break
			} else {
				// Empty or already on fire, propagate flame
				blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Flame})
				gameboard[newY][newX] = Flame // Update the gameboard immediately
			}
		}
	}

	mapUpdateChannel <- blockUpdate

	go time.AfterFunc(1*time.Second, func() {
		clearExplosion(blockUpdate, gameboard, mapUpdateChannel)
	})
}

func clearExplosion(blockUpdates []BlockUpdate, gameboard [][]string, mapUpdateChannel chan<- []BlockUpdate) {
	for i := range blockUpdates {
		blockUpdates[i].Block = "e"
	}
	mapUpdateChannel <- blockUpdates
}
