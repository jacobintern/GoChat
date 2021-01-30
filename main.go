package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello jacob")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/login.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", loginPage)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
