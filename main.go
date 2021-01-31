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
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./views/login.html"))
		tmpl.Execute(w, nil)
		break
	case "POST":
		r.ParseForm()
		fmt.Println("username:", r.Form["acc"])
		fmt.Println("password:", r.Form["pswd"])
		break
	case "PUT":
		fmt.Fprintf(w, "put")
		break
	case "Delete":
		fmt.Fprintf(w, "delete")
		break
	default:
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "register")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/register", register)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
