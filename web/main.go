package main

import (
	"net/http"
	_ "x-conf/web/route"
)

func main() {
	http.ListenAndServe(":8000", nil)
}
