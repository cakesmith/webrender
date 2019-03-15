package output

import (
	"fmt"
	"github.com/cakesmith/webrender/output/color"
	"github.com/cakesmith/webrender/websocket"
	"github.com/pkg/errors"
	"io"
	"math"
	"strings"
)

type Drawer interface {
	DrawRectangle(xy, y1, w, h uint, color color.Color) error
	DrawVert(x, y1, y2 uint, color color.Color) error
	DrawHoriz(x1, x2, y uint, color color.Color) error
	DrawPixel(x, y uint, color color.Color) error
	DrawLine(x1, y1, x2, y2 uint, color color.Color) error
	DrawCircle(cx, cy, r uint, color color.Color) error
}

type Terminal struct {
	io.Writer
	Width, Height uint
}

func (t *Terminal) send(packet []byte) error {
	n, err := t.Writer.Write(packet)
	if len(packet) != n {
		return errors.Wrap(err, fmt.Sprintf("len %v != %v", n, len(packet)))
	}
	return err
}

func (t *Terminal) DrawRectangle(x1, y1, w, h uint, color color.Color) error {

	// flip y axis for canvas orientation.
	y1 = t.Height - y1
	h = -h

	str := fmt.Sprintf("%v %v %v %v %v", x1, y1, w, h, color)

	return t.send(websocket.Command{
		Name:   "r",
		Params: strings.Split(str, " "),
	}.MakePacket())

}

func (t *Terminal) DrawVert(x, y1, y2 uint, color color.Color) error {
	if y1 > y2 {
		y2, y1 = y1, y2
	}
	return t.DrawRectangle(x, y1, 1, y2-y1, color)
}

func (t *Terminal) DrawHoriz(x1, x2, y uint, color color.Color) error {
	if x1 > x2 {
		x2, x1 = x1, x2
	}
	return t.DrawRectangle(x1, y, x2-x1, 1, color)
}

func (t *Terminal) DrawPixel(x, y uint, color color.Color) error {
	return t.DrawRectangle(x, y, 1, 1, color)
}

func (t *Terminal) DrawLine(x1, y1, x2, y2 uint, color color.Color) error {

	var dx, dy, a, b, adyMinusbdx uint

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

	err := t.DrawLine(t.Width/2, t.Height, 0, t.Height/2, color.White)
	if err != nil {
		return err
	}

	err = t.DrawLine(0, t.Height/2, t.Width/2, 0, color.White)
	if err != nil {
		return err
	}

	err = t.DrawLine(t.Width/2, 0, t.Width, t.Height/2, color.White)
	if err != nil {
		return err
	}

	err = t.DrawLine(t.Width, t.Height/2, t.Width/2, t.Height, color.White)
	if err != nil {
		return err
	}

	err = t.DrawCircle(t.Width/2, t.Height/2, 100, color.TerminalGreen)
	if err != nil {
		return err
	}

	err = t.DrawCircle(t.Width/2, t.Height/2, 99, color.Background)
	if err != nil {
		return err
	}

	err = t.DrawVert(t.Width/2, t.Height, t.Height/2, color.White)
	if err != nil {
		return err
	}

	return t.DrawPixel(t.Width/2, t.Height/2, color.Black)

}

func (t *Terminal) Clear(color color.Color) error {
	return t.DrawRectangle(0, 0, t.Width, t.Height, color)
}

// DrawCircle Draws a color filled circle of radius r around cx, cy
// thickness of 0 = filled
func (t *Terminal) DrawCircle(cx, cy, r uint, color color.Color) error {

	for dy := -r; (dy - 1) < r; dy++ {

		z := sqrt((r * r) - (dy * dy))

		err := t.DrawLine(cx-z, cy+dy, cx+z, cy+dy, color)
		if err != nil {
			return err
		}

	}

	return nil
}

func sqrt(x uint) uint {
	return uint(math.Sqrt(float64(x)))
}
