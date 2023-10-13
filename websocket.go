package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Movement struct {
	Type     string `json:"type"`
	Movement struct {
		Direction string `json:"direction"`
	} `json:"movement"`
}

type Player struct {
	ID        int    `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	Speed     int    `json:"speed,omitempty"`
}

type PositionUpdate struct {
	Type string   `json:"type"`
	Data []Player `json:"players"`
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 256,
	ReadBufferSize:  256,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	gameboard := createMap()
	jsonData, err := json.Marshal(gameboard)
	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println(err)
		return
	}

	// playerUpdates := make(chan Player)
	players := []Player{}
	players = append(players, Player{
		ID:    1,
		X:     0, // Initialize to some default or calculated value
		Y:     0, // Initialize to some default or calculated value
		Speed: 1,
	})
	var lastDirection string
	var lastMoveTime time.Time
	// speed := players[0].Speed
	moveCooldown := 200 * time.Millisecond // 200 milliseconds cooldown
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		stepSize := 64

		for range ticker.C {
			if time.Since(lastMoveTime) < moveCooldown {
				continue
			}
			for i, player := range players {
				if player.ID == 1 {
					fmt.Println("here")
					newX, newY := player.X, player.Y
					switch lastDirection {
					case "up":
						newY -= stepSize
					case "down":
						newY += stepSize
					case "left":
						newX -= stepSize
					case "right":
						newX += stepSize
					}
					fmt.Println(newX)
					if !isCollision(gameboard, newX, newY) {

						players[i].X = newX
						players[i].Y = newY
						players[i].Direction = lastDirection
						lastMoveTime = time.Now()
					}
					fmt.Println(players[i].X)
				}
			}

			positionUpdate := PositionUpdate{
				Type: "player_position_update",
				Data: players,
			}
			updateJSON, _ := json.Marshal(positionUpdate)
			conn.WriteMessage(websocket.TextMessage, updateJSON)

		}
	}()
	defer conn.Close()

	for {
		var move Movement
		err := conn.ReadJSON(&move)
		if err != nil {
			return
		}
		fmt.Println(move)
		if move.Type == "keydown" {
			lastDirection = move.Movement.Direction
		}
	}
}
