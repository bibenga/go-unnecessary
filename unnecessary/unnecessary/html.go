package unnecessary

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/net/html"
)

type NodeScanTest func(*html.Node) bool

func ScanNodeTree(root *html.Node, cb NodeScanTest) {
	var stop = false
	var scanner func(*html.Node)
	scanner = func(node *html.Node) {
		stop = cb(node)
		if stop {
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			scanner(child)
			if stop {
				return
			}
		}
	}
	scanner(root)
}

func FindNode(root *html.Node, f NodeScanTest) *html.Node {
	var node *html.Node
	ScanNodeTree(root, func(n *html.Node) bool {
		if f(n) {
			node = n
		}
		return node != nil
	})
	return node
}

func FindNodeByAttr(root *html.Node, attr, value string) *html.Node {
	return FindNode(root, func(n *html.Node) bool {
		for _, a := range n.Attr {
			if a.Key == attr && a.Val == value {
				return true
			}
		}
		return false
	})
}

func FindNodeByTag(root *html.Node, tag string) *html.Node {
	node := FindNode(root, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == tag {
			return true
		}
		return false
	})
	return node
}

func FindNodesByTag(root *html.Node, tag string) []*html.Node {
	var res []*html.Node
	ScanNodeTree(root, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == tag {
			res = append(res, n)
		}
		return false
	})
	return res
}

func NodeGetAttr(n *html.Node, key string) (string, bool) {
	for i := 0; i < len(n.Attr); i++ {
		attr := &n.Attr[i]
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func NodeSetAttr(n *html.Node, key, value string) {
	wasFound := false
	for i := 0; i < len(n.Attr); i++ {
		attr := &n.Attr[i]
		if attr.Key == key {
			wasFound = true
			attr.Val = value
			break
		}
	}
	if !wasFound {
		n.Attr = append(n.Attr, html.Attribute{Namespace: "", Key: key, Val: value})
	}
}

func NodeDelAttr(n *html.Node, key string) {
	for i := 0; i < len(n.Attr); i++ {
		attr := &n.Attr[i]
		if attr.Key == key {
			n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			break
		}
	}
}

func NodeSetText(n *html.Node, text string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		n.RemoveChild(c)
	}
	textNode := html.Node{
		Type: html.TextNode,
		Data: text,
	}
	n.AppendChild(&textNode)
}

func CloneNodeTree(src *html.Node, includeChildren bool) *html.Node {
	dst := &html.Node{
		Type:     src.Type,
		DataAtom: src.DataAtom,
		Data:     src.Data,
		Attr:     make([]html.Attribute, len(src.Attr)),
	}
	copy(dst.Attr, src.Attr)

	if includeChildren {
		for srcChild := src.FirstChild; srcChild != nil; srcChild = srcChild.NextSibling {
			dstChild := CloneNodeTree(srcChild, includeChildren)
			dst.AppendChild(dstChild)
		}
	}

	return dst
}

func ParseHtmlFile(file string) (*html.Node, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return ParseHtml(data)
}

func ParseHtml(r io.Reader) (*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func RenderNode(node *html.Node) (string, error) {
	buf := new(bytes.Buffer)
	err := html.Render(buf, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
