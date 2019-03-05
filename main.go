package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/websocket"
	"net/http"
	"os"
)

var (
	log = logrus.New()
)


func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	hub, err := websocket.NewHub()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("public/app")))
	http.HandleFunc("/ws", websocket.Handler(hub))
	log.Println(http.ListenAndServe(":"+port, nil))
}