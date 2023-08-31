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
	uid, ok := c.Params.Get("uid")

	if ok {
		user := service.UserInfo{UID: uid}
		user.GetUser()
		c.HTML(http.StatusOK, "chatroom.html", user)
	} else {
		GetLogin(c)
	}
}
