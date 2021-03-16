package main

import (
	"log"
	"net/http"

	"github.com/jacobintern/GoChat/controllers"
)

func main() {
	// page
	controllers.RegisterPage()
	// api
	controllers.UserAPI()
	controllers.GetCookies()
	// websocke
	controllers.RegisterchatHandler()

	// static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// if any err log
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
