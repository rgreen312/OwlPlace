package main

import (
	"fmt"
	"os"
	"github.com/rgreen312/owlplace/server/apiserver"
)


func main() {
	server := apiserver.NewApiServer(os.Getenv("MY_POD_IP"))
	server.ListenAndServe()
}
