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
	playerLeaveChannel  chan game.Player
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
		playerLeaveChannel:  make(chan game.Player, 4),
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
	for _, player := range s.Game.Players {
		s.Game.KeysPressed[player.ID] = make(map[string]bool)
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
					bombX := s.Game.Players[i].X
					bombY := s.Game.Players[i].Y

					if s.Game.Map[bombY][bombX] != game.Bomb {
						s.gameMu.Lock()
						s.Game.Players[i].AvailableBombs--
						// fmt.Println(fmt.Sprintf("PlayerID: %d has %d bombs after planting bomb", i, s.Game.Players[i].AvailableBombs))
						s.gameMu.Unlock()

						fireDistance := s.Game.Players[i].FireDistance
						currentMap := s.Game.Map
						go game.PlantBomb(bombX, bombY, fireDistance, currentMap, s.mapUpdateChannel)
						go game.ReplenishBomb(s.Game.Players[i].ID, &s.Game, &s.gameMu)
					}
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
						} else if typeof == "f" && !s.Game.Players[i].Damaged {
							fmt.Println("Hit flame, -1 life")
							s.Game.Players[i].Lives--
							s.HandleDamage(i)
							s.MovePlayer(i, newX, newY, &shouldUpdate)
							if s.Game.Players[i].Lives <= 0 {
								s.lostGame(s.Game.Players[i])
							}
						} else if typeof == "ex" && !s.Game.Players[i].Damaged {
							fmt.Println("Hit explosion, -1 life")
							s.Game.Players[i].Lives--
							s.HandleDamage(i)
							s.MovePlayer(i, newX, newY, &shouldUpdate)
							if s.Game.Players[i].Lives <= 0 {
								s.lostGame(s.Game.Players[i])
							}
						} else if typeof == "p1" {
							fmt.Println("+ 1 bomb")
							go game.ClearPowerup(newX, newY, s.Game.Map, s.mapUpdateChannel)
							// fmt.Println(fmt.Sprintf("PlayerID: %d has %d bombs before powerup", i, s.Game.Players[i].AvailableBombs))
							s.Game.Players[i].AvailableBombs++
							// fmt.Println(fmt.Sprintf("PlayerID: %d has %d bombs after powerup", i, s.Game.Players[i].AvailableBombs))
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

func (s *Server) HandleDamage(i int) {
	s.Game.Players[i].Damaged = true
	go func() {
		time.Sleep(2 * time.Second)
		s.Game.Players[i].Damaged = false
	}()
}

func (s *Server) UpdateGameState() {
	exitChan := make(chan struct{})
	data := game.GameState{}
	for {
		data.BlockUpdate = nil
		data.Players = nil
		select {
		case <-exitChan:
			fmt.Println("gg")
			return
		case gameStateUpdate := <-s.gameStateChannel:
			data = gameStateUpdate
		case mapUpdate := <-s.mapUpdateChannel:
			s.gameMu.Lock()
			s.Game.BlockUpdate = mapUpdate
			for _, update := range mapUpdate {
				// Check for player hits
				if update.Block == "f" || update.Block == "ex" {
					for i, player := range s.Game.Players {
						if s.Game.Players[i].Lives > 0 {
							if player.X == update.X && player.Y == update.Y && !s.Game.Players[i].Damaged {
								s.Game.Players[i].Lives--
								s.HandleDamage(i)
								s.playerUpdateChannel <- s.Game.Players
								if s.Game.Players[i].Lives <= 0 {
									s.lostGame(s.Game.Players[i])
									s.connsMu.Lock()
									s.Conns[s.Game.Players[i].ID].WriteJSON(MessageType{Type: "game_over"})
									s.connsMu.Unlock()
								}
							}
						}
					}
				}
				s.Game.Map[update.Y][update.X] = update.Block
			}
			s.gameMu.Unlock()

			data.Type = "map_state_update"
			data.BlockUpdate = mapUpdate
		case playerUpdate := <-s.playerUpdateChannel:
			s.gameMu.Lock()
			data.Type = "player_state_update"
			data.Players = playerUpdate
			s.Game.Players = playerUpdate
			if s.Game.Alive == 1 {
				for i := range s.Game.Players {
					if s.Game.Players[i].Lives > 0 {
						s.sendUpdatesToPlayers(MessageType{Type: "gg", Name: "Server", Message: fmt.Sprintf("%s won!", s.Game.Players[i].Name)})
						time.AfterFunc(3*time.Second, func() {
							exitChan <- struct{}{}
						})
					}
				}
			}
			s.gameMu.Unlock()
		}

		s.sendUpdatesToPlayers(data)

	}
}

func (s *Server) MonitorPlayerCount() {
	data := game.GameState{
		Type: "status",
	}

	var cancelFunc context.CancelFunc
	var ctx context.Context
	playerCountChange := make(chan int, 4)

	for {
		select {
		case player := <-s.playerCountChannel:
			s.sendUpdatesToPlayers(MessageType{Type: "message", Name: "Server", Message: fmt.Sprintf("%s joined.", player.Name), Player: s.Game.Players})
			if len(s.Game.Players) == 2 && cancelFunc == nil {
				ctx, cancelFunc = context.WithCancel(context.Background())
				go s.startCountdown(ctx, cancelFunc, data, playerCountChange, 30)
			}
			if len(s.Game.Players) == 4 {
				playerCountChange <- 4
			}

		case player := <-s.playerLeaveChannel:
			s.sendUpdatesToPlayers(MessageType{Type: "message", Name: "Server", Message: fmt.Sprintf("%s left.", player.Name), Player: s.Game.Players})
			if len(s.Game.Players) < 2 && cancelFunc != nil {
				cancelFunc()
				cancelFunc = nil
			}
		}
	}
}

func (s *Server) startCountdown(ctx context.Context, cancelFunc context.CancelFunc, data game.GameState, playerCountChange chan int, startValue int) {
	for i := startValue; i >= 0; i-- {
		select {
		case <-ctx.Done():
			fmt.Println("Countdown cancelled")
			return
		case newCount := <-playerCountChange:
			if newCount == 4 {
				i = 10
			}

		default:
			data.CountDown = i
			s.sendUpdatesToPlayers(data)
			time.Sleep(1 * time.Second)
			fmt.Println("Starting game in: ", i)
			if i == 0 {
				s.NewGame()
				s.Game.Alive = len(s.Game.Players)

				for i, player := range s.Game.Players {
					s.gameMu.Lock()
					playa := game.StartingPositions(&player)
					s.Game.Players[i] = *playa
					s.Game.KeysPressed[s.Game.Players[i].ID] = make(map[string]bool)
					s.gameMu.Unlock()
				}
				fmt.Println(s.Game.Players)
				s.Game.Playing = true
				s.ControlChan <- "start"
				s.gameStateChannel <- s.Game
				s.playerUpdateChannel <- s.Game.Players
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
	s.Game.Alive--
	s.playerUpdateChannel <- s.Game.Players
}
