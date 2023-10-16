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
}

func New() *Server {

	return &Server{
		gameMu: sync.Mutex{},
		Game:   game.GameState{},

		connsMu:         sync.RWMutex{},
		Conns:           make(map[int]*websocket.Conn),
		keyEventChannel: make(chan game.Movement, 100),
		ControlChan:     make(chan string),

		gameStateChannel:    make(chan game.GameState, 100),
		mapUpdateChannel:    make(chan []game.BlockUpdate, 100),
		playerUpdateChannel: make(chan []game.Player, 100),
	}
}

func (s *Server) NewGame() {
	gameboard := game.CreateMap()
	players := []game.Player{}
	blockUpdates := []game.BlockUpdate{}

	for id := range s.Conns {
		newPlayer := game.NewPlayer(id, gameboard)
		players = append(players, *newPlayer)
	}
	s.Game.BlockUpdate = blockUpdates
	s.Game.Players = players
	s.Game.Map = gameboard
	s.Game.Type = "new_game"
	s.Game.KeysPressed = make(map[int]map[string]bool)
	fmt.Println("?")
	go s.UpdateGameState()
}

func (s *Server) ListenForKeyPress(ctx context.Context) {
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
			// Bomb plant
			if keys[" "] && (s.Game.Players)[i].AvailableBombs > 0 {
				// (s.Game.Players)[i].AvailableBombs--
				bombX := (s.Game.Players)[i].X
				bombY := (s.Game.Players)[i].Y
				fireDistance := (s.Game.Players)[i].FireDistance
				currentMap := s.Game.Map
				go game.PlantBomb(bombX, bombY, fireDistance, currentMap, s.mapUpdateChannel)
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
						fmt.Println("Hit another player")
					} else if typeof == "Wall" || typeof == "Bomb" {
						fmt.Println("Hit a wall or a bomb")
					} else if typeof == "f" {
						fmt.Println("Hit flame, -1 life")
						(s.Game.Players)[i].Lives--
						s.MovePlayer(i, newX, newY, &shouldUpdate)
					} else if typeof == "ex" {
						fmt.Println("Hit explosion, -1 life")
						(s.Game.Players)[i].Lives--
						s.MovePlayer(i, newX, newY, &shouldUpdate)
					} else if typeof == "p" {
						// fmt.Println("+ 1 bomb")
						// s.BlockUpdate = &[]BlockUpdate{}
						// (*s.Map)[newY][newX] = "e"
						// *s.BlockUpdate = append(*s.BlockUpdate, BlockUpdate{X: newX, Y: newY, Block: "e"})
						// updateChannel <- "map_state_update"
						// (*s.Players)[i].AvailableBombs++
						// MovePlayer(newX, newY)
					}
				} else {
					s.MovePlayer(i, newX, newY, &shouldUpdate)
				}
				s.gameMu.Unlock()
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
		case mapUpdate := <-s.mapUpdateChannel:
			s.gameMu.Lock()
			s.Game.BlockUpdate = mapUpdate
			for _, update := range mapUpdate {
				s.Game.Map[update.Y][update.X] = update.Block
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
