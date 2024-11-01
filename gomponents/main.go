package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	// b "github.com/willoma/bulma-gomponents"
	// e "github.com/willoma/gomplements"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	hxhttp "maragu.dev/gomponents-htmx/http"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

func main() {
	if err := start(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func start() error {
	now := time.Now()
	mux := http.NewServeMux()
	mux.HandleFunc("/", ghttp.Adapt(
		func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			if r.Method == http.MethodPost && hxhttp.IsBoosted(r.Header) {
				now = time.Now()

				hxhttp.SetPushURL(w.Header(), "/?time="+now.Format(timeFormat))

				return partial(now), nil
			}
			return page(now), nil
		},
	))

	mux.HandleFunc("/2", ghttp.Adapt(
		func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return page2(), nil
		},
	))

	log.Println("Starting on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

const timeFormat = "15:04:05"

func page(now time.Time) g.Node {
	return c.HTML5(
		c.HTML5Props{
			Title: now.Format(timeFormat),
			Head: []g.Node{
				h.Script(h.Src("https://unpkg.com/htmx.org")),
			},
			Body: []g.Node{
				b.Container(
					b.Notification(
						b.Primary,
						"This container is ", e.Strong("centered"), " on desktop and larger viewports.",
					),
				),
				h.Div(
					h.Class("max-w-7xl mx-auto p-4 prose lg:prose-lg xl:prose-xl"),
					h.H1(g.Text(`gomponents + HTMX`)),
					h.P(g.Textf(`Time at last full page refresh was %v.`, now.Format(timeFormat))),
					partial(now),
					h.FormEl(
						h.Method("post"),
						h.Action("/"),
						hx.Boost("true"), hx.Target("#partial"), hx.Swap("outerHTML"),
						h.Button(h.Type("submit"), g.Text(`Update time`),
							h.Class("rounded-md border border-transparent bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"),
						),
					),
				),
			},
		},
	)
}

func partial(now time.Time) g.Node {
	return h.P(
		h.ID("partial"),
		g.Textf(`Time was last updated at %v.`, now.Format(timeFormat)),
	)
}

func Navbar() g.Node {
	return b.Navbar(
		b.Dark,
		b.NavbarStart(
			b.NavbarAHref("/", "Home"),
			b.NavbarAHref("/about", "About"),
		),
	)
}

func page2() g.Node {
	return b.HTML(
		b.Script("https://unpkg.com/htmx.org"),
		b.CSSPath("https://cdn.jsdelivr.net/npm/bulma@1.0.0/css/bulma.min.css"),
		b.HTitle("Olala"),
		b.Language("en"),
		b.Head(e.Meta(e.Charset("utf-8"))),
		b.Container(
			b.Navbar(
				b.NavbarBrand(),
				b.NavbarAHref("/", "Home"),
				b.NavbarAHref("/2", "22"),
			),
			b.Content(
				e.H1("Hello"),
				e.P("Hello world"),
			),
		),
	)
}
