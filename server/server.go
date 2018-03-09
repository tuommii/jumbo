package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/gorilla/context"
	"github.com/tuommii/jumbo/database"
)

const (
	tmplDir  = "templates/"
	baseTmpl = "base.html"
	username = "lol"
	password = "lol"
)

// Server ...
type Server struct {
	db database.Database
	// Hold's all templates, filename is key
	templates map[string]*template.Template
	cookies   *sessions.CookieStore
}

// Create new server instance
func Create(db database.Database) *Server {
	return &Server{db: db, cookies: sessions.NewCookieStore([]byte("helloworld"))}
}

// Start server
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmap := template.FuncMap{
		"FormatPercentage": FormatPercentage,
		"FormatDate":       FormatDate,
	}

	// Cache templates
	s.initTemplates(baseTmpl, fmap)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api/create/match", auth(s.apiCreateMatch))
	http.HandleFunc("/api/create/game", auth(s.apiCreateGame))
	http.HandleFunc("/api/create/player", auth(s.apiCreatePlayer))

	http.HandleFunc("/api/delete/match", auth(s.apiDeleteMatch))
	http.HandleFunc("/api/delete/game", auth(s.apiDeleteGame))
	http.HandleFunc("/api/delete/player", auth(s.apiDeletePlayer))

	http.HandleFunc("/api/search", s.apiSearch)
	http.HandleFunc("/favicon.png", s.favicon)
	http.HandleFunc("/", s.apiHome)

	fmt.Println("Listening :8080")
	err := http.ListenAndServe("0.0.0.0:"+port, context.ClearHandler(http.DefaultServeMux))
	return err
}

// initTemplates inits all templates in folder
func (s *Server) initTemplates(base string, fmap template.FuncMap) {
	files, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		log.Fatal("Error while reading templates:", err)
	}

	s.templates = make(map[string]*template.Template)

	for _, f := range files {
		name := f.Name()
		if name != base {
			s.templates[name], err = template.New(name).Funcs(fmap).ParseFiles(
				tmplDir+base,
				tmplDir+name,
			)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

// getParams splits path and clears empty entries
func getParams(path string) []string {
	vars := strings.Split(path, "/")
	var params []string

	// Ignore empty strings
	for i := range vars {
		if vars[i] != "" {
			params = append(params, vars[i])
		}
	}
	return params
}

// FormatPercentage ...
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.0f", value*100)
}

// FormatDate ...
func FormatDate(date string) string {
	arr := strings.Split(date, "T")
	return arr[0]
}
