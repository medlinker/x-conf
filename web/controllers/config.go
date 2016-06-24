package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
	"x-conf/client/goclient"
	"x-conf/web/utils"

	"github.com/coreos/etcd/client"
	"github.com/sosop/libconfig"
)

// ConfigPage 配置页
func ConfigPage(w http.ResponseWriter, r *http.Request) {
	validSess(w, r)
	t, _ := template.ParseFiles("views/config.html")
	t.Execute(w, nil)
}

// CreatePrj project 创建
func CreatePrj(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()
	if r.Method == "POST" {
		prjName := strings.TrimSpace(r.PostFormValue("prjName"))
		if utils.CheckParamsErr(&ret, prjName) {
			goto OVER
		}
		err := goclient.CreateDir("/prjs/" + prjName)
		if utils.CheckErr(err, &ret) {
			goto OVER
		}
		for _, env := range utils.Envs {
			err = goclient.CreateDir(utils.MakeKey(prjName, env))
			utils.CheckErr(err, &ret)
			goclient.Set(utils.MakeKey("publish", prjName, env), fmt.Sprint(time.Now().UnixNano()), nil)
		}
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}

// PrjList 项目列表
func PrjList(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()
	resp, err := goclient.Get("/prjs", &client.GetOptions{Recursive: true})
	prjs := make([]string, 0, 64)
	if utils.CheckErr(err, &ret) {
		goto OVER
	}
	for _, node := range resp.Node.Nodes {
		prjs = append(prjs, strings.Replace(node.Key, "/prjs/", "", -1))
	}
	ret.Data = prjs
OVER:
	utils.Output(w, ret)
}

// Confs 获取配置
func Confs(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()

	env := strings.TrimSpace(r.FormValue("env"))
	prjName := strings.TrimSpace(r.FormValue("prjName"))

	dir := utils.MakeKey(prjName, env)

	resp, err := goclient.Get(dir, &client.GetOptions{Recursive: true})
	confs := make(map[string]string, 128)

	if utils.CheckParamsErr(&ret, env, prjName) {
		goto OVER
	}

	if utils.CheckErr(err, &ret) {
		goto OVER
	}
	for _, node := range resp.Node.Nodes {
		confs[strings.Replace(node.Key, dir+"/", "", -1)] = node.Value
	}
	ret.Data = confs
OVER:
	utils.Output(w, ret)
}

// CreateConf 创建配置
func CreateConf(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()
	if r.Method == "POST" {
		r.ParseForm()
		env := strings.TrimSpace(r.PostFormValue("env"))
		prjName := strings.TrimSpace(r.PostFormValue("prjName"))
		key := strings.TrimSpace(r.PostFormValue("key"))
		value := strings.TrimSpace(r.PostFormValue("value"))
		if utils.CheckParamsErr(&ret, env, prjName, key, value) {
			goto OVER
		}
		_, err := goclient.Set(utils.MakeKey(prjName, env, key), value, nil)
		utils.CheckErr(err, &ret)
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}

// CreateBatchConf 批量创建配置
func CreateBatchConf(w http.ResponseWriter, r *http.Request) {
	validSess(w, r)
	if r.Method == "POST" {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err)
			return
		}
		file, _, err := r.FormFile("confsFile")
		if utils.CheckErr(err, nil) {
			fmt.Println(err)
			return
		}
		defer file.Close()
		env := strings.TrimSpace(r.FormValue("env"))
		prjName := strings.TrimSpace(r.FormValue("prjName"))
		iniConf := libconfig.NewIniConfigAsReader(file)
		for k, v := range iniConf.Entry {
			goclient.Set(utils.MakeKey(prjName, env, k), v.(string), nil)
		}
		t, _ := template.ParseFiles("views/config.html")
		t.Execute(w, nil)
	}

}

// DownloadConfs 下载配置
func DownloadConfs(w http.ResponseWriter, r *http.Request) {
	validSess(w, r)
	r.ParseForm()
	env := strings.TrimSpace(r.FormValue("env"))
	prjName := strings.TrimSpace(r.FormValue("prjName"))
	key := utils.MakeKey(prjName, env)
	resp, _ := goclient.Get(key, &client.GetOptions{Recursive: true})
	data := ""
	for _, node := range resp.Node.Nodes {
		data += fmt.Sprintln(strings.Replace(node.Key, key+"/", "", -1), "=", node.Value)
	}
	w.Header().Add("Content-Disposition", "attachment; filename="+"configure.properties")
	w.Write([]byte(data))
}

// DeleteConf 删除配置
func DeleteConf(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()
	env := strings.TrimSpace(r.PostFormValue("env"))
	prjName := strings.TrimSpace(r.PostFormValue("prjName"))
	key := strings.TrimSpace(r.PostFormValue("key"))
	if utils.CheckParamsErr(&ret, env, prjName, key) {
		goto OVER
	}
	goclient.Delete(utils.MakeKey(prjName, env, key), nil)
OVER:
	utils.Output(w, ret)
}

// Publish 发布配置
func Publish(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	validSess(w, r)
	ret := utils.NewRet()
	env := strings.TrimSpace(r.PostFormValue("env"))
	prjName := strings.TrimSpace(r.PostFormValue("prjName"))
	if utils.CheckParamsErr(&ret, env, prjName) {
		goto OVER
	}
	goclient.Update(utils.MakeKey("publish", prjName, env), fmt.Sprint(time.Now().UnixNano()))
OVER:
	utils.Output(w, ret)
}

func validSess(w http.ResponseWriter, r *http.Request) {
	valid, sess := utils.CheckSessFromCookie(r)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
	sess.Update()
}
