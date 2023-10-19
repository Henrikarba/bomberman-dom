package server

import (
	"log"
	"net/http"
)

type Post struct {
	//
}

type User struct {
	//
}

type Message struct {
	Type string `json:"type"`
	Post []Post `json:"posts"`
	User []User `json:"users"`
}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading JSON:", err)
			return
		}
		handleMessage(msg)
	}
}

func handleMessage(m Message) {
	switch m.Type {
	case "get_posts":
		// get posts.....
	case "get_users":
		// get users.....
	}
}
