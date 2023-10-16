package main

import (
	"bomberman-dom/server"
	log "bomberman-dom/server/logger"
	"context"
	"fmt"
	"net/http"
)

const PORT = 5000

func main() {
	logFile, err := log.InitLogger()
	defer logFile.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	s := server.New()

	http.HandleFunc("/ws", s.WebsocketHandler)

	go func() {
		var listening bool
		var ctx context.Context

		for {
			cmd := <-s.ControlChan
			if cmd == "start" && !listening {
				fmt.Println("Started listening for keypresses")
				log.Info("Started listening for keypresses")
				listening = true
				ctx, s.CancelFunc = context.WithCancel(context.Background())
				go s.ListenForKeyPress(ctx)
			} else if cmd == "stop" && listening && len(s.Game.Players) == 0 {
				fmt.Println("Stopped listening for keypresses")
				log.Info("Stopped listening for keypresses")
				listening = false
				if s.CancelFunc != nil {
					s.CancelFunc()
				}
			}
		}
	}()

	log.Info("Bomberman running on http://localhost:", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
