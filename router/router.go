package router

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Parameter struct {
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Mandatory bool        `json:"mandatory"`
	Default   interface{} `json:"default,omitempty"`
	ParamType string
}

type Parameters struct {
	Url     []Parameter `json:"url,omitempty"`
	Headers []Parameter `json:"headers,omitempty"`
	Query   []Parameter `json:"query,omitempty"`
}

type Action struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Route struct {
	Method     string     `json:"method"`
	Path       string     `json:"path"`
	Parameters Parameters `json:"parameters"`
	Action     Action     `json:"action"`
	Exceptions []string   `json:"exceptions,omitempty"`
}

type Routes struct {
	Routes []Route `json:"routes"`
}

type HttpRouter struct {
	Routes []Route
}

func (r *HttpRouter )LoadRoute(path string, routeFileName string) {
	var files []string
	err := filepath.Walk(path, Lookup(&files, routeFileName))

	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(filepath.Ext(path))
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("File reading error", err)
		}
		tempRouter := new(Routes)
		if err := json.Unmarshal(data, &tempRouter); err != nil {
			panic(err)
		}
		r.Routes =  append(r.Routes, tempRouter.Routes...)
	}
}

func (r *HttpRouter ) BindChiRouter() http.Handler {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	chiRouter.MethodFunc("Get", "/link", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("custom link method"))
	})
	for _, route := range r.Routes {
		fmt.Println(route.Path)
		chiRouter.Group(func(r chi.Router) {
			// Add middleware
			r.Use(middleware.RequestID)
			r.Use(middleware.Logger)
			r.Use(middleware.Recoverer)
			r.Use(middleware.URLFormat)
			r.MethodFunc(route.Method, route.Path, func(w http.ResponseWriter, r *http.Request){
				err, vars := getVariables(r, route)
				if err != nil {
					fmt.Println("error", err)
				}
				fmt.Println(vars)
				fmt.Println(route.Parameters)
				fmt.Fprintf(w, "%s %q", r.Method, html.EscapeString(r.URL.Path))
			})
		})
	}

	return chiRouter
}
