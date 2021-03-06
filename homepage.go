package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func homepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println("Serving Homepage")
	http.ServeFile(writer, request, "./html/homepage.html")
}
