package server

import (
	"bomberman-dom/game"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	Game *game.GameState

	connsMu sync.RWMutex
	Conns   map[int]*websocket.Conn

	keyEventChannel chan game.Movement
	ControlChan     chan string
	CancelFunc      context.CancelFunc

	updateChannel chan string
}

func New() *Server {
	return &Server{
		Game:            &game.GameState{},
		Conns:           make(map[int]*websocket.Conn),
		keyEventChannel: make(chan game.Movement, 100),
		ControlChan:     make(chan string),
		updateChannel:   make(chan string, 100),
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

	s.Game.BlockUpdate = &blockUpdates
	s.Game.Players = &players
	s.Game.Map = &gameboard
	s.Game.Type = "new_game"
	s.Game.KeysPressed = make(map[int]map[string]bool)
	s.updateChannel <- "new_game"
}

func (s *Server) ListenForKeyPress(ctx context.Context) {
	// Initialize s.Game.KeysPressed for all players
	for id := range s.Conns {
		s.Game.KeysPressed[id] = make(map[string]bool)
	}

	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		defer ticker.Stop()
		data := game.GameState{}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Process all pending key events
				for len(s.keyEventChannel) > 0 {
					move := <-s.keyEventChannel
					playerID := move.PlayerID

					for k := range s.Game.KeysPressed[playerID] {
						s.Game.KeysPressed[playerID][k] = false
					}

					for _, key := range move.Keys {
						s.Game.KeysPressed[playerID][key] = true
					}

					game.HandleKeyPress(s.Game, s.updateChannel)
				}
				select {
				case update := <-s.updateChannel:
					data.Type = update
					if update == "map_state_update" {
						data.BlockUpdate = s.Game.BlockUpdate
					} else if update == "new_game" {
						data = *s.Game
					} else if update == "player_state_update" {
						data.Players = s.Game.Players
					}
				default:
					data.Type = ""
					data.BlockUpdate = nil
					data.Players = nil
					data.Map = nil
				}
				// Send updated game state to all players
				s.connsMu.Lock()
				if data.Type != "" { // Only send if there's an update
					for _, conn := range s.Conns {
						if err := conn.WriteJSON(data); err != nil {
							log.Println("Write JSON error:", err)
						}
					}
				}
				s.connsMu.Unlock()
				s.Game.BlockUpdate = nil
			}
		}
	}()
}
