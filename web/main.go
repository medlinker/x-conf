package main

import (
	"fmt"
	"net/http"
	"time"
	"x-conf/client/goclient"
	_ "x-conf/web/route"
	"x-conf/web/utils"
)

func init() {
	// 创建share
	for _, env := range utils.Envs {
		err := goclient.CreateDir(goclient.MakeKey("share", env))
		if err == nil {
			goclient.CreateDir(goclient.MakeKey("prjs", "share"))
			goclient.Set(goclient.MakeKey("publish", "share", env), fmt.Sprint(time.Now().UnixNano()), nil)
		}
	}
}

func main() {
	addr := goclient.IniConf.GetString("addr", ":8000")
	http.ListenAndServe(addr, nil)
}
