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
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Message     string `json:"message,omitempty"`
	PlayerCount int    `json:"player_count"`
}

var playerMap = make(map[int]string)

var availableIDs = make(chan int, 4)

func init() {
	for i := 1; i <= 4; i++ {
		availableIDs <- i
	}
}

func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if len(s.Conns) >= 4 {
		conn.WriteJSON(MessageType{Type: "server_full", Message: "Server is currently full, try again later"})
		return
	}

	playerID := <-availableIDs

	s.AddConn(playerID, conn)
	defer func() {
		conn.Close()
		s.ControlChan <- "stop"
		playerCounter--
		s.RemoveConn(playerID)

		availableIDs <- playerID
	}()

	var lastKeydownTime time.Time
	debounceDuration := 50 * time.Millisecond
	for {
		var rawMessage json.RawMessage
		err := conn.ReadJSON(&rawMessage)
		if err != nil {
			log.Println("Error reading JSON:", err)
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
	if s.Game.Players != nil {
		for i := range s.Game.Players {
			if s.Game.Players[i].ID == userID {
				s.Game.Players[i].Lives = 0
				s.playerUpdateChannel <- s.Game.Players
				break
			}
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
}*/

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
		s.connsMu.Lock()
		s.gameMu.Lock()
		defer s.gameMu.Unlock()
		defer s.connsMu.Unlock()
		var registerMsg game.Player
		err = json.Unmarshal(rawMessage, &registerMsg)
		if err != nil {
			log.Println("Error unmarshaling to RegisterMessage:", err)
			return
		}
		newPlayer := game.NewPlayer(registerMsg.Name, playerID)
		s.Game.Players = append(s.Game.Players, *newPlayer)
		s.Conns[playerID].WriteJSON(MessageType{Type: "playerID", Message: fmt.Sprintf("%d", playerID)})
		s.Game.PlayerCount++
		s.playerCountChannel <- *newPlayer

	case "message":
		var msg MessageType
		err = json.Unmarshal(rawMessage, &msg)
		if err != nil {
			log.Println("Error unmarshaling to message:", err)
			return
		}
		for i := range s.Game.Players {
			if s.Game.Players[i].ID == playerID {
				msg.Name = s.Game.Players[i].Name
			}
		}
		s.sendUpdatesToPlayers(msg)
	}

}
