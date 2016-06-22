package controllers

import (
	"net/http"
	"strings"
	"x-conf/client/goclient"
	"x-conf/web/utils"

	"github.com/coreos/etcd/client"
	"github.com/sosop/libconfig"
)

// CreatePrj project 创建
func CreatePrj(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
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
			err = goclient.CreateDir("/" + prjName + "/" + env)
			utils.CheckErr(err, &ret)
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

// CreateConf 创建配置
func CreateConf(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
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
		_, err := goclient.Set("/"+prjName+"/"+env+"/"+key, value, nil)
		utils.CheckErr(err, &ret)
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}

// CreateBatchConf 批量创建配置
func CreateBatchConf(w http.ResponseWriter, r *http.Request) {
	utils.Header(w)
	ret := utils.NewRet()
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		env := strings.TrimSpace(r.PostFormValue("env"))
		prjName := strings.TrimSpace(r.PostFormValue("prjName"))
		file, _, err := r.FormFile("upload")
		if utils.CheckErr(err, &ret) {
			goto OVER
		}
		defer file.Close()
		iniConf := libconfig.NewIniConfigAsReader(file)
		for k, v := range iniConf.Entry {
			_, err = goclient.Set("/"+prjName+"/"+env+"/"+k, v.(string), nil)
			if utils.CheckErr(err, &ret) {
				goto OVER
			}
		}
	} else {
		utils.SetMethodErr(&ret)
	}
OVER:
	utils.Output(w, ret)
}
