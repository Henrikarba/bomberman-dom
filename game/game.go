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

func PlantBomb(x int, y int, fireDistance int, s *GameState, mapUpdateChannel chan<- []BlockUpdate) {
	if s.BlockUpdate == nil {
		s.BlockUpdate = []BlockUpdate{}
	}

	s.BlockUpdate = append(s.BlockUpdate, BlockUpdate{X: x, Y: y, Block: "B"})
	mapUpdateChannel <- s.BlockUpdate

	time.AfterFunc(3*time.Second, func() {
		handleExplosion(x, y, fireDistance, *s, mapUpdateChannel)
	})
}

func handleExplosion(x int, y int, fireDistance int, s GameState, mapUpdateChannel chan<- []BlockUpdate) {
	if s.BlockUpdate == nil {
		s.BlockUpdate = []BlockUpdate{}
	}

	s.BlockUpdate = append(s.BlockUpdate, BlockUpdate{X: x, Y: y, Block: "ex"})
	mapUpdateChannel <- s.BlockUpdate
	FlameBlocks(&s, x, y, 0, 1, fireDistance, mapUpdateChannel)  // Up
	FlameBlocks(&s, x, y, 0, -1, fireDistance, mapUpdateChannel) // Down
	FlameBlocks(&s, x, y, 1, 0, fireDistance, mapUpdateChannel)  // Right
	FlameBlocks(&s, x, y, -1, 0, fireDistance, mapUpdateChannel) // Left
}

func FlameBlocks(s *GameState, x, y, directionX, directionY, fireDistance int, mapUpdateChannel chan<- []BlockUpdate) {
	for i := 1; i <= fireDistance; i++ {
		newX, newY := x+(i*directionX), y+(i*directionY)
		if newX >= 0 && newX < len(s.Map) && newY >= 0 && newY < len((s.Map)[newY]) {
			blockType := (s.Map)[newY][newX]
			if blockType == "d" || blockType == "e" || blockType == "f" {
				(s.Map)[newY][newX] = "f"
				s.BlockUpdate = append(s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "f"})
				mapUpdateChannel <- s.BlockUpdate

				go time.AfterFunc(1000*time.Millisecond, func() {
					clearFlame(s, newX, newY, mapUpdateChannel)
				})
			} else {
				break
			}
		}
	}
}

func clearFlame(s *GameState, x, y int, mapUpdateChannel chan<- []BlockUpdate) {
	if s.BlockUpdate == nil {
		s.BlockUpdate = []BlockUpdate{}
	}

	s.BlockUpdate = append(s.BlockUpdate, BlockUpdate{X: x, Y: y, Block: "e"})
	mapUpdateChannel <- s.BlockUpdate
}
