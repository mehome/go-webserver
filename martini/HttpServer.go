package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	//"github.com/bradhe/stopwatch"
	//"github.com/odysseus/stopwatch"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// Greeting: XML demo
type Greeting struct {
	XMLName xml.Name `xml:"greeting"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func getMilliSeconds(startTime time.Time, stopTime time.Time) float64 {
	// Returns the elapsed stopwatch time in milliseconds
	// time.Duration
	return float64(stopTime.Sub(startTime).Nanoseconds()) / 1000000
}

func getSeconds(startTime time.Time, stopTime time.Time) float64 {
	// Returns the elapsed stopwatch time in milliseconds
	// time.Duration
	return float64(stopTime.Sub(startTime).Nanoseconds()) / 1000000000
}

func Benchmark_Routing() {
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

	for i, r := range router.GetAllRoutes() {
		fmt.Printf("[%3d]: %-20s %-8s %-30s\n", i,
			r.GetName(), r.Method(), r.Pattern())
	}
	fmt.Printf("\n")

	for _, urls := range testUrls {
		matched := false
		for _, route := range router.GetAllRoutes() {
			match, params := route.Match(urls.method, urls.path)
			if match != martini.NoMatch {
				fmt.Printf("Matched: (%v, %-26v) - (%v), (%v)\n",
					urls.method, urls.path, match, params)
				matched = true
				break
			}
		}
		if !matched {
			fmt.Printf("Not matched: (%v, %-26v)\n", urls.method, urls.path)
		}
	}
	fmt.Printf("\n")

	kMaxIterators := 1000000

	startTotalTime := time.Now()
	for _, urls := range testUrls {
		matched := false
		startTime := time.Now()
		for i := 0; i < kMaxIterators; i++ {
			for _, route := range router.GetAllRoutes() {
				match, _ := route.Match(urls.method, urls.path)
				if match != martini.NoMatch {
					matched = true
					break
				}
			}
			if matched {
				userRouted++
			}
		}
		stopTime := time.Now()
		elapsedTime := getMilliSeconds(startTime, stopTime)

		fmt.Printf("Matched: (%s, %-26s) - (%10d), (%9.3f ms)\n",
			urls.method, urls.path, userRouted, elapsedTime)
	}
	fmt.Printf("\n")

	stopTotalTime := time.Now()
	totalElapsedTime := getSeconds(startTotalTime, stopTotalTime)

	fmt.Printf("userRouted = %d\n\n", userRouted)
	fmt.Printf("Total elapsed time: %9.3f second(s)\n", totalElapsedTime)
}

// Note: You can set the system environment variables
// "MARTINI_ENV" to "production", default value is "development".
// See: /martini/env.go

func api2(res http.ResponseWriter, req *http.Request) {
	html := "Hello world! /api"
	slice := make([]byte, len(html))
	copy(slice, html)
	res.WriteHeader(200)
	_, err := res.Write(slice[0:len(html)])
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// This will set the Content-Type header to "text/html; charset=UTF-8"
func home(r render.Render) {
	// template file: /views/home.tmpl
	r.HTML(200, "home", "jeremy")
}

// This will set the Content-Type header to "application/json; charset=UTF-8"
func api(r render.Render) {
	r.JSON(200, map[string]interface{}{"hello": "world"})
}

// This will set the Content-Type header to "text/xml; charset=UTF-8"
func xmlfunc(r render.Render) {
	r.XML(200, Greeting{One: "hello", Two: "world"})
}

// This will set the Content-Type header to "text/plain; charset=UTF-8"
func text(r render.Render) {
	r.Text(200, "This is a text format file. 'text/plain; charset=UTF-8'")
}

func getBooks(r render.Render, params martini.Params, req *http.Request) {
	r.Text(200, "/Books, action = get, id = "+params["id"])
}

func newBook(r render.Render, params martini.Params, req *http.Request) {
	r.Text(200, "/Books, action = new")
}

func updateBook(r render.Render, params martini.Params, req *http.Request) {
	r.Text(200, "/Books, action = update, id = "+params["id"])
}

func deleteBook(r render.Render, params martini.Params, req *http.Request) {
	r.Text(200, "/Books, action = delete, id = "+params["id"])
}

func main() {
	fmt.Printf("\n")
	fmt.Printf("Martini 1.0\n")
	fmt.Printf("\n")

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

	// This will set the Content-Type header to "text/html; charset=UTF-8"
	app.Get("/", home)
	app.Get("/home", home)
	app.Get("/index.htm", home)
	app.Get("/index.html", home)

	// This will set the Content-Type header to "application/json; charset=UTF-8"
	app.Get("/api", api)

	app.Get("/api/v1/:name", func(params martini.Params) string {
		return "API v1: name = " + params["name"]
	})

	app.Get("/api/v2/**", func(params martini.Params) string {
		return "API v2: ** = " + params["_1"]
	})

	app.Get("/api/v3/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
		return fmt.Sprintf("API v3 %s", params["name"])
	})

	// This will set the Content-Type header to "text/xml; charset=UTF-8"
	app.Get("/xml", xmlfunc)

	// This will set the Content-Type header to "text/plain; charset=UTF-8"
	app.Get("/text", text)

	app.Group("/books", func(r martini.Router) {
		r.Get("/:id", getBooks)
		r.Post("/new", newBook)
		r.Put("/update/:id", updateBook)
		r.Delete("/delete/:id", deleteBook)
	})

	Benchmark_Routing()

	//app.RunOnAddr(":8080")
}
