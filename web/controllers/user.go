package controllers

import (
	"errors"
	"net/http"
	"strings"
	"text/template"
	"time"
	wm "x-conf/web/models"
	"x-conf/web/utils"

	"github.com/coreos/etcd/client"

	"x-conf/client/goclient"
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

// Create 创建用户
func Create(w http.ResponseWriter, r *http.Request) {
	validSess(w, r)
	utils.Header(w)
	ret := utils.NewRet()
	if r.Method == "POST" {
		username := strings.TrimSpace(r.PostFormValue("username"))
		if utils.CheckParamsErr(&ret, username) {
			goto OVER
		}
		key := goclient.MakeKey("users", username)
		goclient.Set(key, wm.EncrytPass("123456"), nil)
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}

// Modify 修改密码
func Modify(w http.ResponseWriter, r *http.Request) {
	validSess(w, r)
	utils.Header(w)
	ret := utils.NewRet()
	if r.Method == "POST" {
		oldPassword := strings.TrimSpace(r.PostFormValue("oldPassword"))
		newPassword := strings.TrimSpace(r.PostFormValue("newPassword"))
		if utils.CheckParamsErr(&ret, oldPassword, newPassword) {
			goto OVER
		}
		if exist, s := utils.CheckSessFromCookie(r); exist {
			_, err := goclient.Set(goclient.MakeKey("users", s.Data.(wm.User).Name), wm.EncrytPass(newPassword), &client.SetOptions{PrevValue: wm.EncrytPass(oldPassword)})
			if err != nil {
				utils.CheckErr(err, &ret)
				goto OVER
			}
		} else {
			utils.CheckErr(errors.New("user error"), &ret)
		}
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}
