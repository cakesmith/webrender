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

		d := display.Terminal{client}

		d.DrawRectangle(0, 0, width, height, display.ColorBackground)

		//for x := 8; x < width; x = x + 8 {
		//	d.DrawVert(x, 0, height, display.ColorTerminalGreen)
		//}
		//
		//for y := 11; y < height; y = y + 11 {
		//	d.DrawHoriz(0, width, y, display.ColorTerminalGreen)
		//}

		d.DrawLine(20, 100, 50, 60, display.ColorWhite)

		d.DrawCircle(width/2, height/2, 100, display.ColorTerminalGreen)

		d.DrawCircle(width/2, height/2, 99, display.ColorBackground)

		d.DrawPixel(width/2, height/2, display.ColorBlack)

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/ws", hub.Handler())

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}


}