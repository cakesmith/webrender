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
	log = logrus.New()

	//TODO use html template to pass these values to index.html
	width  = 512
	height = 330

	charWidth  = 8
	charHeight = 11

	d *display.Terminal
)

type grid struct {
	bottomLeftX int
	bottomLeftY int
	bitmap      map[int]bool
}

type bounds struct {
	left, right, bottom, top int
}

func (g *grid) reset() {
	g.bitmap = make(map[int]bool)
	g.draw(d, display.ColorTerminalGreen)
}

func (g *grid) String() string {

	var n []string

	for y := 0; y < charHeight; y++ {

		var ch byte

		for x := 0; x < charWidth; x++ {

			i := x + charWidth*y

			if g.bitmap[i] {
				ch = ch | (1 << uint(x))
			}

		}

		n = append(n, strconv.Itoa(int(ch)))
	}

	return strings.Join(n, ", ")
}

func (g *grid) which(mx, my int) (int, bounds) {

	for y := 0; y < charHeight; y++ {
		for x := 0; x < charWidth; x++ {

			left := g.bottomLeftX + (charWidth * x)
			right := left + charWidth
			bottom := g.bottomLeftY + (charHeight * y)
			top := bottom + charHeight

			if mx > left && mx < right && my > bottom && my < top {
				return x + (charWidth * y), bounds{left, right, bottom, top}
			}
		}
	}
	return -1, bounds{}
}

func (g *grid) draw(t *display.Terminal, color display.Color) error {

	//TODO quick and dirty, use charWidth and charHeight instead of 8 and 11

	x1 := 28 * charWidth
	y1 := 10 * charHeight
	w1 := 36 * charWidth
	h1 := 21 * charHeight

	t.DrawRectangle(28*charWidth, 10*charHeight, 8*charWidth, 11*charHeight, display.ColorBackground)

	for x := charWidth; x < t.Width; x = x + charWidth {
		if x < 28*8 || x > 36*8 {
			continue
		}

		err := t.DrawVert(x, y1, h1, color)

		if err != nil {
			return err
		}
	}

	for y := charHeight; y < t.Height; y = y + charHeight {
		if y < 10*11 || y > 21*11 {
			continue
		}

		err := t.DrawHoriz(x1, w1, y, color)

		if err != nil {
			return err
		}
	}

	//t.DrawRectangle(0, 0, t.Width, 11*10, display.ColorBackground)
	//t.DrawRectangle(0, t.Height, t.Width, (-11*9)+1, display.ColorBackground)
	//t.DrawRectangle(0, 0, 28*8, t.Height, display.ColorBackground)
	//t.DrawRectangle((36*8)+1, 0, t.Width-((36*8)+1), t.Height, display.ColorBackground)

	return nil
}

type button struct {
	x, y, width, height int
	color               display.Color
	border              display.Color
	text                string
}

func (b button) is(mx, my int) bool {
	if mx > b.x && my > b.y && mx < b.x+b.height && my < b.y+b.width {
		return true
	}
	return false
}

func (b button) draw(d *display.Terminal) {
	d.DrawRectangle(b.x, b.y, b.width, b.height, b.border)
	d.DrawRectangle(b.x+1, b.y+1, b.width-2, b.height-2, b.color)
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

	grid := &grid{
		bitmap:      make(map[int]bool),
		bottomLeftX: (28 * charWidth) + 1,
		bottomLeftY: (10 * charHeight) + 1,
	}

	resetBtn := button{
		x:      320,
		y:      154,
		width:  48,
		height: 33,
		border: display.ColorTerminalGreen,
		color:  display.ColorRed,
	}

	events := &websocket.Events{

		OnKeypress: func(key int) {

			fmt.Printf("keypress: %v\n", key)
			d.PrintChar(key)

		},

		OnClick: func(btn, x, y int) {

			if resetBtn.is(x, y) {
				grid.reset()
			}

			w, b := grid.which(x, y)

			if w != -1 {

				grid.bitmap[w] = !grid.bitmap[w]

				var color display.Color

				if grid.bitmap[w] {
					color = display.ColorTerminalGreen
				} else {
					color = display.ColorBackground
				}

				alt := display.ColorTerminalGreen

				if alt == color {
					alt = display.ColorBlack
				}

				d.DrawRectangle(b.left-1, b.bottom-1, charWidth+1, charHeight+1, alt)
				d.DrawRectangle(b.left, b.bottom, charWidth-1, charHeight-1, color)

				// fix border

				dx := grid.bottomLeftX - 1
				dy := grid.bottomLeftY - 1
				dbx := charWidth*charWidth + dx
				dby := charHeight*charHeight + dy

				d.DrawLine(dx, dy, dx, dby, display.ColorTerminalGreen)
				d.DrawLine(dx, dy, dbx, dy, display.ColorTerminalGreen)
				d.DrawLine(dbx, dy, dbx, dby, display.ColorTerminalGreen)
				d.DrawLine(dx, dby, dbx, dby, display.ColorTerminalGreen)

				fmt.Println(grid.String())

			}

		},
	}

	hub.OnRegister = func(client *websocket.Client) {

		d = display.New(width, height, charWidth, charHeight, client)

		d.Clear(display.ColorBackground)

		grid.draw(d, display.ColorTerminalGreen)

		resetBtn.draw(d)

	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", hub.Handler(events))

	logrus.WithField("port", port).Info("ready")

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
