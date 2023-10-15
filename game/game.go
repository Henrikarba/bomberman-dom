package game

import (
	"fmt"
	"math/rand"
	"time"
	log "bomberman-dom/server/logger"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func AddPoweup() bool {
	randomNumber := rand.Intn(2) + 1
	return randomNumber == 1
}

func HandleKeyPress(s *GameState, updateChannel chan<- string) {
	for i, player := range *s.Players {
		if keys, ok := s.KeysPressed[player.ID]; ok {

			// Bomb plant
			if keys[" "] && player.AvailableBombs > 0 {
				if s.BlockUpdate == nil {
					s.BlockUpdate = &[]BlockUpdate{}
					(*s.Map)[player.Y][player.X] = "B"
					*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: player.X, Y: player.Y, Block: "B"})
					updateChannel <- "map_state_update"

					(*s.Players)[i].AvailableBombs--
					log.Info(player.X, player.Y, "start", player.ID)
					// Explosion
					bombX, bombY := player.X, player.Y
					FlameBlocks := func(x, y int, directionX, directionY int) {
						for i := 1; i <= player.FireDistance; i++ {
							newX, newY := x+(i*directionX), y+(i*directionY)
							if newX >= 0 && newX < len(*s.Map) && newY >= 0 && newY < len((*s.Map)[newY]) {
								if (*s.Map)[newY][newX] == "d" {
									// If it's a "d" block, update it to "e" and exit the loop
									if !AddPoweup() {
										(*s.Map)[newY][newX] = "e"
										*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "e"})
									} else {
										(*s.Map)[newY][newX] = "p"
										*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "p"})
									}
									break
								} else if (*s.Map)[newY][newX] != "ex" && (*s.Map)[newY][newX] == "e" {
									// If it's not "ex" and is "e", update it to "f"
									(*s.Map)[newY][newX] = "f"
									*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "f"})
									flameX, flameY := newX, newY
									time.AfterFunc(1000*time.Millisecond, func() {
										s.BlockUpdate = &[]BlockUpdate{}
										(*s.Map)[flameY][flameX] = "e"
										*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: flameX, Y: flameY, Block: "e"})
										updateChannel <- "map_state_update"
									})
								} else if (*s.Map)[newY][newX] == "f" {
									// If it's "f", update it to "f" with new timer
									(*s.Map)[newY][newX] = "f"
									*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "f"})
									flamereX, flamereY := newX, newY
									time.AfterFunc(1000*time.Millisecond, func() {
										s.BlockUpdate = &[]BlockUpdate{}
										(*s.Map)[flamereY][flamereX] = "e"
										*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: flamereX, Y: flamereY, Block: "e"})
										updateChannel <- "map_state_update"
									})
								} else {
									// If out of bounds, exit the loop
									break
								}
							}
						}
					}

					time.AfterFunc(3*time.Second, func() {
						if s.BlockUpdate == nil {
							s.BlockUpdate = &[]BlockUpdate{}
							(*s.Map)[player.Y][player.X] = "ex"
							*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: bombX, Y: bombY, Block: "ex"})
							updateChannel <- "map_state_update"

							(*s.Players)[i].AvailableBombs++

							FlameBlocks(bombX, bombY, 0, 1)  // Up
							FlameBlocks(bombX, bombY, 0, -1) // Down
							FlameBlocks(bombX, bombY, 1, 0)  // Right
							FlameBlocks(bombX, bombY, -1, 0) // Left
							time.AfterFunc(1100*time.Millisecond, func() {
								s.BlockUpdate = &[]BlockUpdate{}
								(*s.Map)[player.Y][player.X] = "e"
								*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: bombX, Y: bombY, Block: "e"})
								updateChannel <- "map_state_update"
							})
						}
						log.Info(player.X, player.Y, "end", player.ID)
					})
				}
			}

			// Movement
			if time.Since(player.LastMoveTime) >= time.Second*time.Duration(100)/time.Duration(player.Speed) {
				MovePlayer := func(newX, newY int) {
					(*s.Players)[i].X = newX
					(*s.Players)[i].Y = newY
					(*s.Players)[i].LastMoveTime = time.Now()
					updateChannel <- "player_state_update"
				}
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
					if typeof == "Player" {
						fmt.Println("Hit another player")
					} else if typeof == "Wall" || typeof == "Bomb" {
						fmt.Println("Hit a wall or a bomb")
					} else if typeof == "f" {
						fmt.Println("Hit flame, -1 life")
						(*s.Players)[i].Lives--
						MovePlayer(newX, newY)
					} else if typeof == "ex" {
						fmt.Println("Hit explosion, -1 life")
						(*s.Players)[i].Lives--
						MovePlayer(newX, newY)
					} else if typeof == "p" {
						fmt.Println("+ 1 bomb")
						s.BlockUpdate = &[]BlockUpdate{}
						(*s.Map)[newY][newX] = "e"
						*s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "e"})
						updateChannel <- "map_state_update"
						(*s.Players)[i].AvailableBombs++
						MovePlayer(newX, newY)
					}
				} else {
					MovePlayer(newX, newY)
				}
			}
		}
	}

}
