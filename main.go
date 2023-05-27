package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"unnecessary/unnecessary"
)

func page2Serve() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.NoCache)
	router.Use(middleware.Heartbeat("/ping"))
	// router.Mount("/debug", middleware.Profiler())

	fs := http.FileServer(http.Dir("html/static"))
	router.Mount("/static", http.StripPrefix("/static/", fs))

	var page2Count uint64 = 0
	var page2 = func(w http.ResponseWriter, r *http.Request) {
		page, err := unnecessary.NewWicketPage("html/page2.html")
		if err != nil {
			panic(err)
		}

		holder := unnecessary.NewComponent("holder", nil)
		page.Add(holder)
		holder.SetOutputMarkupId(true)

		rawValue := unnecessary.NewComponent("rawValue", unnecessary.NewGenericModel(page2Count))
		holder.Add(rawValue)

		dynamicModel := unnecessary.NewDynamicModel(func() interface{} { return page2Count })
		holder.Add(unnecessary.NewComponent("dynamicValue", dynamicModel))
		holder.Add(unnecessary.NewComponent("textValue", dynamicModel))
		holder.Add(unnecessary.NewComponent("textareaValue", dynamicModel))

		btnPlus1 := unnecessary.NewComponent("btnPlus1", nil)
		holder.Add(btnPlus1)
		// wicket.SetComponentAttribute(btnPlus1, "id", "btnPlus1")
		btnPlus1.SetOutputMarkupId(true)
		// btnPlus1.AddAjaxHandler("click", func(target *unnecessary.AjaxTarget) error {
		// 	page2Count += 1
		// 	rawValue.SetModel(unnecessary.NewGenericModel(page2Count))
		// 	target.Add(holder)
		// 	return nil
		// })
		btnPlus1.AddBejavior(unnecessary.NewAjaxBejavior("click", func(target *unnecessary.AjaxTarget) error {
			page2Count += 1
			rawValue.SetModel(unnecessary.NewGenericModel(page2Count))
			target.Add(holder)
			return nil
		}))

		btnMenos1 := unnecessary.NewComponent("btnMenos1", nil)
		holder.Add(btnMenos1)
		// wicket.SetComponentAttribute(btnMenos1, "id", "btnMenos1")
		btnMenos1.SetOutputMarkupId(true)
		// btnMenos1.AddAjaxHandler("click", func(target *unnecessary.AjaxTarget) error {
		// 	page2Count -= 1
		// 	rawValue.SetModel(unnecessary.NewGenericModel(page2Count))
		// 	target.Add(holder)
		// 	return nil
		// })
		btnMenos1.AddBejavior(unnecessary.NewAjaxBejavior("click", func(target *unnecessary.AjaxTarget) error {
			page2Count -= 1
			rawValue.SetModel(unnecessary.NewGenericModel(page2Count))
			target.Add(holder)
			return nil
		}))

		// loop
		// loop := wicket.NewComponent(holder, "loop", loopModel)
		loop := unnecessary.NewComponent("loop", unnecessary.NewGenericListModel(1.1, 2, 3))
		holder.Add(loop)
		loop.SetPopulateItemCallback(func(loopItem *unnecessary.Component, index int, model unnecessary.Model) {
			v := model.(*unnecessary.GenericModel)
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
			loopItem.Add(unnecessary.NewComponent("loopValue", unnecessary.NewGenericModel(v3)))
		})

		// table
		// tr := wicket.NewComponent("tr", wicket.NewGenericListModel("1", "2"))
		// holder.Add(tr)
		// tr.SetListPopulateItem(func(trItem *wicket.WicketComponent, trIndex int, trModel wicket.Model) {
		// 	trModelStr := trModel.String()
		// 	tdsModel := wicket.NewGenericListModel(fmt.Sprintf("%v:1", trModelStr), fmt.Sprintf("%v:2", trModelStr))
		// 	td := wicket.NewComponent("td", tdsModel)
		// 	trItem.Add(td)
		// 	td.SetListPopulateItem(func(tdItem *wicket.WicketComponent, tdIndex int, tdModel wicket.Model) {
		// 		tdModelStr := fmt.Sprintf("%v:%d", tdModel, page2Count)
		// 		tdItem.SetModel(wicket.NewGenericModel(tdModelStr))
		// 	})
		// })
		tr := unnecessary.NewComponent("tr", unnecessary.NewGenericListModel(
			unnecessary.NewGenericListModel("11", "12"),
			unnecessary.NewGenericListModel("21", "22"),
		))
		holder.Add(tr)
		tr.SetPopulateItemCallback(func(trItem *unnecessary.Component, trIndex int, trModel unnecessary.Model) {
			// td := wicket.NewComponent("td", trModel)
			td := trItem.Add(unnecessary.NewComponent("td", trModel))
			td.SetPopulateItemCallback(func(tdItem *unnecessary.Component, tdIndex int, tdModel unnecessary.Model) {
				tdModelStr := fmt.Sprintf("%s:%d", tdModel.String(), page2Count)
				tdItem.SetModel(unnecessary.NewGenericModel(tdModelStr))
			})
		})

		// notEnabled
		notEnabled := unnecessary.NewComponent("notEnabled", nil)
		notEnabled.SetIsEnabled(false)
		holder.Add(notEnabled)

		if r.Method == "GET" {
			pageStr, err := unnecessary.RenderPage(page)
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
			responseData, err := unnecessary.ProcessAjaxRequest(page, requestData)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Add("Content-Type", "application/json")
			w.Write(responseData)
		}
	}
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/page2", http.StatusFound)
	})
	router.Get("/page2", page2)
	router.Post("/page2", page2)

	log.Print("ready...")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Panicf("Terminated - %v", err)
	} else {
		log.Panicf("Terminated")
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	page2Serve()
}
