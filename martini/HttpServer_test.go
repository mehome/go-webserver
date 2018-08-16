package main

import (
	//"fmt"
	//"net/http"

	//"net/http/httptest"
	//"strings"
	"testing"

	martini "github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
)

func Test_Routing(t *testing.T) {
	//fmt.Printf("Martini 1.0\n")
	/*
		app := martini.Classic()

		// render html templates from templates directory
		app.Use(render.Renderer(render.Options{
			Directory:  "views",                    // Specify what path to load the templates from.
			Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
			Charset:    "UTF-8",                    // Sets encoding for json and html content-types. Default is "UTF-8".
			IndentJSON: true,                       // Output human readable JSON
			IndentXML:  true,                       // Output human readable XML
		}))

		app.Use(martini.Static("static"))

		static := martini.Static("static", martini.StaticOptions{
			IndexFile: "index.html", // Web default doucment filename
			Fallback:  "/404.html",  // Not found page filename
			Exclude:   "/api/v",     // Exclude path
		})
		app.NotFound(static, http.NotFound)
	*/

	userRouted := 0
	router := martini.NewRouter()

	router.Get("/service/candy/:kind", func(params martini.Params) {
		userRouted++
	})

	router.Get("/service/shutdown", func() {
		userRouted++
	})

	router.Get("/", func() {
		userRouted++
	})

	router.Get("/:filename", func(params martini.Params) {
		userRouted++
	})

	// Run benchmark of various urls
	var testUrls = []struct {
		method string
		path   string
		match  martini.RouteMatch
	}{
		{"GET", "/service/candy/lollipop", martini.ExactMatch},
		{"GET", "/service/candy/gum", martini.ExactMatch},
		{"GET", "/service/candy/seg_ratta", martini.ExactMatch},
		{"GET", "/service/candy/lakrits", martini.ExactMatch},

		{"GET", "/service/shutdown", martini.ExactMatch},
		{"GET", "/", martini.ExactMatch},
		{"GET", "/some_file.html", martini.ExactMatch},
		{"GET", "/another_file.jpeg", martini.ExactMatch},
	}

	for _, urls := range testUrls {
		matched := false
		for _, route := range router.GetAllRoutes() {
			match, params := route.Match(urls.method, urls.path)
			if match != martini.NoMatch {
				//t.Errorf("expected: (%v) got: (%v, %v)", urls.match, match, params)
				t.Logf("Matched: (%v, %-26v) - (%v), (%v)\n", urls.method, urls.path, match, params)
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("Not matched: (%v), (%v)\n", urls.method, urls.path)
		}
	}

	t.Logf("userRouted = %v\n", userRouted)
}
