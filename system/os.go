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

type Pixel struct {
	R int
	G int
	B int
}

var (
	ColorBlack = Pixel{
		0,0,0,
	}

	ColorWhite = Pixel{
		255,255,255,
	}
)

func (p Pixel) String() string {
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



func (display *DisplayWriter) DrawRectangle(x1, y1, w, h int, color Pixel) error {

	sx1 := strconv.Itoa(x1)
	sy1 := strconv.Itoa(y1)
	sw := strconv.Itoa(w)
	sh := strconv.Itoa(h)
	sc := color.String()

	return display.Send("r", sx1, sy1, sw, sh, sc)

}

func (display *DisplayWriter) DrawPixel(x, y int, pixel Pixel) error {
	return display.DrawRectangle(x, y, 1, 1, pixel)
}

