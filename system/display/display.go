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
	Width, Height int
}

type Command struct {
	Name string
	Params []string
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
	return t.DrawRectangle(x, y1, 1, y2-y1, color)
}

func (t *Terminal) DrawHoriz(x1, x2, y int, color Color) error {
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
			 a = a+1
			 adyMinusbdx = adyMinusbdx + dy

		} else {
			 b = b+1
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

	return t.DrawPixel(t.Width/2, t.Height/2, ColorBlack)

}

func (t *Terminal) Clear(color Color) error {
	return t.DrawRectangle(0, 0, t.Width, t.Height, ColorBackground)

}

func (t *Terminal) CharGrid(charWidth, charHeight int, color Color) error {
	
	for x := charWidth; x < t.Width; x = x + charWidth {
		err := t.DrawVert(x, 0, t.Height, color)
		if err != nil {
			return err
		}
	}

	for y := charHeight; y < t.Height; y = y + charHeight {
		err := t.DrawHoriz(0, t.Width, y, color)
		if err != nil {
			return err
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
func (t *Terminal) DrawCircle(cx, cy, r int, color Color) error {

	r = Abs(r)

	for dy := -r; (dy - 1) < r; dy++ {

		z := Sqrt((r*r) - (dy*dy))

		err := t.DrawLine(cx-z, cy+dy, cx+z, cy+dy, color)
		if err != nil {
			return err
		}

	}

	return nil
}


