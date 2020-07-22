package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func docserve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	d.AddNode(c)
}

var d = NewDocument()

func main() {
	log.SetFlags(0)
	http.HandleFunc("/", docserve)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
