package game

import (
	"fmt"
	"time"
)

type GameState struct {
	Type string `json:"type"`

	Map         [][]string    `json:"map,omitempty"`
	BlockUpdate []BlockUpdate `json:"block_updates,omitempty"`
	Players     []Player      `json:"players,omitempty"`
}

type BlockUpdate struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Block string `json:"block"`
}

func HandleKeyPress(players []Player, gameboard *[][]string, keysPressed map[int]map[string]bool, blockUpdates *[]BlockUpdate) {
	for i, player := range players {
		if keys, ok := keysPressed[player.ID]; ok {
			if keys["enter"] && player.AvailableBombs > 0 {
				(*gameboard)[player.Y][player.X] = "B"
				*blockUpdates = append(*blockUpdates, BlockUpdate{X: player.X, Y: player.Y, Block: "B"})
			}

			if time.Since(player.LastMoveTime) >= time.Second/time.Duration(player.Speed*2) {
				newX, newY := player.X, player.Y
				if keys["w"] {
					newY -= 1
					players[i].Direction = "up"
				} else if keys["s"] {
					newY += 1
					players[i].Direction = "down"
				} else if keys["a"] {
					newX -= 1
					players[i].Direction = "left"
				} else if keys["d"] {
					newX += 1
					players[i].Direction = "right"
				}

				collision, typeof := IsCollision(*gameboard, newX, newY, players, player.ID)
				if collision {
					if typeof == "Flame" {
						fmt.Println("Hit flame")
					} else if typeof == "Player" {
						fmt.Println("Hit another player")
					} else if typeof == "Wall" || typeof == "Bomb" {
						fmt.Println("Hit a wall or a bomb")
					}
				} else {
					players[i].X = newX
					players[i].Y = newY
				}

			}
		}
	}
}
