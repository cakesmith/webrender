package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/system/display"
	"github.com/cakesmith/webrender/websocket"
	"net/http"
	"os"
)

var (
	log    = logrus.New()
	width  = 512
	height = 330
)

type charGrid struct {
	bottomLeftX int
	bottomLeftY int
	width       int // in characters
	height      int // in characters
	grid        map[int]bool
}

var d display.Terminal

type bounds struct {
	left, right, bottom, top int
}

func (chars charGrid) which(mx, my int) (int, bounds) {
	for y := 0; y < chars.height; y++ {
		for x := 0; x < chars.width; x++ {

			left := chars.bottomLeftX + (8 * x)
			right := left + 8
			bottom := chars.bottomLeftY + (11 * y)
			top := bottom + 11

			if mx > left && mx < right && my > bottom && my < top {
				return x + (8 * y), bounds{left, right, bottom, top}
			}
		}
	}
	return -1, bounds{}
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	}

	hub, err := websocket.NewHub()
	if err != nil {
		log.Fatal(err)
	}

	chars := charGrid{
		width:       8,
		height:      11,
		grid:        make(map[int]bool),
		bottomLeftX: (28 * 8) + 1,
		bottomLeftY: (10 * 11) + 1,
	}

	events := &websocket.Events{

		OnClick: func(btn, x, y int) {
			w, b := chars.which(x, y)
			if w != -1 {
				chars.grid[w] = !chars.grid[w]
				var color display.Color
				if chars.grid[w] {
					color = display.ColorTerminalGreen
				} else {
					color = display.ColorBackground
				}
				d.DrawRectangle(b.left, b.bottom, 7, 10, color)

			}
		},
	}

	hub.OnRegister = func(client *websocket.Client) {

		d = display.Terminal{
			Writer: client,
			Width:  width,
			Height: height,
		}

		d.Clear(display.ColorBackground)

		d.CharGrid(8, 11, display.ColorTerminalGreen)

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", hub.Handler(events))

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
