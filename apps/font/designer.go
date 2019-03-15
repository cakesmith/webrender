package font

import (
	"github.com/cakesmith/webrender/apps"
	"github.com/cakesmith/webrender/output/color"
	"math"
	"strconv"
	"strings"
)

type Grid struct {
	apps.Component
	Background            color.Color
	CharWidth, CharHeight uint
	bitmap                map[uint]bool
}

type bounds struct {
	left, right, bottom, top uint
}

func (g *Grid) reset() {
	g.bitmap = make(map[uint]bool)
	g.Draw()
}

func (g *Grid) String() string {

	var n []string

	for y := uint(0); y < g.CharHeight; y++ {

		var ch byte

		for x := uint(0); x < g.CharWidth; x++ {

			i := x + g.CharWidth*y

			if g.bitmap[i] {
				ch = ch | (1 << x)
			}

		}

		n = append(n, strconv.Itoa(int(ch)))
	}

	return strings.Join(n, ", ")
}

func (g *Grid) bounds(n int) bounds {

	x := uint(math.Mod(float64(n), float64(g.CharWidth)))
	y := uint(math.Floor(float64(n) / float64(g.CharWidth)))

	return g.calcBounds(x, y)

}

func (g *Grid) calcBounds(x, y uint) bounds {
	left := g.X + (g.CharWidth * x)
	right := left + g.CharWidth
	bottom := g.Y + (g.CharHeight * y)
	top := bottom + g.CharHeight

	return bounds{left, right, bottom, top}
}

func (g *Grid) which(mx, my uint) uint {

	for y := uint(0); y < g.CharHeight; y++ {
		for x := uint(0); x < g.CharWidth; x++ {

			b := g.calcBounds(x, y)

			if mx > b.left && mx < b.right && my > b.bottom && my < b.top {
				return x + (g.CharWidth * y)
			}
		}
	}
	return -1
}

func (g *Grid) Draw() error {

	//TODO quick and dirty, use g.CharWidth and g.CharHeight instead of 8 and 11

	x1 := 28 * g.CharWidth
	y1 := 10 * g.CharHeight
	w1 := 36 * g.CharWidth
	h1 := 21*g.CharHeight + 1

	g.DrawRectangle(28*g.CharWidth, 10*g.CharHeight, 8*g.CharWidth, 11*g.CharHeight, g.Background)

	for x := g.CharWidth; x < g.Width; x = x + g.CharWidth {
		if x < 28*8 || x > 36*8 {
			continue
		}

		err := g.DrawVert(x, y1, h1, g.Color)

		if err != nil {
			return err
		}
	}

	for y := g.CharHeight; y < g.Height; y = y + g.CharHeight {
		if y < 10*11 || y > 21*11 {
			continue
		}

		err := g.DrawHoriz(x1, w1, y, g.Color)

		if err != nil {
			return err
		}
	}

	return nil
}
