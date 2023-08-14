package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobintern/GoChat/service"
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
	uid := service.UID{UID: c.Query("uid")}

	if len(uid.UID) > 0 {
		data := uid.GetUser()
		c.HTML(http.StatusOK, "chatroom.html", data)
	} else {
		GetLogin(c)
	}
}
