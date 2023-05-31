package unnecessary

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func render2(component *Component) {
	if !component.isEnabled {
		component.parent.node.RemoveChild(component.node)
		return
	}
	if component.populateItemCallback != nil {
		var itemsModel []Model = component.model.List()
		if len(itemsModel) == 0 {
			component.node.Parent.RemoveChild(component.node)
		} else {
			parentNode := component.node.Parent
			templateNode := CloneNodeTree(component.node, true)
			nodeNextSibling := component.node.NextSibling
			for itemIndex, itemModel := range itemsModel {
				var itemNode *html.Node
				if itemIndex == 0 {
					itemNode = component.node
				} else {
					itemNode = CloneNodeTree(templateNode, true)
					parentNode.InsertBefore(itemNode, nodeNextSibling)
				}
				itemId := fmt.Sprintf("__%s__%d", component.Id, itemIndex)
				NodeSetAttr(itemNode, NodeIdAttrName, itemId)

				itemComponent := NewComponent(itemId, nil)
				component.parent.Add(itemComponent)
				component.populateItemCallback(itemComponent, itemIndex, itemModel)
				render2(itemComponent)
			}
		}

	} else if component.model != nil {
		value := component.model.String()
		if component.node.DataAtom == atom.Input {
			NodeSetAttr(component.node, "value", value)
		} else {
			NodeSetText(component.node, value)
		}
	}

	if component.children != nil {
		for _, child := range component.children {
			render2(child)
		}
	}
}

func collectPageScript(component *Component) {
	if !component.isEnabled {
		return
	}
	if component.page != nil && len(component.behaviors) > 0 {
		script := ""
		for _, behavior := range component.behaviors {
			script += *behavior.GetScript(true)
		}
		script = fmt.Sprintf("(function(){%s;})();", script)

		scriptTag := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Script,
			Data:     "script",
			Attr: []html.Attribute{
				{Namespace: "", Key: "id", Val: fmt.Sprintf("s-%s", component.GetMarkupId())},
				{Namespace: "", Key: "type", Val: "text/javascript"},
			},
		}
		scriptTag.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: script,
		})
		head := FindNodeByTag(component.page.node, "head")
		head.AppendChild(&scriptTag)
	}

	for _, child := range component.children {
		collectPageScript(child)
	}
}

func collectAjaxScript(component *Component) string {
	if !component.isEnabled {
		return ""
	}
	script := ""
	if component.behaviors != nil {
		for _, behavior := range component.behaviors {
			script += *behavior.GetScript(false)
		}
	}
	for _, child := range component.children {
		script += collectAjaxScript(child)
	}
	return script
}

func RenderPage(page *Component) (string, error) {
	render2(page)
	collectPageScript(page)
	return RenderNode(page.node)
}

func RenderComponent(component *Component) (string, error) {
	render2(component)
	return RenderNode(component.node)
}
