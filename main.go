package main

import (
	presener "mvcModule/presener"
	"net/http"
)

func main() {
	_ = presener.Init()

	http.ListenAndServe(":8080", nil)
}
