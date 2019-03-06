package system

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
)

var (
	log = logrus.New()
)

type Command string

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

type DisplayWriter struct {
	io.Writer
}

func (display *DisplayWriter) Send(command Command, params ...string) error {
	packet := []byte(strings.Join(append([]string{string(command)}, params...), " "))
	n, err := display.Writer.Write(packet)
	if len(packet) != n {
		return errors.Wrap(err, fmt.Sprintf("len %v != %v", n, len(packet)))
	}
	return err
}

func (display *DisplayWriter) DrawVert(x, y1, y2 int, color Color) error {
	return display.DrawRectangle(x, y1, 1, y2-y1, color)
}

func (display *DisplayWriter) DrawHoriz(y, x1, x2 int, color Color) error {
	return display.DrawRectangle(x1, y, x2-x1, 1, color)
}



func (display *DisplayWriter) DrawRectangle(x1, y1, w, h int, color Color) error {

	sx1 := strconv.Itoa(x1)
	sy1 := strconv.Itoa(y1)
	sw := strconv.Itoa(w)
	sh := strconv.Itoa(h)
	sc := color.String()

	return display.Send("r", sx1, sy1, sw, sh, sc)

}

func (display *DisplayWriter) DrawPixel(x, y int, color Color) error {
	return display.DrawRectangle(x, y, 1, 1, color)
}

