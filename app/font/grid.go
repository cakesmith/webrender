package font

import (
	"fmt"
	"github.com/cakesmith/webrender/app/component"
	"image"
	"image/color"
)

type Cell struct {
	*component.Component
	Background color.Color
	active     bool
}

func NewCell(bounds image.Rectangle, background, c color.Color) *Cell {

	cell := &Cell{
		Component:  &component.Component{
			Rectangle:  bounds,
			Color:      c,
			Init:       nil,
			OnKeypress: nil,
		},
		Background: background,
		active:     false,
	}

	cell.Init = func() {
		cell.Draw()
	}

	cell.Draw = func() {
		var c color.Color
		if cell.active {
			c = cell.Component.Color
		} else {
			c = cell.Background
		}
		cell.DrawRectangle(cell.Component.Rectangle, c)
	}

	cell.OnClick = func(btn, x, y int) {
		cell.active = !cell.active
		cell.Draw()
	}

	return cell
}

type Grid struct {
	*component.Component
	Background             color.Color
	CharWidth, CharHeight int
}

type GridOptions struct {

	CellWidth,
	CellHeight,
	XCells,
	YCells,
	X,
	Y,
	Thickness int

	BackgroundColor,
	LineColor,
	ActiveColor color.Color
}

//NewGrid creates a new grid
//xCells = number of cells wide
//yCells = number of cells tall
//(x, y) is the top left point of the grid
//thickness is the thickness of the grid lines
func NewGrid(g GridOptions) *Grid {

	width := g.CellWidth * g.XCells
	height := g.CellHeight * g.YCells

	grid := &Grid{
		Component:  &component.Component{
			Rectangle:  image.Rect(g.X, g.Y, g.X + width, g.Y + height),
			Color:      g.LineColor,
		},
		Background: g.BackgroundColor,
		CharWidth:  0,
		CharHeight: 0,
	}

	grid.Init = func() {

		numCells := g.XCells * g.YCells

		for y := 0; y < g.YCells; y++ {
			for x := 0; x < g.XCells; x++ {

				cx := x * g.CellWidth
				cy := y * g.CellHeight

				bounds := image.Rect(cx, cy, cx +g.CellWidth, cy +g.CellHeight)
				cell := NewCell(bounds, g.BackgroundColor, g.ActiveColor)
				grid.Container.Add(cell.Component)

			}
		}

		fmt.Println(numCells)

		//grid.Draw()

	}

	return grid
}

func (grid *Grid) String() string {
	return "not implemented"
}
