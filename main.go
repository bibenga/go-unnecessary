package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	u "unnecessary/unnecessary"
)

var page2Count uint64 = 0

func page2(w http.ResponseWriter, r *http.Request) {
	page, err := u.NewWicketPage("html/page2.html")
	if err != nil {
		panic(err)
	}

	holder := u.NewComponent("holder", nil)
	page.Add(holder)
	holder.SetOutputMarkupId(true)

	rawValue := u.NewComponent("rawValue", u.NewGenericModel(page2Count))
	holder.Add(rawValue)

	dynamicModel := u.NewDynamicModel(func() interface{} { return page2Count })
	holder.Add(u.NewComponent("dynamicValue", dynamicModel))
	holder.Add(u.NewComponent("textValue", dynamicModel))
	holder.Add(u.NewComponent("textareaValue", dynamicModel))

	btnPlus1 := u.NewComponent("btnPlus1", nil)
	holder.Add(btnPlus1)
	btnPlus1.SetOutputMarkupId(true)
	btnPlus1.AddBehavior(u.NewAjaxBehavior("click", func(target *u.AjaxTarget) error {
		page2Count += 1
		rawValue.SetModel(u.NewGenericModel(page2Count))
		target.Add(holder)
		return nil
	}))

	btnMenos1 := u.NewComponent("btnMenos1", nil)
	holder.Add(btnMenos1)
	btnMenos1.SetOutputMarkupId(true)
	btnMenos1.AddBehavior(u.NewAjaxBehavior("click", func(target *u.AjaxTarget) error {
		page2Count -= 1
		rawValue.SetModel(u.NewGenericModel(page2Count))
		target.Add(holder)
		return nil
	}))

	// simple loop
	// loop := wicket.NewComponent(holder, "loop", loopModel)
	loop := u.NewComponent("loop", u.NewGenericListModel(1.1, 2, 3))
	holder.Add(loop)
	loop.SetPopulateItemCallback(func(loopItem *u.Component, index int, model u.Model) {
		v := model.(*u.GenericModel)
		v2 := v.Value
		var v3 float64
		switch v := v2.(type) {
		case int:
			v3 = math.Pow(float64(page2Count), float64(v))
		case float64:
			v3 = math.Pow(float64(page2Count), float64(v))
		default:
			panic("Olala")
		}
		loopItem.Add(u.NewComponent("loopValue", u.NewGenericModel(v3)))
	})

	// table
	tr := u.NewComponent("tr", u.NewGenericListModel(
		u.NewGenericListModel("11", "12"),
		u.NewGenericListModel("21", "22"),
	))
	holder.Add(tr)
	tr.SetPopulateItemCallback(func(trItem *u.Component, trIndex int, trModel u.Model) {
		td := trItem.Add(u.NewComponent("td", trModel))
		td.SetPopulateItemCallback(func(tdItem *u.Component, tdIndex int, tdModel u.Model) {
			tdModelStr := fmt.Sprintf("%s:%d", tdModel.String(), page2Count)
			tdItem.SetModel(u.NewGenericModel(tdModelStr))
		})
	})

	// not enabled component
	notEnabled := u.NewComponent("notEnabled", nil)
	notEnabled.SetIsEnabled(false)
	holder.Add(notEnabled)

	if r.Method == "GET" {
		pageStr, err := u.RenderPage(page)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(pageStr))
	} else {
		requestData, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		responseData, err := u.ProcessAjaxRequest(page, requestData)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write(responseData)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.NoCache)
	router.Use(middleware.Heartbeat("/ping"))
	router.Mount("/debug", middleware.Profiler())

	fs := http.FileServer(http.Dir("html/static"))
	router.Mount("/static", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/page2", http.StatusFound)
	})

	router.Get("/page2", page2)
	router.Post("/page2", page2)

	log.Print("ready at 8000 port: http://127.0.0.1:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Panicf("Terminated - %v", err)
	} else {
		log.Panicf("Terminated")
	}
}
