package main

import (
	"net/http"
	"os"
	"github.com/Sirupsen/logrus"
)

var (
	log = logrus.New()
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is where the websocket connects."))
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleHello)
	log.Println(http.ListenAndServe(":"+port, nil))
}