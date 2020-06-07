package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type DateEvent struct {
	Title   string `json:"title"`
	Date    string `json:"start"`
	AllDay  bool   `json:"allDay"`
	BgColor string `json:"color"`
	TColor  string `json:"textColor"`
}

func calendar(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "calendar.html", nil)
	if err != nil {
		log.Fatalln("Calendar Page err")
	}
}

func getDateEvent(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		log.Println("Fail to get session on cookie")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID := dbm.getUserID(c.Value)

	if r.Method == "POST" {
		rst, err := dbm.getDateEvent(userID)

		if err != nil {
			log.Fatalf("Fail to Get Date Event %s\n", err.Error())
		}

		json.NewEncoder(w).Encode(rst)

	}

}
