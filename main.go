package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app/component"
	"github.com/cakesmith/webrender/output"
	"github.com/cakesmith/webrender/websocket"
	"image"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	log = logrus.New()

	//TODO use html template to pass these values to index.html
	//width  = 512
	//height = 330

)

type App struct {
	*component.Container
	*websocket.Client
}

func (app *App) OnRecv(cmd []byte) {

	split := strings.Split(string(cmd), " ")

	switch split[0] {

	// keypress
	case "k":

		key, err := strconv.Atoi(string(split[1]))
		if err != nil {
			log.WithField("command", string(cmd)).Error(err)
			return
		}

		app.Container.OnKeypress(key)

	// mouse click
	case "mc":

		btn, err := strconv.Atoi(string(split[1]))
		if err != nil {
			log.WithField("command", string(cmd)).Error(err)
			return
		}

		x, err := strconv.Atoi(string(split[2]))
		if err != nil {
			log.WithField("command", string(cmd)).Error(err)
			return
		}

		y, err := strconv.Atoi(string(split[3]))
		if err != nil {
			log.WithField("command", string(cmd)).Error(err)
			return
		}

		app.Container.OnClick(btn, x, y)

	}

}

func (app *App) OnRegister() {
	width := 512
	height := 330

	fmt.Println("creating main container")

	app.Container = &component.Container{
		Rectangle:  image.Rect(0, 0, width, height),
		Terminal:   output.Terminal{Writer: app.Client},
		Components: []*component.Component{},
	}

	var redBorder = component.Border{
		Color:     output.ColorRed,
		Thickness: 1,
	}

	btn := component.NewButton(output.ColorBackground, redBorder, 0, 0, 50, 75)
	app.Container.Add(btn.Component)

	btn.Component.Init = func() {
		btn.Draw()
	}

	btn.Component.OnClick = func(btn, x, y int) {
		fmt.Printf("click btn %v received: x %v - y %v\n", btn, x, y)
	}

	app.Container.Init()

}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	}

	client := websocket.NewClient()

	app := &App{
		Client: client,
	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", client.MakeHandler(app))

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
