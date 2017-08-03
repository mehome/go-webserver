package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// Greeting: XML demo
type Greeting struct {
	XMLName xml.Name `xml:"greeting"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

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
	r.Text(200, "hello, world")
}

func getBooks(r render.Render) {

}

func newBook(r render.Render) {

}

func updateBook(r render.Render) {

}

func deleteBook(r render.Render) {

}

func main() {
	fmt.Printf("Martini 1.0\n")
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
		IndexFile: "index.html", // Web defualt doucment filename
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

	app.RunOnAddr(":8080")
}
