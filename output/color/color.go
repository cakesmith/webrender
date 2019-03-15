package color

import "fmt"

type Color struct {
	R uint8
	G uint8
	B uint8
}

var (
	Background = Color{
		40, 40, 40,
	}

	TerminalGreen = Color{
		51, 255, 51,
	}

	Black = Color{
		0, 0, 0,
	}

	White = Color{
		255, 255, 255,
	}

	Red = Color{
		200, 0, 0,
	}
)

func (c Color) String() string {
	return fmt.Sprintf("%v-%v-%v", c.R, c.G, c.B)
}
