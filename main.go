package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", websocketHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%v", 5000), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
