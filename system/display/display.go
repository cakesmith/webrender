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
	CharMap               charMap
}

func Bit(x, j int) bool {
	return !(x&(1<<uint(j)) == 0)
}

func (t *Terminal) print(ch int) {

	startX, startY := t.cursorX*t.charWidth, t.cursorY*t.charHeight

	stopX, stopY := startX+t.charWidth, startY+t.charHeight

	printMe := t.CharMap.get(ch)

	for y := 0; startY+y < stopY; y++ {
		for x := 0; startX+x < stopX; x++ {

			color := ColorTerminalGreen

			pixel := Bit(printMe[y], x)

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
		CharMap: make(map[int]character),
	}

	t.CharMap.init()

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
	ch.add(32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)         // space
	ch.add(33, 0, 0, 8, 0, 8, 8, 8, 8, 8, 8, 8)         // !
	ch.add(34, 0, 0, 0, 0, 0, 0, 0, 20, 20, 20, 0)      // "
	ch.add(35, 0, 0, 0, 0, 20, 62, 20, 62, 20, 0, 0)    // #
	ch.add(36, 0, 0, 8, 30, 40, 40, 28, 10, 10, 60, 8)  // $
	ch.add(37, 0, 0, 0, 33, 82, 36, 8, 18, 37, 66, 0)   // %
	ch.add(38, 0, 0, 92, 34, 82, 74, 12, 8, 20, 20, 8)  // &
	ch.add(39, 0, 0, 0, 0, 0, 0, 0, 0, 16, 16, 8)       // '
	ch.add(40, 0, 56, 12, 6, 3, 1, 1, 3, 6, 12, 56)     // (
	ch.add(41, 0, 7, 12, 24, 48, 32, 32, 48, 24, 12, 7) // )
	ch.add(42, 0, 0, 0, 73, 42, 28, 127, 28, 42, 73, 0) // *
	ch.add(43, 0, 0, 0, 0, 0, 8, 28, 8, 0, 0, 0)        // +
	ch.add(44, 0, 4, 8, 8, 0, 0, 0, 0, 0, 0, 0)         // ,
	ch.add(45, 0, 0, 0, 0, 0, 30, 0, 0, 0, 0, 0)        // -
	ch.add(46, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0)         // .
	ch.add(47, 0, 0, 1, 2, 4, 8, 16, 32, 64, 128, 0)    // /

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

	ch.add(58, 0, 0, 0, 0, 8, 0, 0, 0, 8, 0, 0) // :
	ch.add(59, 0, 4, 8, 8, 0, 0, 0, 8, 0, 0, 0) // ;
	ch.add(60, 0, 0, 16, 8, 4, 2, 1, 2, 4, 8, 16) // <
	ch.add(61, 0, 0, 0, 0, 0, 28, 0, 28, 0, 0, 0) // =
	ch.add(62, 0, 0, 2, 4, 8, 16, 32, 16, 8, 4, 2) // >
	ch.add(63, 0, 8, 0, 8, 8, 24, 96, 64, 68, 68, 56) // ?
	ch.add(64, 0, 0, 124, 2, 33, 89, 85, 85, 73, 34, 28) // @

	ch.add(65, 0, 0, 65, 65, 65, 127, 34, 34, 20, 28, 8)  // A
	ch.add(66, 0, 0, 15, 17, 33, 17, 15, 17, 33, 17, 15)  // B
	ch.add(67, 0, 0, 28, 34, 65, 1, 1, 1, 65, 34, 28)     // C
	ch.add(68, 0, 0, 15, 17, 33, 33, 33, 33, 33, 17, 15)  // D
	ch.add(69, 0, 0, 63, 1, 1, 1, 31, 1, 1, 1, 63)        // E
	ch.add(70, 0, 0, 1, 1, 1, 1, 31, 1, 1, 1, 63)         // F
	ch.add(71, 0, 0, 28, 34, 65, 57, 1, 1, 65, 34, 28)    // G
	ch.add(72, 0, 0, 33, 33, 33, 33, 63, 33, 33, 33, 33)  // H
	ch.add(73, 0, 0, 31, 4, 4, 4, 4, 4, 4, 4, 31)         // I
	ch.add(74, 0, 0, 30, 33, 33, 32, 32, 32, 32, 32, 32)  // J
	ch.add(75, 0, 0, 65, 33, 17, 11, 5, 9, 17, 33, 65)    // K
	ch.add(76, 0, 0, 63, 1, 1, 1, 1, 1, 1, 1, 1)          // L
	ch.add(77, 0, 0, 65, 65, 65, 65, 73, 93, 119, 99, 99) // M
	ch.add(78, 0, 0, 65, 65, 97, 113, 89, 77, 71, 67, 67) // N
	ch.add(79, 0, 0, 12, 18, 33, 33, 33, 33, 33, 18, 12)  // O
	ch.add(80, 0, 0, 1, 1, 1, 1, 31, 33, 33, 33, 31)      // P
	ch.add(81, 0, 0, 44, 18, 41, 33, 33, 33, 33, 18, 12)  // Q
	ch.add(82, 0, 0, 33, 17, 9, 5, 31, 33, 33, 17, 15)    // R
	ch.add(83, 0, 0, 15, 16, 32, 32, 30, 1, 1, 2, 60)     // S
	ch.add(84, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 127)         // T
	ch.add(85, 0, 0, 12, 18, 33, 33, 33, 33, 33, 33, 33)  // U
	ch.add(86, 0, 0, 8, 28, 54, 34, 99, 65, 65, 65, 65)   // V
	ch.add(87, 0, 0, 65, 99, 119, 85, 93, 73, 73, 65, 65) // W
	ch.add(88, 0, 0, 65, 99, 54, 28, 8, 28, 54, 99, 65)   // X
	ch.add(89, 0, 0, 8, 8, 8, 8, 20, 34, 99, 65, 65)      // Y
	ch.add(90, 0, 0, 63, 1, 2, 4, 8, 16, 32, 63, 0)       // Z

	ch.add(91, 0, 0, 7, 1, 1, 1, 1, 1, 1, 1, 7)          // [
	ch.add(92, 0, 0, 128, 64, 32, 16, 8, 4, 2, 1, 0)     // \
	ch.add(93, 0, 0, 56, 32, 32, 32, 32, 32, 32, 32, 56) // ]
	ch.add(94, 0, 0, 0, 0, 0, 0, 0, 0, 34, 20, 8)        // ^
	ch.add(95, 0, 126, 0, 0, 0, 0, 0, 0, 0, 0, 0)        // _

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
