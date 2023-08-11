package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobintern/GoChat/service"
)

func Login(c *gin.Context) {
	user := &service.Acc{}
	user.Acc = c.PostForm("acc")
	user.Pswd = c.PostForm("pswd")
}

// case "POST":
// if acc := service.ValidUser(r); acc != nil {
// 	acc.SetUsrCookie(w)
// 	http.Redirect(w, r, "/chatroom?uid="+acc.ID, http.StatusSeeOther)
// } else {
// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
// }
// break

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
	// for _, cookie := range req.Cookies() {
	// 	fmt.Println("Found a cookie named:", cookie.Name)
	// 	fmt.Println("Found a cookie expired:", cookie.Expires)
	// }
	r, err := json.Marshal(req.Cookies())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprint(w, string(r))
}
