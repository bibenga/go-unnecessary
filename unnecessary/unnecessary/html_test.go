package unnecessary

import (
	"strings"
	"testing"
)

func TestParseHtml(t *testing.T) {
	htmlData := "<html><head></head><body>olala</body></html>"
	htmlReader := strings.NewReader(htmlData)
	html, err := ParseHtml(htmlReader)
	if err != nil {
		t.Error(err)
	}
	if html == nil {
		t.Errorf("html is nil")
	}
}

func TestParseHtmlAndRenderNode(t *testing.T) {
	htmlData := "<html><head></head><body>olala</body></html>"
	htmlReader := strings.NewReader(htmlData)
	html, err := ParseHtml(htmlReader)
	if err != nil {
		t.Error(err)
	}
	if html == nil {
		t.Errorf("html is nil")
	}
	htmlRendered, err := RenderNode(html)
	if err != nil {
		t.Error(err)
	}
	if htmlData != htmlRendered {
		t.Errorf("source and destination htmls are mismatch")
	}
}
