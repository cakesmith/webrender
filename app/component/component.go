package component

import (
	"github.com/Sirupsen/logrus"
	"github.com/cakesmith/webrender/app/state"
	"github.com/cakesmith/webrender/output"
	"image"
	"image/color"
	"sync"
)

var (
	log = logrus.New()
)

type Drawable interface {
	Draw()
}

type Initializable interface {
	Init()
}

type Container struct {
	sync.Mutex
	image.Rectangle
	components []*Component
	Focused    []*Component
	state.Stateful
	output.Terminal
}

func (container *Container) Width() int {
	return container.Rectangle.Max.X - container.Rectangle.Min.X
}

func (container *Container) Height() int {
	return container.Rectangle.Max.Y - container.Rectangle.Min.Y
}

func (container *Container) Add(component *Component) {

	if container.components == nil {
		container.components = []*Component{}
	}

	container.components = append(container.components, component)
	component.Container = container

	if component.Init != nil {
		component.Init()
	}

}

func (container *Container) Draw() {
	for _, comp := range container.components {
		if comp.Draw != nil {
			comp.Draw()
		}
	}
}

func (container *Container) Init() {
	for _, comp := range container.components {
		if comp.Init != nil {
			comp.Init()
		}
	}
}

func (container *Container) OnKeypress(key int) {
	if container.Focused != nil && len(container.Focused) > 0 {
		for _, f := range container.Focused {
			if f.OnKeypress != nil {
				f.OnKeypress(key)
			}
		}
	} else {
		for _, comp := range container.components {
			if comp.OnKeypress != nil {
				comp.OnKeypress(key)
			}
		}
	}
}

func (container *Container) OnClick(btn, x, y int) {
	for _, comp := range container.components {

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
	Init       func()
	OnKeypress func(key int)
	OnClick    func(btn, x, y int)
	Draw       func()
}

func (comp *Component) Width() int {
	return comp.Rectangle.Max.X - comp.Rectangle.Min.X
}

func (comp *Component) Height() int {
	return comp.Rectangle.Max.Y - comp.Rectangle.Min.Y
}

func (comp *Component) Center() {

	w := comp.Width()
	h := comp.Height()

	gw := comp.Container.Width()
	gh := comp.Container.Height()

	x := int(gw/2) - int(w/2)
	y := int(gh/2) - int(h/2)

	comp.Rectangle = image.Rect(x, y, x+int(w), y+int(h))
}

type Border struct {
	color.Color
	Thickness int
}
