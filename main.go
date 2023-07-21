package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jacobintern/GoChat/controllers"
)

func main() {
	r := mux.NewRouter()
	// page
	r.HandleFunc("/hello", controllers.HomePage).Methods(http.MethodGet)
	controllers.RegisterPage()
	// api
	r.HandleFunc("/api/GetUserList", controllers.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/GetCookies", controllers.GetUsrCookies).Methods(http.MethodGet)
	// websocket
	controllers.RegisterchatHandler()

	// register api
	http.Handle("/", r)

	// static
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// if any err log
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
