package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app"
	"github.com/cakesmith/webrender/app/font"
	"github.com/cakesmith/webrender/websocket"
	"image"
	"image/draw"
	"net/http"
	"os"
)

var (
	log = logrus.New()

	//TODO use html template to pass these values to index.html
	//width  = 512
	//height = 330
	//
	//d *display.Terminal
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	}

	//grid := &grid{
	//	bitmap:      make(map[int]bool),
	//	bottomLeftX: (28 * charWidth) + 1,
	//	bottomLeftY: (10 * charHeight) + 1,
	//}

	//resetBtn := button{
	//	x:      320,
	//	y:      154,
	//	width:  48,
	//	height: 33,
	//	border: display.ColorTerminalGreen,
	//	color:  display.ColorRed,
	//}
	//
	//skipKeys := []int{
	//	20, // caps lock
	//	16, // shift
	//}

	//events := &websocket.Events{
	//
	//	OnKeypress: func(key int) {
	//
	//		fmt.Printf("keypress: %v\n", key)
	//
	//		for _, k := range skipKeys {
	//			if key == k {
	//				fmt.Println("skipping key")
	//				return
	//			}
	//		}
	//
	//		d.PrintChar(key)
	//
	//		grid.reset()
	//
	//		printMe := d.CharMap[key]
	//
	//		for y := 0; y < height/charHeight; y++ {
	//			for x := 0; x < width/charWidth; x++ {
	//
	//				pixel := display.Bit(printMe[y], x)
	//
	//				if pixel {
	//
	//					n := y*charWidth + x
	//					grid.bitmap[n] = true
	//
	//					b := grid.bounds(n)
	//
	//					d.DrawRectangle(b.left-1, b.bottom-1, charWidth+1, charHeight+1, display.ColorBlack)
	//					d.DrawRectangle(b.left, b.bottom, charWidth-1, charHeight-1, display.ColorTerminalGreen)
	//
	//					// fix border
	//
	//					dx := grid.bottomLeftX - 1
	//					dy := grid.bottomLeftY - 1
	//					dbx := charWidth*charWidth + dx
	//					dby := charHeight*charHeight + dy
	//
	//					d.DrawLine(dx, dy, dx, dby, display.ColorTerminalGreen)
	//					d.DrawLine(dx, dy, dbx, dy, display.ColorTerminalGreen)
	//					d.DrawLine(dbx, dy, dbx, dby, display.ColorTerminalGreen)
	//					d.DrawLine(dx, dby, dbx, dby, display.ColorTerminalGreen)
	//
	//				}
	//
	//			}
	//		}
	//
	//	},
	//
	//	OnClick: func(btn, x, y int) {
	//
	//		if resetBtn.is(x, y) {
	//			grid.reset()
	//		}
	//
	//		w := grid.which(x, y)
	//
	//		if w != -1 {
	//
	//			grid.bitmap[w] = !grid.bitmap[w]
	//
	//			var color display.Color
	//
	//			if grid.bitmap[w] {
	//				color = display.ColorTerminalGreen
	//			} else {
	//				color = display.ColorBackground
	//			}
	//
	//			alt := display.ColorTerminalGreen
	//
	//			if alt == color {
	//				alt = display.ColorBlack
	//			}
	//
	//			b := grid.bounds(w)
	//
	//			d.DrawRectangle(b.left-1, b.bottom-1, charWidth+1, charHeight+1, alt)
	//			d.DrawRectangle(b.left, b.bottom, charWidth-1, charHeight-1, color)
	//
	//			// fix border
	//
	//			dx := grid.bottomLeftX - 1
	//			dy := grid.bottomLeftY - 1
	//			dbx := charWidth*charWidth + dx
	//			dby := charHeight*charHeight + dy
	//
	//			d.DrawLine(dx, dy, dx, dby, display.ColorTerminalGreen)
	//			d.DrawLine(dx, dy, dbx, dy, display.ColorTerminalGreen)
	//			d.DrawLine(dbx, dy, dbx, dby, display.ColorTerminalGreen)
	//			d.DrawLine(dx, dby, dbx, dby, display.ColorTerminalGreen)
	//
	//			fmt.Println(grid.String())
	//
	//		}
	//
	//	},
	//
	//	OnReady: func(width, height int) {
	//
	//	},
	//}

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

	client := websocket.NewClient()

	client.OnRegister = func() {

		width := 512
		height := 330

		//cont := app.Container{
		//	App: font.Designer{},
		//}

	}

		//OnRecv: func(cmd []byte) {
		//	split := strings.Split(string(cmd), " ")
		//
		//	switch split[0] {
		//
		//	// keypress
		//	case "k":
		//
		//		k, err := strconv.Atoi(string(split[1]))
		//		if err != nil {
		//			log.WithField("command", string(cmd)).Error(err)
		//			return
		//		}
		//
		//		c.Events.OnKeypress(k)
		//
		//	// mouse click
		//	case "mc":
		//
		//		btn, err := strconv.Atoi(string(split[1]))
		//		if err != nil {
		//			log.WithField("command", string(cmd)).Error(err)
		//			return
		//		}
		//
		//		x, err := strconv.Atoi(string(split[2]))
		//		if err != nil {
		//			log.WithField("command", string(cmd)).Error(err)
		//			return
		//		}
		//
		//		y, err := strconv.Atoi(string(split[3]))
		//		if err != nil {
		//			log.WithField("command", string(cmd)).Error(err)
		//			return
		//		}
		//
		//		if c.Events.OnClick != nil {
		//			c.Events.OnClick(btn, x, y)
		//		}
		//
		//	}
		//
		//},


	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.HandleFunc("/ws", client.Handler)

	logrus.WithField("port", port).Info("ready")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Error(err)
	}

}
