package route

import (
	"net/http"
	"x-conf/web/controllers"
)

func init() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/x/conf/config", controllers.ConfigPage)
	http.HandleFunc("/x/conf/prjs", controllers.PrjList)
	http.HandleFunc("/x/conf/prj", controllers.CreatePrj)
	http.HandleFunc("/x/conf/", controllers.CreatePrj)
}
