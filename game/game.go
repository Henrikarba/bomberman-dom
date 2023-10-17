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

	PlayerCount int `json:"player_count,omitempty"`
	CountDown   int `json:"countdown,omitempty"`
}

type BlockUpdate struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Block string `json:"block"`
}

func AddPoweup() int {
	randomNumber := rand.Intn(3) + 1
	return randomNumber
}

func ClearPowerup(x int, y int, gameboard [][]string, mapUpdateChannel chan<- []BlockUpdate) {
	blockUpdate := []BlockUpdate{
		{
			X:     x,
			Y:     y,
			Block: "e",
		},
	}

	mapUpdateChannel <- blockUpdate
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
			if blockType == Wall {
				break
			} else if blockType == Block {
				test := AddPoweup()
				if test == 1 {
					blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Power1})
					gameboard[newY][newX] = Power1
				} else if test == 2 {
					blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Power2})
					gameboard[newY][newX] = Power2
				} else if test == 3 {
					blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Power3})
					gameboard[newY][newX] = Power3
				} else {
					blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Flame})
				}
				break
			} else {
				// Empty or already on fire, propagate flame
				blockUpdate = append(blockUpdate, BlockUpdate{X: newX, Y: newY, Block: Flame})
				gameboard[newY][newX] = Flame
			}
		}
	}

	mapUpdateChannel <- blockUpdate

	go time.AfterFunc(1*time.Second, func() {
		clearExplosion(blockUpdate, gameboard, mapUpdateChannel)
	})
}

func clearExplosion(blockUpdates []BlockUpdate, gameboard [][]string, mapUpdateChannel chan<- []BlockUpdate) {
	clearUpdates := make([]BlockUpdate, len(blockUpdates))
	for i, update := range blockUpdates {
		if update == (BlockUpdate{X: update.X, Y: update.Y, Block: "f"}) || update == (BlockUpdate{X: update.X, Y: update.Y, Block: "ex"}) {
			clearUpdates[i] = BlockUpdate{X: update.X, Y: update.Y, Block: "e"}
		}
	}
	mapUpdateChannel <- clearUpdates
}
