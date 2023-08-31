package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jacobintern/GoChat/controllers"
	"github.com/jacobintern/GoChat/service"
)

func main() {
	r := gin.Default()
	// page
	page := r.Group("")
	{
		page.GET("/login", controllers.GetLogin)
		page.GET("/register", controllers.GetRegister)
		page.GET("/chatroom/:uid", controllers.GetRoom)
	}

	// api
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("register", controllers.Register)
		api.GET("/GetUserList", controllers.GetUsers)
		api.GET("/GetCookies", controllers.GetUsrCookies)
	}

	// websocket
	go service.Hub.Run()
	r.GET("/ws", controllers.HandShake)

	// static
	r.LoadHTMLGlob("views/*")
	r.Static("/static", "./static")

	// if any err log
	err := r.Run(":8888")
	if err != nil {
		fmt.Println(err)
	}
}
