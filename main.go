package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app"
	"github.com/cakesmith/webrender/app/demo"
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

		ta := demo.TextArea(main.Rectangle)

		main.Add(ta.Component)

	})

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
