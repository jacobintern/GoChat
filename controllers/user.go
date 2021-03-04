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
	} else {
		fmt.Fprint(w, string(r))
	}
}

// UserAPI is
func UserAPI() {
	http.HandleFunc("api/GetUserList", GetUsers)
}
