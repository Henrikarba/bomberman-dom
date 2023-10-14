package server

import (
	"bomberman-dom/game"
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

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
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

	gameboard := game.CreateMap()
	newPlayer := game.NewPlayer(1, gameboard)
	newPlayer2 := game.NewPlayer(2, gameboard)
	newPlayer3 := game.NewPlayer(3, gameboard)
	newPlayer4 := game.NewPlayer(4, gameboard)
	players := []game.Player{}
	players = append(players, *newPlayer, *newPlayer2, *newPlayer3, *newPlayer4)

	keysPressed := make(map[int]map[string]bool)
	keysPressed[playerID] = make(map[string]bool)
	keyEventChannel := make(chan game.Movement, 100) // Buffered channel
	var blockUpdates []game.BlockUpdate

	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		for range ticker.C {
			for len(keyEventChannel) > 0 {
				move := <-keyEventChannel
				for k := range keysPressed[playerID] {
					keysPressed[playerID][k] = false

				}
				for _, key := range move.Keys {
					keysPressed[playerID][key] = true
					game.HandleKeyPress(players, &gameboard, keysPressed, &blockUpdates)
				}
			}

			gameState := game.GameState{
				Players: players,
			}

			if len(blockUpdates) > 0 {
				gameState.Type = "map_state_update"
				gameState.BlockUpdate = blockUpdates
			} else {
				gameState.Type = "player_position_update"
			}

			// Send the game state
			if err := conn.WriteJSON(gameState); err != nil {
				log.Println("Write JSON error:", err)
			}

			// Clear the block updates
			blockUpdates = []game.BlockUpdate{}
		}
	}()

	// Send initial game state
	lobby := game.GameState{
		Type:    "new_game",
		Map:     gameboard,
		Players: players,
	}
	conn.WriteJSON(lobby)

	// Main loop for listening to key events
	var lastKeydownTime time.Time
	debounceDuration := 50 * time.Millisecond
	for {
		var move = game.Movement{}
		err := conn.ReadJSON(&move)
		if err != nil {
			close(keyEventChannel)
			log.Println("Error reading JSON:", err)
			return
		}
		fmt.Println(move.Keys)
		if move.Type == "keydown" {
			currentTime := time.Now()
			if currentTime.Sub(lastKeydownTime) >= debounceDuration {
				lastKeydownTime = currentTime
				keyEventChannel <- move
			}
		} else {
			keyEventChannel <- move
		}
	}
}
