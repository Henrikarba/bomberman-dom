package server

import (
	"bomberman-dom/game"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	gameMu sync.Mutex
	Game   game.GameState

	connsMu sync.RWMutex
	Conns   map[int]*websocket.Conn

	keyEventChannel chan game.Movement
	ControlChan     chan string
	CancelFunc      context.CancelFunc

	gameStateChannel    chan game.GameState
	mapUpdateChannel    chan []game.BlockUpdate
	playerUpdateChannel chan []game.Player
	playerCountChannel  chan game.Player
}

func New() *Server {

	return &Server{
		gameMu: sync.Mutex{},
		Game:   game.GameState{},

		connsMu:         sync.RWMutex{},
		Conns:           make(map[int]*websocket.Conn),
		keyEventChannel: make(chan game.Movement, 100),
		ControlChan:     make(chan string, 4),

		gameStateChannel:    make(chan game.GameState, 100),
		mapUpdateChannel:    make(chan []game.BlockUpdate, 100),
		playerUpdateChannel: make(chan []game.Player, 100),
		playerCountChannel:  make(chan game.Player, 4),
	}
}

func (s *Server) NewGame() {
	fmt.Println("Initializing New Game")
	s.Game.KeysPressed = make(map[int]map[string]bool)

	gameboard := game.CreateMap()
	blockUpdates := []game.BlockUpdate{}

	s.Game.BlockUpdate = blockUpdates
	s.Game.Map = gameboard
	s.Game.Type = "new_game"
	go s.UpdateGameState()
}

func (s *Server) ListenForKeyPress(ctx context.Context) {
	if s.Game.KeysPressed == nil {
		fmt.Println("s.Game.KeysPressed is nil")
		return
	}
	fmt.Println("Inside ListenForKeyPress, s.Game.KeysPressed:", s.Game.KeysPressed)
	// Initialize s.Game.KeysPressed for all players
	for id := range s.Conns {
		s.Game.KeysPressed[id] = make(map[string]bool)
	}

	for {
		select {
		case <-ctx.Done():
			// Exit the loop if the context is cancelled
			return
		case move := <-s.keyEventChannel:
			// Process the key event
			playerID := move.PlayerID
			for k := range s.Game.KeysPressed[playerID] {
				s.Game.KeysPressed[playerID][k] = false
			}
			for _, key := range move.Keys {
				s.Game.KeysPressed[playerID][key] = true
			}
			s.HandleKeyPress()
		}
	}
}

func (s *Server) HandleKeyPress() {
	shouldUpdate := false // Flag to check if update is needed
	for i, player := range s.Game.Players {
		if keys, ok := s.Game.KeysPressed[player.ID]; ok {
			if player.Lives > 0 {

				// Bomb plant
				if keys[" "] && (s.Game.Players)[i].AvailableBombs > 0 {
					s.gameMu.Lock()
					s.Game.Players[i].AvailableBombs--
					s.gameMu.Unlock()
					bombX := s.Game.Players[i].X
					bombY := s.Game.Players[i].Y
					fireDistance := s.Game.Players[i].FireDistance
					currentMap := s.Game.Map
					go game.PlantBomb(bombX, bombY, fireDistance, currentMap, s.mapUpdateChannel)
					go func(playerID int) {
						time.AfterFunc(3500*time.Millisecond, func() {
							s.gameMu.Lock()
							defer s.gameMu.Unlock()

							for i := range s.Game.Players {
								if s.Game.Players[i].ID == playerID {
									s.Game.Players[i].AvailableBombs++
									break
								}
							}
						})
					}(s.Game.Players[i].ID)
				}

				// Movement
				if time.Since(player.LastMoveTime) >= time.Second*time.Duration(100)/time.Duration(player.Speed) {
					newX, newY := player.X, player.Y
					if keys["w"] {
						newY -= 1
						s.Game.Players[i].Direction = "up"
					} else if keys["s"] {
						newY += 1
						s.Game.Players[i].Direction = "down"
					} else if keys["a"] {
						newX -= 1
						s.Game.Players[i].Direction = "left"
					} else if keys["d"] {
						newX += 1
						s.Game.Players[i].Direction = "right"
					}

					s.gameMu.Lock()
					collision, typeof := game.IsCollision(s.Game.Map, newX, newY, s.Game.Players, player.ID)
					if collision {
						if typeof == "Player" {
							// fmt.Println("Hit another player")
						} else if typeof == "Wall" || typeof == "Bomb" {
							// fmt.Println("Hit a wall or a bomb")
						} else if typeof == "f" {
							fmt.Println("Hit flame, -1 life")
							s.Game.Players[i].Lives--
							s.Game.Players[i].Damaged = true
							go func() {
								time.Sleep(2 * time.Second)
								s.Game.Players[i].Damaged = false
							}()
							s.MovePlayer(i, newX, newY, &shouldUpdate)
							if s.Game.Players[i].Lives <= 0 {
								s.lostGame(s.Game.Players[i])
							}
						} else if typeof == "ex" {
							fmt.Println("Hit explosion, -1 life")
							s.Game.Players[i].Lives--
							s.Game.Players[i].Damaged = true
							go func() {
								time.Sleep(2 * time.Second)
								s.Game.Players[i].Damaged = false
							}()
							s.MovePlayer(i, newX, newY, &shouldUpdate)
							if s.Game.Players[i].Lives <= 0 {
								s.lostGame(s.Game.Players[i])
							}
						} else if typeof == "p1" {
							fmt.Println("+ 1 bomb")
							go game.ClearPowerup(newX, newY, s.Game.Map, s.mapUpdateChannel)
							s.Game.Players[i].AvailableBombs++
							s.MovePlayer(i, newX, newY, &shouldUpdate)
						} else if typeof == "p2" {
							fmt.Println("+ 100 speed")
							go game.ClearPowerup(newX, newY, s.Game.Map, s.mapUpdateChannel)
							s.Game.Players[i].Speed += 100
							s.MovePlayer(i, newX, newY, &shouldUpdate)
						} else if typeof == "p3" {
							fmt.Println("+ 1 range")
							go game.ClearPowerup(newX, newY, s.Game.Map, s.mapUpdateChannel)
							s.Game.Players[i].FireDistance++
							s.MovePlayer(i, newX, newY, &shouldUpdate)
						}
					} else {
						s.MovePlayer(i, newX, newY, &shouldUpdate)
					}
					s.gameMu.Unlock()
				}
			}
		}
	}
	if shouldUpdate {
		s.playerUpdateChannel <- s.Game.Players
	}
}

func (s *Server) MovePlayer(i, newX, newY int, shouldUpdate *bool) {
	s.Game.Players[i].X = newX
	s.Game.Players[i].Y = newY
	s.Game.Players[i].LastMoveTime = time.Now()
	*shouldUpdate = true
}

func (s *Server) UpdateGameState() {
	data := game.GameState{}
	for {
		data.BlockUpdate = nil
		data.Players = nil
		select {
		case gameStateUpdate := <-s.gameStateChannel:
			data = gameStateUpdate
			break
		case mapUpdate := <-s.mapUpdateChannel:
			s.gameMu.Lock()
			s.Game.BlockUpdate = mapUpdate
			for _, update := range mapUpdate {
				// Check for player hits
				if update.Block == "f" || update.Block == "ex" {
					for i, player := range s.Game.Players {
						if player.X == update.X && player.Y == update.Y {
							s.Game.Players[i].Lives--
							s.Game.Players[i].Damaged = true
							go func() {
								time.Sleep(2 * time.Second)
								s.Game.Players[i].Damaged = false
							}()
							s.playerUpdateChannel <- s.Game.Players
							if s.Game.Players[i].Lives <= 0 {
								s.lostGame(s.Game.Players[i])
							}
						}
					}
				}
				if update.Y != -300 && update.X != -300 {
					s.Game.Map[update.Y][update.X] = update.Block
				}
			}
			s.gameMu.Unlock()

			data.Type = "map_state_update"
			data.BlockUpdate = mapUpdate
			break
		case playerUpdate := <-s.playerUpdateChannel:
			s.gameMu.Lock()
			data.Type = "player_state_update"
			data.Players = playerUpdate
			s.Game.Players = playerUpdate
			s.gameMu.Unlock()
			break
		}

		s.sendUpdatesToPlayers(data)

	}
}

func (s *Server) MonitorPlayerCount() {
	data := game.GameState{
		Type: "status",
	}
	for {
		select {
		case player := <-s.playerCountChannel:
			s.gameMu.Lock()
			currentCount := s.Game.PlayerCount
			s.gameMu.Unlock()

			if currentCount >= 2 {
				s.NewGame()
				for id := range s.Conns {
					player := game.NewPlayer(id, s.Game.Map, player.Name)
					s.Game.Players = append(s.Game.Players, *player)
					s.Game.KeysPressed[id] = make(map[string]bool)
				}

				for i := 6; i >= 0; i-- {
					data.CountDown = i
					s.sendUpdatesToPlayers(data)
					time.Sleep(1 * time.Second)
					fmt.Println("Starting game in: ", i)
					if i == 0 {
						s.ControlChan <- "start"
						s.gameStateChannel <- s.Game
					}
				}
			}
		}
	}
}

func (s *Server) sendUpdatesToPlayers(data interface{}) {

	s.connsMu.Lock()
	for _, conn := range s.Conns {
		if err := conn.WriteJSON(data); err != nil {
			log.Println("Write JSON error:", err)
		}
	}
	s.connsMu.Unlock()
	s.gameMu.Lock()
	s.Game.BlockUpdate = nil
	s.gameMu.Unlock()
}

func (s *Server) lostGame(player game.Player) {
	fmt.Println("game lost", player.ID)
	s.playerUpdateChannel <- s.Game.Players
}
