package main

import (
	"net/http"

	"github.com/rgreen312/owlplace/server/apiserver"
)

func main() {
	http.HandleFunc("/hello", apiserver.Hello)
	http.HandleFunc("/headers", apiserver.Headers)

	http.ListenAndServe(":3000", nil)
}
