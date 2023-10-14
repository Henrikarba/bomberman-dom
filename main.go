package main

import (
	"bomberman-dom/server"
	"fmt"
	"log"
	"net/http"
)

const PORT = 5000

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", server.WebsocketHandler)

	log.Printf("Bomberman running on http://localhost:%d", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
