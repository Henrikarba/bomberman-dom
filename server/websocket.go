package server

import (
	"bomberman-dom/game"
	"encoding/json"
	"fmt"
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

type MessageType struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	playerID := getNextPlayerID()
	if playerID > 4 {
		conn.WriteJSON(MessageType{Type: "server_full", Message: "Server is currently full, try again later"})
		playerCounter--
		return
	}

	// Add, remove players
	s.AddConn(playerID, conn)

	var lastKeydownTime time.Time
	debounceDuration := 50 * time.Millisecond
	for {
		var rawMessage json.RawMessage
		err := conn.ReadJSON(&rawMessage)
		if err != nil {
			log.Println("Error reading JSON:", err)
			s.ControlChan <- "stop"
			s.RemoveConn(playerID)
			playerCounter--
			return
		}
		s.handleMessage(rawMessage, playerID, &lastKeydownTime, debounceDuration)
	}
}

func (s *Server) AddConn(userID int, conn *websocket.Conn) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()
	s.Conns[userID] = conn
	fmt.Printf("Added connection with id %d\n", userID)
	fmt.Println("s.Conns = ", s.Conns)
}

func (s *Server) RemoveConn(userID int) {
	// Lock both mutexes
	s.gameMu.Lock()
	s.connsMu.Lock()
	defer s.gameMu.Unlock()
	defer s.connsMu.Unlock()

	// Remove player from game state
	for i := range s.Game.Players {
		if s.Game.Players[i].ID == userID {
			delete(s.Game.KeysPressed, userID)
			s.Game.Players[i].X = -300
			s.Game.Players[i].Y = -300
			s.playerUpdateChannel <- s.Game.Players
			s.Game.Players = s.removePlayerByID(s.Game.Players, userID)
			break
		}
	}

	// Remove connection
	delete(s.Conns, userID)

	fmt.Printf("Removed connection with id %d\n", userID)
	fmt.Println("s.Conns = ", s.Conns)
}

func (s *Server) removePlayerByID(players []game.Player, playerID int) []game.Player {

	for i, player := range players {
		if player.ID == playerID {
			s.Game.PlayerCount--
			return append(players[:i], players[i+1:]...)
		}
	}
	return players
}

func (s *Server) handleMessage(rawMessage json.RawMessage, playerID int, lastKeydownTime *time.Time, debounceDuration time.Duration) {
	var genMsg MessageType
	err := json.Unmarshal(rawMessage, &genMsg)
	if err != nil {
		log.Println("Error unmarshaling to MessageType:", err)
		return
	}

	switch genMsg.Type {
	case "keydown", "keyup":
		var move game.Movement
		err = json.Unmarshal(rawMessage, &move)
		if err != nil {
			log.Println("Error unmarshaling to Movement:", err)
			return
		}
		move.PlayerID = playerID
		if move.Type == "keydown" {
			currentTime := time.Now()
			if currentTime.Sub(*lastKeydownTime) >= debounceDuration {
				lastKeydownTime = &currentTime
				s.keyEventChannel <- move
			}
		} else {
			s.keyEventChannel <- move
		}

	case "register":

		var registerMsg game.Player
		err = json.Unmarshal(rawMessage, &registerMsg)
		if err != nil {
			log.Println("Error unmarshaling to RegisterMessage:", err)
			return
		}
		s.Game.PlayerCount++
		s.playerCountChannel <- registerMsg
	}
}
