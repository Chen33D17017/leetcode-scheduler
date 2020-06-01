package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const dateFormat = "2006-01-02"

type Problem struct {
	problemQuery
	Deadline    string
	ProblemType []string
	ReviewLevel int
	Done        bool
	LogID       int
}

type DoneLog struct {
	LogID    int
	userID   int
	CostTime string `json:"costTime"`
	Level    string `json:"level"`
}

func home(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		log.Fatalln("Home Page err")
	}
}

func addProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["problemID"]
	val, err := strconv.Atoi(problemID)
	if err != nil {
		http.Error(w, "Wrong Search Condition", http.StatusServiceUnavailable)
		return
	}
	if r.Method == "POST" {
		c, err := r.Cookie("session")
		if err != nil {
			log.Println("Fail to get session on cookie")
		}
		userID := dbm.getUserID(c.Value)
		queryRst, err := dbm.problemInfo(val)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Problem Not Exist", http.StatusServiceUnavailable)
			return
		}
		newProblem := Problem{
			queryRst, time.Now().Format(dateFormat), nil, 0, false, 0,
		}

		err = dbm.insertProblem(userID, &newProblem)
		if err != nil {
			log.Printf("Insert new Problem err: %s", err.Error())
			http.Error(w, "Aleady in your schedule queue", http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(newProblem)
	}
}

func doneProblem(w http.ResponseWriter, r *http.Request) {
	var rData DoneLog
	vars := mux.Vars(r)
	c, err := r.Cookie("session")
	if err != nil {
		log.Println("Fail to get session on cookie")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rData.userID = dbm.getUserID(c.Value)
	if r.Method == "POST" {

		logID, err := strconv.Atoi(vars["target"])
		if err != nil {
			log.Printf("doneProblem: Wrong type of target: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rData.LogID = logID
		err = json.NewDecoder(r.Body).Decode(&rData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = dbm.doneProblem(rData)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = dbm.addProblemLog(rData)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

}

func getUndo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c, err := r.Cookie("session")
		if err != nil {
			log.Printf("Get cookie err: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID := dbm.getUserID(c.Value)
		rst, err := dbm.getProblems(userID)
		if err != nil {
			log.Fatalf("Fail to get problem sets: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rst)
	}
}

func deleteLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID, err := strconv.Atoi(vars["target"])
	if err != nil {
		log.Printf("doneProblem: Wrong type of target: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		log.Println("Fail to get session on cookie")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := dbm.getUserID(c.Value)
	err = dbm.deleteLog(logID, userID)
	if err != nil {
		log.Printf("deleteLog: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
