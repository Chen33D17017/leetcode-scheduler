package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type SummaryData struct {
	UserName     string
	RemainNum    int
	DoneNum      int
	CostTimeDate []string
	DateLabel    []string
	DateDone     []int
}

func summary(w http.ResponseWriter, r *http.Request) {
	if !alreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	c, err := r.Cookie("session")
	if err != nil {
		log.Println("Fail to get session on cookie")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := dbm.getUserID(c.Value)
	data := getSummaryDate(userID)
	err = tpl.ExecuteTemplate(w, "summary.html", data)
}

func getSummaryDate(userID int) SummaryData {
	rst := SummaryData{}
	rst.DateLabel = getLast7Date()

	userName, err := dbm.getUserName(userID)
	if err != nil {
		log.Fatalf("getSummaryDate err: Fail to get username %s\n", err.Error())
	}
	rst.UserName = userName

	remainNum, err := dbm.getUndoNum(userID)
	if err != nil {
		log.Fatalf("getSummaryDate err: Fail to get undoNum %s\n", err.Error())
	}
	rst.RemainNum = remainNum

	doneNum, err := dbm.getTotalProblem(userID)
	if err != nil {
		log.Fatalf("getSummaryDate err: get done problem number %s", err.Error())
	}
	rst.DoneNum = doneNum
	rst.DateDone = getDateDoneNum(userID, rst.DateLabel)
	rst.CostTimeDate = getDateCostTime(userID, rst.DateLabel)
	return rst
}

func time2Float(target string) float32 {
	minuteFormat := "04:05"
	minuteTime := strings.Join(strings.Split(target, ":")[:2], ":")
	rst, _ := time.Parse(minuteFormat, minuteTime)
	rstTime := float32(rst.Minute()) + float32(rst.Second())/float32(60)
	return rstTime
}

func getDateCostTime(userID int, dateLabel []string) []string {
	costTimeRst := make([]string, 7)
	for i, dateLabel := range dateLabel {
		tmpRst := float32(0)
		datetimes, err := dbm.getDateCostTime(userID, dateLabel)
		if err != nil {
			log.Fatalf("getDateCostTime err: %s \n", err.Error())
			return costTimeRst
		}

		for _, datetime := range datetimes {
			tmpRst += time2Float(datetime)
		}
		costTimeRst[i] = fmt.Sprintf("%.2f", tmpRst)
	}
	return costTimeRst
}

func getLast7Date() []string {
	rst := make([]string, 7)
	today := time.Now()
	for i := 6; i >= 0; i-- {
		rst[6-i] = today.AddDate(0, 0, -i).Format(dateFormat)
	}
	return rst
}

func getDateDoneNum(userID int, target []string) []int {
	rst := make([]int, 7)
	for i, t := range target {
		tmp, err := dbm.getDateDoneNum(userID, t)
		if err != nil {
			log.Fatalf("getDoneNumFromDate err @ date %s: %s", t, err.Error())
			return rst
		}
		rst[i] = tmp
	}
	return rst
}
