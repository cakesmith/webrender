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
		Center:          true,
		Thickness:       1,
		BackgroundColor: output.ColorRed,
		LineColor:       output.ColorWhite,
		ActiveColor:     output.ColorTerminalGreen,
	}

	grid := font.NewGrid(opts)

	return grid
}
