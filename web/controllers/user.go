package controllers

import (
	"net/http"
	"strings"
	"text/template"
	"time"
	wm "x-conf/web/models"
	"x-conf/web/utils"
)

// Login 用户登陆
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		name := strings.TrimSpace(r.PostFormValue("name"))
		pass := strings.TrimSpace(r.PostFormValue("pass"))
		if utils.CheckParamsErr(nil, name, pass) {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}
		user := wm.User{Name: name, Pass: pass}
		if user.ValidPass() {
			uuid := utils.GenerateUUID()
			utils.SessMap[uuid] = utils.NewSession(uuid, time.Minute*30, time.Now().UnixNano(), user)
			cookie := http.Cookie{Name: "SESSIONID", Value: uuid, Path: "/", MaxAge: 86400}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/x/conf/config", http.StatusMovedPermanently)
		} else {
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}
	} else if r.Method == "GET" {
		if v, _ := utils.CheckSessFromCookie(r); v {
			http.Redirect(w, r, "/x/conf/config", http.StatusMovedPermanently)
		}
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	}
}

// Logout 下线
func Logout(w http.ResponseWriter, r *http.Request) {
	_, s := utils.CheckSessFromCookie(r)
	s.Delete()
	t, _ := template.ParseFiles("views/login.html")
	t.Execute(w, nil)
}
