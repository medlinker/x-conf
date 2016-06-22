package controllers

import (
	"net/http"
	"strings"
	wm "x-conf/web/models"
	"x-conf/web/utils"
)

// Login 用户登陆
func Login(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	ret := utils.NewRet()
	if r.Method == "POST" {
		r.ParseForm()
		name := strings.TrimSpace(r.PostFormValue("name"))
		pass := strings.TrimSpace(r.PostFormValue("pass"))
		if utils.CheckParamsErr(&ret, name, pass) {
			goto OVER
		}
		user := wm.User{Name: name, Pass: pass}
		if user.ValidPass() {
			// TODO cookies
		} else {
			ret.Code = utils.ErrUser
			ret.Msg = "username or password error"
		}
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}
