package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template
var dbm dbManager

func init() {
	var err error
	tpl = template.Must(template.ParseGlob("template/*.html"))
	dbm, err = newDBManager("config.json")
	if err != nil {
		log.Fatalf("Fail to construcut db manager err: %s", err.Error())
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/home", home)
	r.HandleFunc("/calendar", calendar)
	r.HandleFunc("/_checkUserExist", checkUserExist)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/addProblem/{problemID}", addProblem)
	r.HandleFunc("/DoneProblem/{target}", doneProblem)
	r.HandleFunc("/_getUndo", getUndo)
	r.HandleFunc("/_getDateEvent", getDateEvent)
	r.HandleFunc("/deleteLog/{target}", deleteLog)
	http.Handle("/", r)

	// deploy asset
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/favicon.ico", favicon)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
