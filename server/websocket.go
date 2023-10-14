package server

import (
	"bomberman-dom/game"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 256,
	ReadBufferSize:  256,
}

var playerCounter = 1
var playerCounterMutex = &sync.Mutex{}

func getNextPlayerID() int {
	playerCounterMutex.Lock()
	defer playerCounterMutex.Unlock()
	id := playerCounter
	playerCounter++
	return id
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	playerID := getNextPlayerID()
	if playerID > 4 {
		log.Println("Maximum players reached")
		return
	}

	s.AddConn(playerID, conn)
	defer s.RemoveConn(playerID)
	if playerID == 1 {
		s.NewGame()
	} else {
		newPlayer := game.NewPlayer(playerID, *s.Game.Map)
		updatedPlayers := append(*s.Game.Players, *newPlayer)
		s.Game.Players = &updatedPlayers
	}

	s.Game.KeysPressed[playerID] = make(map[string]bool)
	s.ControlChan <- "start"

	s.Game.Type = "new_game"
	conn.WriteJSON(s.Game)
	var lastKeydownTime time.Time
	debounceDuration := 50 * time.Millisecond
	for {
		var move = game.Movement{}
		err := conn.ReadJSON(&move)
		if err != nil {
			log.Println("Error reading JSON:", err)
			playerCounter--
			s.ControlChan <- "stop"
			return
		}
		move.PlayerID = playerID

		if move.Type == "keydown" {
			currentTime := time.Now()
			if currentTime.Sub(lastKeydownTime) >= debounceDuration {
				lastKeydownTime = currentTime
				s.keyEventChannel <- move
			}
		} else {
			s.keyEventChannel <- move
		}
	}
}

func (s *Server) AddConn(userID int, conn *websocket.Conn) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	s.Conns[userID] = conn
}

func (s *Server) RemoveConn(userID int) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	delete(s.Conns, userID)
}
