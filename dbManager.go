package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	// "strconv"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type config struct {
	StoreType string `json:"type"`
	ID        string `json:"id"`
	EndPoint  string `json:"endpoint"`
	PW        string `json:"password"`
	DB        string `json:"database"`
}

type problemQuery struct {
	ID           int `json:"id"`
	DifficultyID int `json:"level_id"`
	Difficulty   string
	ProblemName  string `json:"problem_name"`
	URL          string `json:"url"`
}

type dbManager struct {
	*sql.DB
}

func checkErr(err error) {
	log.Println(err)
}

func newDBManager(fileName string) (dbManager, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal("Fail to open config file")
	}

	configData := config{}
	err = json.Unmarshal([]byte(file), &configData)
	if err != nil {
		log.Fatalf("Read json err: %s", err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", configData.ID, configData.PW, configData.EndPoint, configData.DB)
	db, err := sql.Open(configData.StoreType, dsn)
	if err != nil {
		return dbManager{}, err
	}

	dbm := dbManager{db}

	// check for connection
	err = dbm.Ping()
	if err != nil {
		log.Fatalf("Fail to connect to mysql: %s", err.Error())
		return dbManager{}, err
	}

	return dbm, nil
}

func (dbm *dbManager) userRegist(username, email, password string) error {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	_, err = dbm.Exec("INSERT INTO `leetcode_user`(user_name, email, password) VALUES(?, ?, ?)", username, email, bs)
	if err != nil {
		return fmt.Errorf("register err : %s", err.Error())
	}
	return nil
}

func (dbm *dbManager) setUserSession(userEmail, session string) error {
	stmt, err := dbm.Prepare("update leetcode_user set `session`=? where email=?")
	if err != nil {
		return fmt.Errorf("setSession err : %s", err.Error())
	}

	_, err = stmt.Exec(session, userEmail)
	if err != nil {
		return fmt.Errorf("setSession err : %s", err.Error())
	}
	return nil
}

func (dbm *dbManager) delUserSession(userID int) error {
	stmt, err := dbm.Prepare("UPDATE `leetcode_user` set `session`=? where user_id=?")
	if err != nil {
		return fmt.Errorf("delUserSession err : %s", err.Error())
	}
	_, err = stmt.Exec(nil, userID)
	if err != nil {
		return fmt.Errorf("setSession err: %s", err.Error())
	}
	return nil
}

func (dbm *dbManager) getUserSession(userEmail string) (string, error) {
	var rst string
	err := dbm.QueryRow("SELECT session on `leetcode_user` WHERE email=?", userEmail).Scan(&rst)
	if err != nil {
		return "", fmt.Errorf("getUserSession err : %s", err.Error())
	}

	return rst, nil
}

func (dbm *dbManager) getUserID(session string) int {
	var rst int
	err := dbm.QueryRow("SELECT user_id FROM `leetcode_user` WHERE `session`=? LIMIT 1", session).Scan(&rst)
	if err != nil {
		return 0
	}
	return rst
}

func (dbm *dbManager) checkLogin(email, password string) bool {
	psRst := make([]byte, 60)
	err := dbm.QueryRow("SELECT `password` FROM `leetcode_user` WHERE email=?", email).Scan(&psRst)
	if err != nil {
		log.Printf("checkLogin err : %s", err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword(psRst, []byte(password))
	if err != nil {
		log.Printf("Wrong password with %s: %s", email, err)
		return false
	}

	return true
}

func (dbm *dbManager) checkExist(problemID int) bool {
	rows, err := dbm.Query("SELECT COUNT(*) FROM `leetcode_problem` WHERE id=?", problemID)
	if err != nil {
		fmt.Printf("checkExist: err %s", err.Error())
		return false
	}

	defer rows.Close()
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("checkExist: result err : %s", err.Error())
			return false
		}
	}

	if id != 0 {
		return true
	}
	return false
}

func (dbm *dbManager) problemInfo(problemID int) (problemQuery, error) {
	var rst problemQuery
	var err error
	err = dbm.QueryRow("SELECT * FROM `leetcode_problem` WHERE id=?", problemID).Scan(&rst.ID, &rst.DifficultyID, &rst.ProblemName, &rst.URL)
	if err != nil {
		return rst, fmt.Errorf("problemInfo err: Problem Not Exist %s", err.Error())
	}

	err = dbm.QueryRow("SELECT level FROM `problem_level` WHERE id=?", rst.DifficultyID).Scan(&rst.Difficulty)
	if err != nil {
		return rst, fmt.Errorf("problemInfo err: unexpected difficulty %s", err.Error())
	}
	return rst, nil
}

func (dbm *dbManager) insertProblem(userID int, p *Problem) error {
	var rst int
	stmt, err := dbm.Prepare("INSERT INTO `problem_log`(`user_id`, `problem_id`, `date`, `review_level`, `time`, `done`) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("insertProbelm prepare err : %s", err.Error())
	}
	_, err = stmt.Exec(userID, p.ID, p.Deadline, p.ReviewLevel, nil, false)
	if err != nil {
		return fmt.Errorf("setSession err: %s", err.Error())
	}

	err = dbm.QueryRow("SELECT `id` FROM `problem_log` WHERE user_id=? AND problem_id=? AND `date`=?", userID, p.ID, p.Deadline).Scan(&rst)
	if err != nil {
		return fmt.Errorf("insertProblem get id err: %s", err.Error())
	}
	p.LogID = rst
	return nil
}

func (dbm *dbManager) getProblems(userID int) ([]Problem, error) {
	rst := make([]Problem, 0)
	rows, err := dbm.Query("SELECT id, problem_id, date, review_level FROM `problem_log` WHERE `user_id`=? AND `date`<=CURDATE() AND `done`=False ORDER BY DATE(date) DESC;", userID)
	if err != nil {
		return rst, fmt.Errorf("getProblem err: query fail: %s", err.Error())
	}

	defer rows.Close()
	var id, pID, rLV int
	var date string

	for rows.Next() {
		err = rows.Scan(&id, &pID, &date, &rLV)
		if err != nil {
			return rst, fmt.Errorf("getProblems err: scan fail: %s", err.Error())
		}

		pINFO, err := dbm.problemInfo(pID)
		if err != nil {
			return rst, fmt.Errorf("getProblems err: get problem info %s", err.Error())
		}

		// TODO: Category of Problem
		rst = append(rst, Problem{pINFO, date, nil, rLV, false, id})
	}
	return rst, nil
}

func (pq problemQuery) String() string {
	return fmt.Sprintf("%s : %s : %s", pq.ProblemName, pq.Difficulty, pq.URL)
}

func (dbm *dbManager) doneProblem(data DoneLog) error {
	stmt, err := dbm.Prepare("UPDATE `problem_log` set `done`=?, `date`=?, time=? where id=? and user_id=?")
	if err != nil {
		return fmt.Errorf("doneProblem err : %s", err.Error())
	}
	_, err = stmt.Exec(true, time.Now().Format(dateFormat), strings.TrimSpace(data.CostTime), data.LogID, data.userID)
	if err != nil {
		return fmt.Errorf("setSession err: %s", err.Error())
	}

	return nil
}

func (dbm *dbManager) addProblemLog(prev DoneLog) error {
	var level, problemID int
	err := dbm.QueryRow("SELECT `review_level`, `problem_id`  FROM `problem_log` WHERE id=? AND user_id=?", prev.LogID, prev.userID).Scan(&level, &problemID)
	if err != nil {
		return fmt.Errorf("addProblemLog err: %s", err.Error())
	}

	switch prev.Level {
	case "easy":
		level++
	case "hard":
		if level > 0 {
			level--
		}
	}

	nextDate := time.Now().AddDate(0, 0, 3*(level+1)).Format(dateFormat)

	// INSERT NEW LOG
	stmt, err := dbm.Prepare("INSERT INTO `problem_log`(`user_id`, `problem_id`, `date`, `review_level`, `time`, `done`) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("addProblemLog insert prepare err : %s", err.Error())
	}
	_, err = stmt.Exec(prev.userID, problemID, nextDate, level, nil, false)
	if err != nil {
		return fmt.Errorf("addProblemLog insert exec err: %s", err.Error())
	}

	return nil
}

func (dbm *dbManager) deleteLog(logID, userID int) error {
	stmt, err := dbm.Prepare("DELETE FROM `problem_log` WHERE user_id=? AND id=?")
	if err != nil {
		return fmt.Errorf("deleteLog err: %s", err.Error())
	}

	_, err = stmt.Exec(userID, logID)
	if err != nil {
		return fmt.Errorf("deleteLog exec err: %s", err.Error())
	}

	return nil
}

func (dbm *dbManager) getDateEvent(userID int) ([]DateEvent, error) {
	rst := make([]DateEvent, 0)
	rows, err := dbm.Query("SELECT problem_name, color, date FROM `problem_log` JOIN `leetcode_problem` ON problem_id=leetcode_problem.id JOIN `problem_level` ON leetcode_problem.level_id=problem_level.id WHERE `user_id`=? and date <= CURDATE() and done=1 union SELECT problem_name, color, date FROM `problem_log` JOIN `leetcode_problem` ON problem_id=leetcode_problem.id JOIN `problem_level` ON leetcode_problem.level_id=problem_level.id WHERE `user_id`=? and date > CURDATE();", userID, userID)
	if err != nil {
		return rst, fmt.Errorf("getDateEvent: err %s", err.Error())
	}

	defer rows.Close()
	var catcher DateEvent
	for rows.Next() {
		err = rows.Scan(&catcher.Title, &catcher.BgColor, &catcher.Date)
		catcher.AllDay = true
		catcher.TColor = "#ffffff"
		rst = append(rst, catcher)
		if err != nil {
			return rst, fmt.Errorf("getDateEvent: result err : %s", err.Error())
		}
	}

	return rst, nil
}

func (dbm *dbManager) checkUserExist(email string) (int, error) {
	var rst int
	err := dbm.QueryRow("SELECT COUNT(*) FROM `leetcode_user` WHERE email=?;", email).Scan(&rst)
	if err != nil {
		return 1, fmt.Errorf("addProblemLog err: %s", err.Error())
	}
	return rst, nil
}

func (dbm *dbManager) getDateCostTime(userID int, date string) ([]string, error) {
	rst := make([]string, 0)
	rows, err := dbm.Query("SELECT time FROM `problem_log` WHERE `user_id`=? AND date=? AND done=1;", userID, date)
	if err != nil {
		return rst, fmt.Errorf("getDateCostTime err: query fail: %s", err.Error())
	}

	defer rows.Close()
	var practiceTime string

	for rows.Next() {
		err = rows.Scan(&practiceTime)
		if err != nil {
			return rst, fmt.Errorf("getDateCostTime err: scan fail: %s", err.Error())
		}

		rst = append(rst, practiceTime)
	}
	return rst, nil
}

func (dbm *dbManager) getTotalProblem(userID int) (int, error) {
	var rst int
	err := dbm.QueryRow("SELECT COUNT(DISTINCT problem_id) FROM `problem_log` WHERE `user_id`=? AND done=1;", userID).Scan(&rst)
	if err != nil {
		return 0, fmt.Errorf("dbm getTotalProblem err: %s", err.Error())
	}
	return rst, nil
}

func (dbm *dbManager) getDateDoneNum(userID int, sDate string) (int, error) {
	var rst int
	err := dbm.QueryRow("SELECT COUNT(*) FROM `problem_log` WHERE `user_id`=? AND date =? AND done=1;", userID, sDate).Scan(&rst)

	if err != nil {
		return 0, fmt.Errorf("dbm getDateDoneNum err : %s", err.Error())
	} else {
		return rst, nil
	}
}

func (dbm *dbManager) getUserName(userID int) (string, error) {
	var rst string
	err := dbm.QueryRow("SELECT user_name from `leetcode_user` WHERE `user_id`=?;", userID).Scan(&rst)
	if err != nil {
		return "", fmt.Errorf("dbm.getUserName err : %s", err.Error())
	} else {
		return rst, nil
	}
}

func (dbm *dbManager) getUndoNum(userID int) (int, error) {
	var rst int
	err := dbm.QueryRow("SELECT COUNT(*) FROM `problem_log` WHERE `user_id`=? AND date<=curdate() AND done=0;", userID).Scan(&rst)
	if err != nil {
		return 0, fmt.Errorf("dbm.getUndoNum err: %s", err.Error())
	} else {
		return rst, nil
	}
}
