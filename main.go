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

		for y := 11; y < height; y = y + 11 {
			for x := 0; x < width; x++ {
				go d.DrawPixel(x, y, system.ColorBlack)
			}
		}

		for x := 8; x < width; x = x + 8 {
			for y := 0; y < height; y++ {
				go d.DrawPixel(x, y, system.ColorBlack)
			}
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