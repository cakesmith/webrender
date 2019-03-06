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
	R int
	G int
	B int
}

var (

	ColorBackground = Color{
		40, 40, 40,
	}

	ColorTerminalGreen = Color{
		51, 255, 51,
	}

	ColorBlack = Color{
		0,0,0,
	}

	ColorWhite = Color{
		255,255,255,
	}
)

func (p Color) String() string {
	return fmt.Sprintf("%v-%v-%v", p.R, p.G, p.B)
}

type Terminal struct {
	io.Writer
}

type Command struct {
	Name string
	Params []string
}

func (c Command) makePacket() []byte {
	return []byte(strings.Join(append([]string{c.Name}, c.Params...), " "))
}

func (terminal *Terminal) send(c Command) error {
	packet := c.makePacket()
	n, err := terminal.Writer.Write(packet)
	if len(packet) != n {
		return errors.Wrap(err, fmt.Sprintf("len %v != %v", n, len(packet)))
	}
	return err
}

func (terminal *Terminal) DrawRectangle(x1, y1, w, h int, color Color) error {

	sx1 := strconv.Itoa(x1)
	sy1 := strconv.Itoa(y1)
	sw := strconv.Itoa(w)
	sh := strconv.Itoa(h)
	sc := color.String()

	return terminal.send(Command{
		Name: "r",
		Params: []string{
			sx1, sy1, sw, sh, sc,
		},
	})


}

func (terminal *Terminal) DrawVert(x, y1, y2 int, color Color) error {
	return terminal.DrawRectangle(x, y1, 1, y2-y1, color)
}

func (terminal *Terminal) DrawHoriz(x1, x2, y int, color Color) error {
	return terminal.DrawRectangle(x1, y, x2-x1, 1, color)
}

func (terminal *Terminal) DrawPixel(x, y int, color Color) error {
	return terminal.DrawRectangle(x, y, 1, 1, color)
}

func (terminal *Terminal) DrawLine(x1, y1, x2, y2 int, color Color) error {

	var dx, dy, a, b, adyMinusbdx int

	adyMinusbdx = 0

	dx = x2 - x1

	if dx == 0 {
		return terminal.DrawVert(x1, y1, y2, color)
	}

	dy = y2 - y1

	if dy == 0 {
		return terminal.DrawHoriz(x1, x2, y1, color)
	}

	a, b = 0, 0

	if dx < 0 {

		if dy < 0 {
			// dx < 0 and dy < 0
			// swap both points

			return terminal.DrawLine(x2, y2, x1, y1, color)

		} else {
			// dx < 0 and dy > 0
			// this is only subtly different than when dx > 0
			// because we draw the pixel at x1-a instead of x1+a

			dx = -dx

			for (a <= dx) && (b <= dy) {

				err := terminal.DrawPixel(x1-a, y1+b, color)
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

			return terminal.DrawLine(x2, y2, x1, y1, color)

		}

	}

	for (a <= dx) && (b <= dy) {

		err := terminal.DrawPixel(x1+a, y1+b, color)
		if err != nil {
			return err
		}

		if adyMinusbdx < 0 {
			 a = a+1
			 adyMinusbdx = adyMinusbdx + dy

		} else {
			 b = b+1
			 adyMinusbdx = adyMinusbdx - dx
		}

	}

	return nil
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
func (terminal *Terminal) DrawCircle(cx, cy, r int, color Color) error {

	r = Abs(r)

	for dy := -r; (dy - 1) < r; dy++ {

		z := Sqrt((r*r) - (dy*dy))

		err := terminal.DrawLine(cx-z, cy+dy, cx+z, cy+dy, color)
		if err != nil {
			return err
		}

	}

	return nil
}