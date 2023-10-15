package game

import (
	"fmt"
	"time"
)

type GameState struct {
	Type string `json:"type,omitempty"`

	Map         *[][]string             `json:"map,omitempty"`
	BlockUpdate *[]BlockUpdate          `json:"block_updates,omitempty"`
	Players     *[]Player               `json:"players,omitempty"`
	KeysPressed map[int]map[string]bool `json:"-"`
}

type BlockUpdate struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Block string `json:"block"`
}

func HandleKeyPress(s *GameState, updateChannel chan<- string) {
	for i, player := range *s.Players {
		if keys, ok := s.KeysPressed[player.ID]; ok {

			// Bomb plant
			if keys["enter"] && player.AvailableBombs > 0 {
				if s.BlockUpdate == nil {
					s.BlockUpdate = &[]BlockUpdate{}
					(*s.Map)[player.Y][player.X] = "B"
					*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: player.X, Y: player.Y, Block: "B"})
					updateChannel <- "map_state_update"
				}
			}

			// Movement
			if time.Since(player.LastMoveTime) >= time.Second*time.Duration(100)/time.Duration(player.Speed) {
				newX, newY := player.X, player.Y
				if keys["w"] {
					newY -= 1
					(*s.Players)[i].Direction = "up"
				} else if keys["s"] {
					newY += 1
					(*s.Players)[i].Direction = "down"
				} else if keys["a"] {
					newX -= 1
					(*s.Players)[i].Direction = "left"
				} else if keys["d"] {
					newX += 1
					(*s.Players)[i].Direction = "right"
				}

				collision, typeof := IsCollision(*s.Map, newX, newY, *s.Players, player.ID)
				if collision {
					if typeof == "Flame" {
						fmt.Println("Hit flame")
					} else if typeof == "Player" {
						fmt.Println("Hit another player")
					} else if typeof == "Wall" || typeof == "Bomb" {
						fmt.Println("Hit a wall or a bomb")
					}
				} else {
					(*s.Players)[i].X = newX
					(*s.Players)[i].Y = newY
					(*s.Players)[i].LastMoveTime = time.Now()
					updateChannel <- "player_state_update"
				}
			}
		}
	}

}
