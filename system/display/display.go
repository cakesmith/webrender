package display

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"io"
	"math"
	"strconv"
	"strings"
)

var (
	log = logrus.New()
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

var (
	ColorBackground = Color{
		40, 40, 40,
	}

	ColorTerminalGreen = Color{
		51, 255, 51,
	}

	ColorBlack = Color{
		0, 0, 0,
	}

	ColorWhite = Color{
		255, 255, 255,
	}

	ColorRed = Color{
		200, 0, 0,
	}
)

func (p Color) String() string {
	return fmt.Sprintf("%v-%v-%v", p.R, p.G, p.B)
}

type Terminal struct {
	io.Writer
	Width, Height         int
	cursorX, cursorY      int
	charWidth, charHeight int
	charMap               charMap
}

func bit(x, j int) bool {
	return !(x&(1<<uint(j)) == 0)
}

func (t *Terminal) print(ch int) {

	startX, startY := t.cursorX*t.charWidth, t.cursorY*t.charHeight

	stopX, stopY := startX+t.charWidth, startY+t.charHeight

	printMe := t.charMap.get(ch)

	for y := 0; startY+y < stopY; y++ {
		for x := 0; startX+x < stopX; x++ {

			color := ColorTerminalGreen

			pixel := bit(printMe[y], x)

			if !pixel {
				color = ColorBackground
			}

			t.DrawPixel(startX+x, startY+y, color)

		}
	}

}

//Advances the cursor to the beginning of the next line.
func (t *Terminal) println() {
	if t.cursorY < t.Height/t.charHeight {
		t.cursorX = 0
		t.cursorY++
	}
}

func (t *Terminal) PrintChar(ch int) {

	if t.cursorX < t.Width/t.charWidth {

		t.print(ch)
		t.cursorX++

	} else {
		if t.cursorY < t.Height/t.charHeight {
			t.cursorX = 0
			t.cursorY = 0
			t.print(ch)
			t.cursorX = 1
		} else {
			t.println()
		}
	}
}

type Command struct {
	Name   string
	Params []string
}

func New(width, height, charWidth, charHeight int, writer io.Writer) *Terminal {

	t := &Terminal{
		Writer: writer,
		Width:  width, Height: height,
		charWidth: charWidth, charHeight: charHeight,
		charMap: make(map[int]character),
	}

	t.charMap.init()

	return t
}

func (c Command) makePacket() []byte {
	return []byte(strings.Join(append([]string{c.Name}, c.Params...), " "))
}

func (t *Terminal) send(c Command) error {
	packet := c.makePacket()
	n, err := t.Writer.Write(packet)
	if len(packet) != n {
		return errors.Wrap(err, fmt.Sprintf("len %v != %v", n, len(packet)))
	}
	return err
}

func (t *Terminal) DrawRectangle(x1, y1, w, h int, color Color) error {

	// flip y axis for canvas orientation.
	y1 = t.Height - y1
	h = -h

	sx1 := strconv.Itoa(x1)
	sy1 := strconv.Itoa(y1)
	sw := strconv.Itoa(w)
	sh := strconv.Itoa(h)
	sc := color.String()

	return t.send(Command{
		Name: "r",
		Params: []string{
			sx1, sy1, sw, sh, sc,
		},
	})

}

func (t *Terminal) DrawVert(x, y1, y2 int, color Color) error {
	if y1 > y2 {
		y2, y1 = y1, y2
	}
	return t.DrawRectangle(x, y1, 1, y2-y1, color)
}

func (t *Terminal) DrawHoriz(x1, x2, y int, color Color) error {
	if x1 > x2 {
		x2, x1 = x1, x2
	}
	return t.DrawRectangle(x1, y, x2-x1, 1, color)
}

func (t *Terminal) DrawPixel(x, y int, color Color) error {
	return t.DrawRectangle(x, y, 1, 1, color)
}

func (t *Terminal) DrawLine(x1, y1, x2, y2 int, color Color) error {

	var dx, dy, a, b, adyMinusbdx int

	adyMinusbdx = 0

	dx = x2 - x1

	if dx == 0 {
		return t.DrawVert(x1, y1, y2, color)
	}

	dy = y2 - y1

	if dy == 0 {
		return t.DrawHoriz(x1, x2, y1, color)
	}

	a, b = 0, 0

	if dx < 0 {

		if dy < 0 {
			// dx < 0 and dy < 0
			// swap both points

			return t.DrawLine(x2, y2, x1, y1, color)

		} else {
			// dx < 0 and dy > 0
			// this is only subtly different than when dx > 0
			// because we draw the pixel at x1-a instead of x1+a

			dx = -dx

			for (a <= dx) && (b <= dy) {

				err := t.DrawPixel(x1-a, y1+b, color)
				if err != nil {
					return err
				}

				if adyMinusbdx < 0 {
					a = a + 1
					adyMinusbdx = adyMinusbdx + dy
				} else {
					b = b + 1
					adyMinusbdx = adyMinusbdx - dx
				}

			}

			return nil
		}

	} else {

		if dy < 0 {
			// dx > 0 and dy < 0
			// swap both points

			return t.DrawLine(x2, y2, x1, y1, color)

		}

	}

	for (a <= dx) && (b <= dy) {

		err := t.DrawPixel(x1+a, y1+b, color)
		if err != nil {
			return err
		}

		if adyMinusbdx < 0 {
			a = a + 1
			adyMinusbdx = adyMinusbdx + dy

		} else {
			b = b + 1
			adyMinusbdx = adyMinusbdx - dx
		}

	}

	return nil
}

func (t *Terminal) TestPattern() error {

	err := t.DrawLine(t.Width/2, t.Height, 0, t.Height/2, ColorWhite)
	if err != nil {
		return err
	}

	err = t.DrawLine(0, t.Height/2, t.Width/2, 0, ColorWhite)
	if err != nil {
		return err
	}

	err = t.DrawLine(t.Width/2, 0, t.Width, t.Height/2, ColorWhite)
	if err != nil {
		return err
	}

	err = t.DrawLine(t.Width, t.Height/2, t.Width/2, t.Height, ColorWhite)
	if err != nil {
		return err
	}

	err = t.DrawCircle(t.Width/2, t.Height/2, 100, ColorTerminalGreen)
	if err != nil {
		return err
	}

	err = t.DrawCircle(t.Width/2, t.Height/2, 99, ColorBackground)
	if err != nil {
		return err
	}

	err = t.DrawVert(t.Width/2, t.Height, t.Height/2, ColorWhite)
	if err != nil {
		return err
	}

	return t.DrawPixel(t.Width/2, t.Height/2, ColorBlack)

}

func (t *Terminal) Clear(color Color) error {
	return t.DrawRectangle(0, 0, t.Width, t.Height, ColorBackground)

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Sqrt(x int) int {
	return int(math.Sqrt(float64(x)))

}

// DrawCircle Draws a color filled circle of radius r around cx, cy
// thickness of 0 = filled
func (t *Terminal) DrawCircle(cx, cy, r int, color Color) error {

	r = Abs(r)

	for dy := -r; (dy - 1) < r; dy++ {

		z := Sqrt((r * r) - (dy * dy))

		err := t.DrawLine(cx-z, cy+dy, cx+z, cy+dy, color)
		if err != nil {
			return err
		}

	}

	return nil
}

type character map[int]int

type charMap map[int]character

func (ch charMap) init() {

	log.Info("initializing character map...")

	// Blank square for non printable characters
	ch.add(0, 0, 0, 63, 63, 63, 63, 63, 63, 63, 63, 63)

	// Assign bitmap for each character in set
	ch.add(32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0) //

	ch.add(48, 0, 0, 12, 18, 33, 33, 45, 33, 33, 18, 12) // 0
	ch.add(49, 0, 0, 31, 4, 4, 4, 4, 4, 5, 6, 4)         // 1
	ch.add(50, 0, 0, 63, 2, 4, 8, 16, 32, 32, 33, 30)    // 2
	ch.add(51, 0, 0, 14, 17, 32, 16, 12, 16, 32, 17, 14) // 3
	ch.add(52, 0, 0, 16, 16, 16, 63, 17, 18, 20, 24, 16) // 4
	ch.add(53, 0, 0, 15, 16, 32, 32, 16, 15, 1, 1, 63)   // 5
	ch.add(54, 0, 0, 30, 33, 33, 33, 31, 1, 1, 1, 62)    // 6
	ch.add(55, 0, 0, 1, 2, 4, 8, 16, 32, 64, 127, 0)     // 7
	ch.add(56, 0, 0, 12, 18, 33, 18, 12, 18, 33, 18, 12) // 8
	ch.add(57, 0, 0, 30, 32, 32, 32, 62, 33, 33, 33, 30) // 9

}

func (ch charMap) get(i int) character {
	if i < 32 || i > 126 {
		i = 0
	}
	return ch[i]
}

func (ch charMap) add(index int, v ...int) {

	if len(v) < 11 {
		log.Fatalf("cannot add character %v", index)
	}

	ch[index] = map[int]int{
		0:  v[0],
		1:  v[1],
		2:  v[2],
		3:  v[3],
		4:  v[4],
		5:  v[5],
		6:  v[6],
		7:  v[7],
		8:  v[8],
		9:  v[9],
		10: v[10],
	}
}
