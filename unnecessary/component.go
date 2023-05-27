package unnecessary

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

const NodeIdAttrName = "unnecessary:id"

type PopulateItemCallback func(item *Component, index int, model Model)

type Component struct {
	Id                   string
	node                 *html.Node
	page                 *Component
	parent               *Component
	children             []*Component
	model                Model
	populateItemCallback PopulateItemCallback
	bejaviors            map[string]*Bejavior
	isEnabled            bool
	counter              int // page level property
}

func NewWicketPage(file string) (*Component, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	doc, err := html.Parse(data)
	if err != nil {
		return nil, err
	}

	page := Component{
		node:      doc,
		isEnabled: true,
	}
	return &page, nil
}

func NewComponent(id string, model Model) *Component {
	component := Component{
		Id:        id,
		model:     model,
		isEnabled: true,
	}
	return &component
}

func (component *Component) GetPage() *Component {
	if component.page == nil {
		return component
	} else {
		return component.page
	}
}

func (component *Component) Add(child *Component) *Component {
	node := FindNodeByAttr(component.node, NodeIdAttrName, child.Id)
	if node == nil {
		panic(fmt.Errorf("node with '%s' equal to '%s' not found", NodeIdAttrName, child.Id))
	}
	child.node = node
	child.page = component.GetPage()
	child.parent = component
	component.children = append(component.children, child)
	return child
}

func (component *Component) SetOutputMarkupId(output bool) {
	if output {
		value := fmt.Sprintf("id-%d", component.GetNextCounterValue())
		component.SetAttr("id", value)
	} else {
		component.DelAttr("id")
	}
}

func (component *Component) GetMarkupId() string {
	value, ok := component.GetAttr("id")
	if !ok {
		panic(fmt.Errorf("attribute \"id\" not found in murkup"))
	}
	return value
}

func (component *Component) SetModel(model Model) {
	component.model = model
}

func (component *Component) SetPopulateItemCallback(callback PopulateItemCallback) {
	component.populateItemCallback = callback
}

func (component *Component) GetAttr(key string) (string, bool) {
	return NodeGetAttr(component.node, key)
}

func (component *Component) SetAttr(key, value string) {
	NodeSetAttr(component.node, key, value)
}

func (component *Component) DelAttr(key string) {
	NodeDelAttr(component.node, key)
}

func (component *Component) AddBejavior(bejavior *Bejavior) {
	component.SetOutputMarkupId(true)
	page := component.GetPage()
	if page.bejaviors == nil {
		page.bejaviors = make(map[string]*Bejavior)
	}
	if component.bejaviors == nil {
		component.bejaviors = make(map[string]*Bejavior)
	}
	bejavior.Id = fmt.Sprintf("handler-%d", component.GetNextCounterValue())
	bejavior.component = component
	bejavior.page = page
	page.bejaviors[bejavior.Id] = bejavior
	component.bejaviors[bejavior.Id] = bejavior
}

func (component *Component) SetIsEnabled(isEnabled bool) {
	component.isEnabled = isEnabled
}

func (component *Component) GetNextCounterValue() int {
	if component.page == nil {
		value := component.counter
		component.counter++
		return value
	} else {
		return component.GetPage().GetNextCounterValue()
	}
}
