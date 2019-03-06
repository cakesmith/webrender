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
	width = 512
	height = 330
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

		d.DrawRectangle(0, 0, width, height, system.ColorBackground)

		for x := 8; x < width; x = x + 8 {
			go d.DrawVert(x, 0, height, system.ColorTerminalGreen)
		}

		for y := 11; y < height; y = y + 11 {
			go d.DrawHoriz(y, 0, width, system.ColorTerminalGreen)
		}

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/ws", hub.Handler())

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}


}