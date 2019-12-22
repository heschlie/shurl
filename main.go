package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"shurl/models"
	shurl "shurl/pkg"
	"strings"
)

var urls = make(map[string]models.Shurl)

func main() {
	m := &mux.Router{}
	m.HandleFunc("/", home)
	m.HandleFunc("/{hash}", shortened)

	http.ListenAndServe(":8080", m)
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("static/form.html"))
    if r.Method != http.MethodPost {
    	err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, "There was an issue", 500)
		}
		return
	}
	url := r.FormValue("url")
	// Only use returned lower value for checking http to maintain case in saved URL
	if !strings.HasPrefix(strings.ToLower(url), "http") {
		url = "http://" + url
	}

	s := models.Shurl{
		Hash: shurl.RandStringRunes(8),
		Url: url,
		Hits: 0,
	}
	fmt.Printf("%+v", s)

	urls[s.Hash] = s
	fmt.Fprint(w, s.Hash)
}

func shortened(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	s, ok := urls[hash]; if ok {
		s.Hits += 1
		http.Redirect(w, r, s.Url, 301)
	} else {
		t := template.Must(template.ParseFiles("static/form.html"))
		err := t.Execute(w, struct {NotFound bool}{true})
		if err != nil {
			http.Error(w, "Something broke", 500)
		}
	}

}
