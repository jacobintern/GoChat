package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func homePage(h http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(h, "Hello word")
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("login.jacob")
		log.Println(t.Execute(w, nil))
		break
	case "POST":
		fmt.Println("username:", r.Form["acc"])
		fmt.Println("password:", r.Form["pswd"])
	}
}

func main() {
	go http.HandleFunc("/", homePage)
	go http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
