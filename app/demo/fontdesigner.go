package demo

import (
	"github.com/cakesmith/webrender/app/font"
	"github.com/cakesmith/webrender/output"
)

func FontDesigner() *font.Grid {

	opts := font.GridOptions{
		CellWidth:       16,
		CellHeight:      22,
		XCells:          8,
		YCells:          11,
		Center: true,
		Thickness:       1,
		BackgroundColor: output.ColorRed,
		LineColor:       output.ColorBlack,
		ActiveColor:     output.ColorTerminalGreen,
	}

	grid := font.NewGrid(opts)

	grid.Draw = func() {
		//grid.DrawRectangle(grid.Rectangle, grid.Background)
	}

	return grid
}
