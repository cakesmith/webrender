package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/system/display"
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
		log.Fatal("$PORT environment variable must be set")
	}

	hub, err := websocket.NewHub()
	if err != nil {
		log.Fatal(err)
	}

	hub.OnRegister = func(client *websocket.Client) {

		d := display.Terminal{
			Writer: client,
			Width: width,
			Height: height,
		}

		d.Clear(display.ColorBackground)

		d.CharGrid(8, 11, display.ColorTerminalGreen)

		d.TestPattern()

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/ws", hub.Handler())

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}


}