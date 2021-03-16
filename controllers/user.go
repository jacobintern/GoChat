package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jacobintern/GoChat/service"
)

// GetUsers is
func GetUsers(w http.ResponseWriter, req *http.Request) {
	userList := service.Broadcaster.GetUserList()
	r, err := json.Marshal(userList)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprint(w, string(r))
}

// GetUsrCookies is
func GetUsrCookies(w http.ResponseWriter, req *http.Request) {
	r, err := json.Marshal(req.Cookies())

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprint(w, string(r))
}

// UserAPI is
func UserAPI() {
	http.HandleFunc("/api/GetUserList", GetUsers)
}

// GetCookies is
func GetCookies() {
	http.HandleFunc("/api/GetCookies", GetUsrCookies)
}
