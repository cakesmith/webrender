package component

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app/container"
	"github.com/cakesmith/webrender/output"
	"image"
	"image/color"
)

var (
	log = logrus.New()
)

type Drawable interface {
	Draw()
}

type Composable interface {
	GetComponents() []Component
}

type Initializable interface {
	Init()
}

type Focusable interface {
	OnKeypress(key int)
}

type Clickable interface {
	OnClick(btn, x, y int)
}


type Container struct {
	image.Rectangle
	Components []*Component
	Focused Focusable
	state.Stateful
	output.Terminal
}

func (container *Container) Add(component *Component) {
	if container.Components == nil {
		container.Components = []*Component{}
	}
	container.Components = append(container.Components, component)
}

func (container *Container) Init() {
	for _, comp := range container.Components {
		comp.Container = container
		comp.Init()
	}
}

func (container *Container) OnKeypress(key int) {
	if container.Focused != nil {
		container.Focused.OnKeypress(key)
	} else {
		for _, comp := range container.Components {
			if comp.OnKeypress != nil {
				comp.OnKeypress(key)
			}
		}
	}
}

func (container *Container) OnClick(btn, x, y int) {
	for _, comp := range container.Components {

		bounds := comp.Bounds()

		if (image.Point{X: x, Y: y}).In(bounds) {

			dx := x - bounds.Min.X
			dy := y - bounds.Min.Y

			if comp.OnClick != nil {
				comp.OnClick(btn, dx, dy)
			}

		}
	}
}

type Component struct {
	image.Rectangle
	color.Color
	*Container
	Init func()
	OnKeypress func(key int)
	OnClick func(btn, x, y int)
	Draw func()
}

func (component *Component) Set(x, y int, color color.Color) {
	if (image.Point{X:x, Y:y}).In(component.Bounds()) {
		dx := x + component.Bounds().Min.X
		dy := y + component.Bounds().Min.Y
		component.Container.Set(dx, dy, color)
	}
}

func (component *Component) DrawRectangle(x1, y1, w, h int, c color.Color) {
	component.Container.DrawRectangle(x1, y1, w, h, c)
}

type Border struct {
	color.Color
	Thickness int
}


