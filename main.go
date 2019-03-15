package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/system/display"
	"github.com/cakesmith/webrender/websocket"
	"math"
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
	//
	charWidth  = 8
	charHeight = 11
	//
	//d *display.Terminal
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

func (g *grid) bounds(n int) bounds {

	x := int(math.Mod(float64(n), float64(charWidth)))
	y := int(math.Floor(float64(n) / float64(charWidth)))

	return g.calcBounds(x, y)

}

func (g *grid) calcBounds(x, y int) bounds {
	left := g.bottomLeftX + (charWidth * x)
	right := left + charWidth
	bottom := g.bottomLeftY + (charHeight * y)
	top := bottom + charHeight

	return bounds{left, right, bottom, top}
}

func (g *grid) which(mx, my int) int {

	for y := 0; y < charHeight; y++ {
		for x := 0; x < charWidth; x++ {

			b := g.calcBounds(x, y)

			if mx > b.left && mx < b.right && my > b.bottom && my < b.top {
				return x + (charWidth * y)
			}
		}
	}
	return -1
}

func (g *grid) draw(t *display.Terminal, color display.Color) error {

	//TODO quick and dirty, use charWidth and charHeight instead of 8 and 11

	x1 := 28 * charWidth
	y1 := 10 * charHeight
	w1 := 36 * charWidth
	h1 := 21*charHeight + 1

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

	return nil
}

func (g *grid) drawChar(t *display.Terminal, ch int) {

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

	skipKeys := []int{
		20, // caps lock
		16, // shift
	}



	events := &websocket.Events{

		OnKeypress: func(key int) {

			fmt.Printf("keypress: %v\n", key)

			for _, k := range skipKeys {
				if key == k {
					fmt.Println("skipping key")
					return
				}
			}

			d.PrintChar(key)

			grid.reset()

			printMe := d.CharMap[key]

			for y := 0; y < height/charHeight; y++ {
				for x := 0; x < width/charWidth; x++ {

					pixel := display.Bit(printMe[y], x)

					if pixel {

						n := y*charWidth + x
						grid.bitmap[n] = true

						b := grid.bounds(n)

						d.DrawRectangle(b.left-1, b.bottom-1, charWidth+1, charHeight+1, display.ColorBlack)
						d.DrawRectangle(b.left, b.bottom, charWidth-1, charHeight-1, display.ColorTerminalGreen)

						// fix border

						dx := grid.bottomLeftX - 1
						dy := grid.bottomLeftY - 1
						dbx := charWidth*charWidth + dx
						dby := charHeight*charHeight + dy

						d.DrawLine(dx, dy, dx, dby, display.ColorTerminalGreen)
						d.DrawLine(dx, dy, dbx, dy, display.ColorTerminalGreen)
						d.DrawLine(dbx, dy, dbx, dby, display.ColorTerminalGreen)
						d.DrawLine(dx, dby, dbx, dby, display.ColorTerminalGreen)

					}

				}
			}

		},

		OnClick: func(btn, x, y int) {

			if resetBtn.is(x, y) {
				grid.reset()
			}

			w := grid.which(x, y)

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

				b := grid.bounds(w)

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

		OnReady: func(width, height int) {

		},
	}

	//hub.OnRegister = func(client *websocket.Client) {

		//
		//
		//d = display.New(width, height, charWidth, charHeight, client)
		//
		//d.Clear(display.ColorBackground)
		//
		//grid.draw(d, display.ColorTerminalGreen)
		//
		//resetBtn.draw(d)

	//}

	client := &websocket.Client{

		OnRecv: func(cmd []byte) {
			split := strings.Split(string(cmd), " ")

			switch split[0] {

			// keypress
			case "k":

				k, err := strconv.Atoi(string(split[1]))
				if err != nil {
					log.WithField("command", string(cmd)).Error(err)
					return
				}

				c.Events.OnKeypress(k)

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

				if c.Events.OnClick != nil {
					c.Events.OnClick(btn, x, y)
				}

			}

		},
	}

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", client.Handler)

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
