package os

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
	"sync"
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
	sync.Mutex
	memo map[string]bool
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

func (display *DisplayWriter) DrawPixel(x, y int, pixel Pixel) error {

	log.WithFields(logrus.Fields{"x": x, "y": y, "pixel": pixel}).Debug("drawing pixel")

	sx := strconv.Itoa(x)
	sy := strconv.Itoa(y)
	sp := pixel.String()

	key := strings.Join([]string{sx, sy, sp}, ",")

	display.Lock()

	if display.memo == nil {
		display.memo = make(map[string]bool)
	}

	m := display.memo[key]

	if !m {
		display.memo[key] = true
		display.Unlock()
		return display.Send("d", sx, sy, sp)
	}

	display.Unlock()
	return nil

}

