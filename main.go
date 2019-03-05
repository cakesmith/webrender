package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/system"
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

	hub.OnRegister = func(client *websocket.Client) {
		d := system.DisplayWriter{
			Writer: client,
		}

		for x := 0; x < 256; x++ {
			go d.DrawPixel(x, x, system.ColorBlack)
			go d.DrawPixel(256-x, x, system.ColorBlack)
		}
	}

	http.Handle("/", http.FileServer(http.Dir("public/app")))
	http.HandleFunc("/ws", hub.Handler())
	log.Println(http.ListenAndServe(":"+port, nil))
}