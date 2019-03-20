package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app/component"
	"github.com/cakesmith/webrender/output"
	"github.com/cakesmith/webrender/websocket"
	"image"
	"strconv"
	"strings"
)

var (
	log = logrus.New()
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

	w, err := app.Client.Request((websocket.Command{
		Name: "w",
	}).MakePayload())

	if err != nil {
		log.Error(err)
	}

	h, err := app.Client.Request((websocket.Command{
		Name: "h",
	}).MakePayload())

	if err != nil {
		log.Error(err)
	}

	width, err := strconv.Atoi(string(w))

	if err != nil {
		log.Error(err)
	}

	height, err := strconv.Atoi(string(h))

	if err != nil {
		log.Error(err)
	}

	log.WithFields(logrus.Fields{"width": width, "height": height}).Println("creating main container")

	app.Container = &component.Container{
		Rectangle: image.Rect(0, 0, width, height),
		Terminal:  output.Terminal{Writer: app.Client},
	}

}

