package main

import (
	"bomberman-dom/server"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", server.WebsocketHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%v", 5000), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
