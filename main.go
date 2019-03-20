package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app"
	"github.com/cakesmith/webrender/app/component"
	"github.com/cakesmith/webrender/output"
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
		log.Fatal("$PORT environment variable must be set")
	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		// Each connection gets its own client
		client := websocket.NewClient()

		// and therefore, its own app.
		main := &app.App{
			Client: client,
		}

		// handle the rest of the request, with actions
		main.MakeHandler(main)(w, r)

		taBorder := component.Border{
			output.ColorTerminalGreen,
			1,
		}

		ta := component.NewTextArea(output.ColorBackground, output.ColorTerminalGreen, taBorder, 8, 11, 0, 0, main.Container.Max.X, main.Container.Max.Y)

		main.Add(ta.Component)

	})

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
