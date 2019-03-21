package font

import (
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
		Component: &component.Component{
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
	Background            color.Color
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

	Center bool

	BackgroundColor,
	LineColor,
	ActiveColor color.Color
}

//NewGrid creates a new grid
//xCells = number of cells wide
//yCells = number of cells tall
//(x, y) is the top left point of the grid
//thickness is the thickness of the grid lines
func NewGrid(opts GridOptions) *Grid {

	width := opts.CellWidth * opts.XCells
	height := opts.CellHeight * opts.YCells

	grid := &Grid{
		Component: &component.Component{
			Rectangle: image.Rect(opts.X, opts.Y, opts.X+width, opts.Y+height),
			Color:     opts.LineColor,
		},
		Background: opts.BackgroundColor,
		CharWidth:  0,
		CharHeight: 0,
	}

	grid.Init = func() {

		//numCells := g.XCells * g.YCells

		if opts.Center {
			grid.Center()
		}

		for y := grid.Min.Y; y < grid.Min.Y+opts.YCells*opts.CellHeight; y = y + opts.CellHeight {
			for x := grid.Min.X; x < grid.Min.X+opts.XCells*opts.CellWidth; x = x + opts.CellWidth {

				bounds := image.Rect(x, y, x+opts.CellWidth, y+opts.CellHeight)
				cell := NewCell(bounds, opts.BackgroundColor, opts.ActiveColor)
				grid.Container.Add(cell.Component)

			}
		}

		grid.Container.Draw()

	}

	return grid
}

func (grid *Grid) String() string {
	return "not implemented"
}
