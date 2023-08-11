package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func GetRoom(c *gin.Context) {
	c.HTML(http.StatusOK, "chatroom.html", gin.H{
		"title": "ChatRoom",
	})
}

// case "POST":
// 	if len(service.CreateUser(r).InsertedID.(primitive.ObjectID).Hex()) > 0 {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	} else {
// 		http.Redirect(w, r, "/register", http.StatusSeeOther)
// 	}
// 	break

// ChatRoom is chat room
// case "GET":
// 	u := service.UID{UID: r.URL.Query().Get("uid")}
// 	if len(u.UID) > 0 {
// 		path, _ := filepath.Abs("views/chatroom.html")
// 		data := u.GetUser()
// 		tmpl := template.Must(template.ParseFiles(path))
// 		tmpl.Execute(w, data)
// 	} else {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	}
