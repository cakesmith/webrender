package font

import (
	"github.com/cakesmith/webrender/app"
	"github.com/cakesmith/webrender/app/component"
)

type Designer struct {
	app.App
	components []component.Component
}

func (d *Designer) GetComponents() []component.Component {
	return d.components
}

func (d *Designer) Init() {

	resetBtn := &component.Button{

	}


	d.components = []component.Component{
		resetBtn.Component,
	}
}

