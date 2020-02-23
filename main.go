package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"shurl/internal/db"
	shurl "shurl/pkg"
	"strings"
)

var dq = &db.Queries{}

func main() {
	host := "localhost"
	port := 5432
	user := "shurl"
	pass := "micron"
	dbname := "shurl"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							host, port, user, pass, dbname)
	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.WithError(err).Fatalf("Failed to connect to DB %s", host)
	}
	dq = db.New(database)

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

	ctx := context.Background()

	params := db.CreateUrlParams{
		Hash: shurl.RandStringRunes(8),
		Url:  url,
	}

	s, err := dq.CreateUrl(ctx, params)
	if err != nil {
		log.WithError(err).Error("Failed to register shortened URL")
	}

	fmt.Printf("%+v", s)
	fmt.Fprint(w, s.Hash)
}

func shortened(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	log.Infof("requested hash: %s", hash)

	ctx := context.Background()
	s, err := dq.GetUrlFromHash(ctx, hash)
	if err != nil {
		log.WithError(err).Warn("Failed to locate hash")
		t := template.Must(template.ParseFiles("static/form.html"))
		err := t.Execute(w, struct {NotFound bool}{true})
		if err != nil {
			http.Error(w, "Something broke", 500)
		}
		return
	}

	err = dq.UpdateUrl(ctx, db.UpdateUrlParams{
		ID:   s.ID,
		Hits: s.Hits + 1,
	})
	if err != nil {
		log.WithError(err).Warnf("Failed to update shurl hits for %s", s.ID)
	}
	http.Redirect(w, r, s.Url, 301)
}
