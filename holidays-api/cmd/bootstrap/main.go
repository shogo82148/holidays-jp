package main

import (
	"net/http"
	_ "time/tzdata"

	holidays "github.com/shogo82148/holidays-jp/holidays-api"
	"github.com/shogo82148/ridgenative"
)

func main() {
	h := holidays.NewHandler()
	http.Handle("/", h)
	ridgenative.ListenAndServe(":8080", nil)
}
