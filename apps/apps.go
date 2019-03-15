package apps

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/output"
	"github.com/cakesmith/webrender/output/color"
)

var (
	log = logrus.New()
)

type Component struct {
	X, Y, Width, Height uint
	color.Color
	Drawable
}

type Drawable interface {
	output.Drawer
	Hit(mx, my int) bool
	Draw()
}
