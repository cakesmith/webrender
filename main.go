package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/system/display"
	"github.com/cakesmith/webrender/websocket"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	log    = logrus.New()
	width  = 512
	height = 330
	charWidth = 8
	charHeight = 11
)

type charGrid struct {
	bottomLeftX int
	bottomLeftY int
	grid        map[int]bool
}

var d display.Terminal

type bounds struct {
	left, right, bottom, top int
}

func (chars charGrid) String() string {

	var n []string

	for y := 0; y < charHeight; y++ {

		var ch byte

		for x := 0; x < charWidth; x ++ {

			i := x + charWidth * y

			if chars.grid[i] {
				ch = ch | (1 << uint(x))
			}

		}

		n = append(n, strconv.Itoa(int(ch)))
	}

	return strings.Join(n, ", ")
}

func (chars charGrid) which(mx, my int) (int, bounds) {

	for y := 0; y < charHeight; y++ {
		for x := 0; x < charWidth; x++ {

			left := chars.bottomLeftX + (charWidth * x)
			right := left + charWidth
			bottom := chars.bottomLeftY + (charHeight * y)
			top := bottom + charHeight

			if mx > left && mx < right && my > bottom && my < top {
				return x + (charWidth * y), bounds{left, right, bottom, top}
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
		grid:        make(map[int]bool),
		bottomLeftX: (28 * charWidth) + 1,
		bottomLeftY: (10 * charHeight) + 1,
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

				alt := display.ColorTerminalGreen

				if alt == color {
					alt = display.ColorBlack
				}

				d.DrawRectangle(b.left-1, b.bottom-1, charWidth + 1, charHeight + 1, alt)
				d.DrawRectangle(b.left, b.bottom, charWidth -1, charHeight - 1, color)

				dx := chars.bottomLeftX - 1
				dy := chars.bottomLeftY - 1
				dbx := charWidth * charWidth + dx
				dby := charHeight * charHeight + dy

				// fix border
				d.DrawLine(dx, dy, dx, dby, display.ColorTerminalGreen)
				d.DrawLine(dx, dy, dbx, dy, display.ColorTerminalGreen)
				d.DrawLine(dbx, dy, dbx, dby, display.ColorTerminalGreen)
				d.DrawLine(dx, dby, dbx, dby, display.ColorTerminalGreen)

				fmt.Println(chars.String())

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

		d.CharGrid(charWidth, charHeight, display.ColorTerminalGreen)

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", hub.Handler(events))

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
