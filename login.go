package main

import (
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type loginErr struct {
	EmailErr    string
	PasswordErr string
}

func login(w http.ResponseWriter, r *http.Request) {
	tpr := loginErr{}

	if alreadyLogin(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	http.SetCookie(w, c)

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		// TODO: check the password with data in database
		ok := dbm.checkLogin(email, password)
		if !ok {
			tpr.EmailErr = "Wrong Email or password"
			tpr.PasswordErr = "Wrong Email or password"
		} else {
			err = dbm.setUserSession(email, c.Value)
			if err != nil {
				log.Printf("Store Cookie err : %s", err.Error())
			}
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}
	err = tpl.ExecuteTemplate(w, "login.html", tpr)
	if err != nil {
		log.Fatal("login page err : {}", err.Error())
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if alreadyLogin(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := dbm.userRegist(username, email, password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Fatal("registration page err: {}", err.Error())
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("session")
	userID := dbm.getUserID(c.Value)
	err := dbm.delUserSession(userID)
	if err != nil {
		log.Printf("Fail to delete session on db : %s\n", err.Error())
	}
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func checkUserExist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postData := r.FormValue("msg")
		rst, err := dbm.checkUserExist(postData)
		if err != nil {
			log.Printf("Fail to check email on database: %s\n", err.Error())
		}
		if rst > 0 {
			w.Write([]byte("Not OK"))
		} else {
			w.Write([]byte("OK"))
		}
	}
}

func alreadyLogin(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	rst := dbm.getUserID(c.Value)
	if rst == 0 {
		return false
	}
	return true
}

// go get github.com/satori/go.uuid
// go get golang.org/x/crypto/bcrypt
