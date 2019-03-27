package font

import (
	"github.com/cakesmith/webrender/app/component"
	"image"
	"image/color"
	"strconv"
	"strings"
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
		//cell.Draw()
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
	Background     color.Color
	XCells, YCells int
	cells          []*Cell
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
		XCells:     opts.XCells,
		YCells:     opts.YCells,
	}

	grid.Init = func() {

		grid.cells = []*Cell{}

		//numCells := g.XCells * g.YCells

		if opts.Center {
			grid.Center()
		}

		for y := grid.Min.Y; y < grid.Min.Y+opts.YCells*opts.CellHeight; y = y + opts.CellHeight {

			for x := grid.Min.X; x < grid.Min.X+opts.XCells*opts.CellWidth; x = x + opts.CellWidth {
				bounds := image.Rect(x, y, x+opts.CellWidth, y+opts.CellHeight)

				cell := NewCell(bounds, opts.BackgroundColor, opts.ActiveColor)
				grid.Container.Add(cell.Component)

				grid.cells = append(grid.cells, cell)

			}

		}

		grid.Draw()
		//grid.Container.Draw()

	}

	grid.Draw = func() {

		//grid.DrawRectangle(grid.Component.Rectangle, opts.LineColor)
		//
		//grid.DrawRectangle(grid.Component.Rectangle.Inset(opts.Thickness), opts.BackgroundColor)

		for x := grid.Min.X; x < grid.Max.X; x += opts.CellWidth {

			grid.DrawVert(x, grid.Min.Y, grid.Max.Y, opts.LineColor)
		}
		for y := grid.Min.Y; y < grid.Max.Y; y += opts.CellHeight {
			grid.DrawHoriz(grid.Min.X, grid.Max.X, y, opts.LineColor)
		}
	}

	return grid
}

func (grid *Grid) String() string {

	ret := []string{}

	for y := 0; y < grid.YCells; y++ {
		bits := ""
		for x := 0; x < len(grid.cells); x += grid.YCells {
			if grid.cells[x+y].active {
				bits += "1"
			} else {
				bits += "0"
			}
		}
		i, _ := strconv.ParseUint(bits, 2, 8)
		ret = append(ret, strconv.Itoa(int(i)))

	}
	return strings.Join(ret, ", ")
}
