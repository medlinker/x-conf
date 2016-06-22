package route

import (
	"net/http"
	"x-conf/web/controllers"
)

func init() {
	http.HandleFunc("/x/conf/login", controllers.Login)
	http.HandleFunc("/x/conf/prjs", controllers.PrjList)
	http.HandleFunc("/x/conf/prj", controllers.CreatePrj)
}
