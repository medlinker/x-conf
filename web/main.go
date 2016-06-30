package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	"x-conf/client/goclient"
	_ "x-conf/web/route"
	"x-conf/web/utils"

	"github.com/sosop/libconfig"
)

func init() {

	xconf := flag.String("conf", "x-conf.conf", "config file")
	flag.Parse()
	goclient.IniConf = libconfig.NewIniConfig(*xconf)
	goclient.Config()

	// 创建share
	for _, env := range utils.Envs {
		err := goclient.CreateDir(goclient.MakeKey("share", env))
		if err == nil {
			goclient.CreateDir(goclient.MakeKey("prjs", "share"))
			goclient.Set(goclient.MakeKey("publish", "share", env), fmt.Sprint(time.Now().UnixNano()), nil)
		}
	}

	initUser()
}

func main() {
	addr := goclient.IniConf.GetString("addr", ":8000")
	http.ListenAndServe(addr, nil)
}

func initUser() {
	_, err := os.Stat(".lock")
	if err != nil && os.IsNotExist(err) {
		goclient.Set(goclient.MakeKey("users", "admin"), "098f6bcd4621d373cade4e832627b4f6", nil)
		f, _ := os.Create(".lock")
		f.Close()
	}
}
