package output

import (
	"fmt"
	"github.com/cakesmith/webrender/websocket"
	"image"
	"image/color"
	"io"
	"strings"
)

var (
	ColorBackground = color.RGBA{
		40, 40, 40, 0,
	}

	ColorTerminalGreen = color.RGBA{
		51, 255, 51, 0,
	}

	ColorBlack = color.RGBA{
		0, 0, 0, 0,
	}

	ColorWhite = color.RGBA{
		255, 255, 255, 0,
	}

	ColorRed = color.RGBA{
		200, 0, 0, 0,
	}
)

type Drawer interface {
	DrawRectangle(rectangle image.Rectangle, rgba color.RGBA)
	//DrawVert(x, y1, y2 int, rgba color.RGBA)
	//DrawHoriz(x1, x2, y int, rgba color.RGBA)
	DrawPixel(x, y int, rgba color.RGBA)
	//DrawLine(x1, y1, x2, y2 int, rgba color.RGBA)
	//DrawCircle(cx, cy, r int, rgba color.RGBA)
}

type Terminal struct {
	io.Writer
}

//func (t *Terminal) Set(x, y int, c color.Color) {
//	t.DrawRectangle(x, y, 1, 1, c)
//}

func (t *Terminal) DrawRectangle(rect image.Rectangle, c color.Color) {
	r, g, b, a := c.RGBA()
	r8 := uint8(r)
	g8 := uint8(g)
	b8 := uint8(b)
	a8 := uint8(a)

	w := rect.Max.X - rect.Min.X
	h := rect.Max.Y - rect.Min.Y

	str := fmt.Sprintf("%v %v %v %v %v:%v:%v:%v", rect.Min.X, rect.Min.Y, w, h, r8, g8, b8, a8)
	t.Writer.Write(websocket.Command{
		Name:   "r",
		Params: strings.Split(str, " "),
	}.MakePacket())

}

//func (t *Terminal) DrawVert(x, y1, y2 int, color color.Color) error {
//	if y1 > y2 {
//		y2, y1 = y1, y2
//	}
//	return t.DrawRectangle(x, y1, 1, y2-y1, color)
//}
//
//func (t *Terminal) DrawHoriz(x1, x2, y int, color color.Color) error {
//	if x1 > x2 {
//		x2, x1 = x1, x2
//	}
//	return t.DrawRectangle(x1, y, x2-x1, 1, color)
//}
//
func (t *Terminal) DrawPixel(x, y int, color color.Color) {
	t.DrawRectangle(image.Rect(x, y, x+1, y+1), color)
}

//
//func (t *Terminal) DrawLine(x1, y1, x2, y2 int, color color.Color) error {
//
//	var dx, dy, a, b, adyMinusbdx int
//
//	adyMinusbdx = 0
//
//	dx = x2 - x1
//
//	if dx == 0 {
//		return t.DrawVert(x1, y1, y2, color)
//	}
//
//	dy = y2 - y1
//
//	if dy == 0 {
//		return t.DrawHoriz(x1, x2, y1, color)
//	}
//
//	a, b = 0, 0
//
//	if dx < 0 {
//
//		if dy < 0 {
//			// dx < 0 and dy < 0
//			// swap both points
//
//			return t.DrawLine(x2, y2, x1, y1, color)
//
//		} else {
//			// dx < 0 and dy > 0
//			// this is only subtly different than when dx > 0
//			// because we draw the pixel at x1-a instead of x1+a
//
//			dx = -dx
//
//			for (a <= dx) && (b <= dy) {
//
//				err := t.DrawPixel(x1-a, y1+b, color)
//				if err != nil {
//					return err
//				}
//
//				if adyMinusbdx < 0 {
//					a = a + 1
//					adyMinusbdx = adyMinusbdx + dy
//				} else {
//					b = b + 1
//					adyMinusbdx = adyMinusbdx - dx
//				}
//
//			}
//
//			return nil
//		}
//
//	} else {
//
//		if dy < 0 {
//			// dx > 0 and dy < 0
//			// swap both points
//
//			return t.DrawLine(x2, y2, x1, y1, color)
//
//		}
//
//	}
//
//	for (a <= dx) && (b <= dy) {
//
//		err := t.DrawPixel(x1+a, y1+b, color)
//		if err != nil {
//			return err
//		}
//
//		if adyMinusbdx < 0 {
//			a = a + 1
//			adyMinusbdx = adyMinusbdx + dy
//
//		} else {
//			b = b + 1
//			adyMinusbdx = adyMinusbdx - dx
//		}
//
//	}
//
//	return nil
//}
//
//func (t *Terminal) TestPattern() error {
//
//	err := t.DrawLine(t.Width/2, t.Height, 0, t.Height/2, color.White)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawLine(0, t.Height/2, t.Width/2, 0, color.White)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawLine(t.Width/2, 0, t.Width, t.Height/2, color.White)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawLine(t.Width, t.Height/2, t.Width/2, t.Height, color.White)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawCircle(t.Width/2, t.Height/2, 100, color.TerminalGreen)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawCircle(t.Width/2, t.Height/2, 99, color.Background)
//	if err != nil {
//		return err
//	}
//
//	err = t.DrawVert(t.Width/2, t.Height, t.Height/2, color.White)
//	if err != nil {
//		return err
//	}
//
//	return t.DrawPixel(t.Width/2, t.Height/2, color.Black)
//
//}
//
//func (t *Terminal) Clear(color color.Color) error {
//	return t.DrawRectangle(0, 0, t.Width, t.Height, color)
//}
//
//// DrawCircle Draws a color filled circle of radius r around cx, cy
//// thickness of 0 = filled
//func (t *Terminal) DrawCircle(cx, cy, r int, color color.Color) error {
//
//	for dy := -r; (dy - 1) < r; dy++ {
//
//		z := sqrt((r * r) - (dy * dy))
//
//		err := t.DrawLine(cx-z, cy+dy, cx+z, cy+dy, color)
//		if err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
//
//func sqrt(x int) int {
//	return int(math.Sqrt(float64(x)))
//}
