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

		ta.Init = func() {

			ta.Buffer = []int{}

			//***** TEST PATTERN *****

			chw := ta.Width() / ta.CharWidth
			chh := ta.Height() / ta.CharHeight

			str := "ALL WORK AND NO PLAY MAKES JACK A DULL BOY. "

			log.Println(str)

			for i := 0; i+len(str) < chw*chh; i = i + len(str) {
				ta.PrintString(str)
			}

			end := chw - ta.CursorX

			for i := 0; i < end; i++ {
				ta.PrintString(str[i : i+1])
			}

			//************************

		}

		ta.Draw = func() {

			ta.DrawRectangle(ta.Component.Rectangle, ta.Border.Color)
			ta.DrawRectangle(ta.Component.Rectangle.Inset(ta.Border.Thickness), ta.Component.Color)

			ta.CursorX, ta.CursorY = 0, 0

			for _, chr := range ta.Buffer {
				ta.PrintChar(chr)
			}
		}

		ta.OnKeypress = func(key int) {
			log.WithField("key", key)
			ta.Buffer = append(ta.Buffer, key)
			ta.PrintChar(key)
		}
		
		main.Add(ta.Component)

	})

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
