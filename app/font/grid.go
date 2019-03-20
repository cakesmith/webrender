package font

import (
	"github.com/cakesmith/webrender/app/component"
	"image/color"
	"strconv"
	"strings"
)

type Cell struct {
	component.Component
	Background color.RGBA
	active     bool
}

type Grid struct {
	component.Component
	Background            color.RGBA
	CharWidth, CharHeight int
	bitmap                map[int]map[int]*Cell
}

func (g *Grid) Init() {
	g.bitmap = make(map[int]map[int]*Cell)
	g.Draw()
}

func (g *Grid) String() string {

	var n []string

	for y := 0; y < g.CharHeight; y++ {

		var ch byte

		for x := 0; x < g.CharWidth; x++ {

			if g.bitmap[y][x].active {
				ch = ch | (1 << uint(x))
			}

		}

		n = append(n, strconv.Itoa(int(ch)))
	}

	return strings.Join(n, ", ")
}
