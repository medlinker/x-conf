package route

import (
	"net/http"
	"x-conf/web/controllers"
)

func init() {
	http.HandleFunc("/x/conf/prj", controllers.CreatePrj)
}
