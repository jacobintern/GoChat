package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jacobintern/GoChat/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *gin.Context) {
	user := &service.Acc{}
	user.Acc = c.PostForm("acc")
	user.Pswd = c.PostForm("pswd")
	service.ValidUser(user)

	if user.ID != "" {
		c.SetCookie(uuid.New().String(), user.ID, 10, "", "", true, true)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "login successs", "uid": user.ID})
		return
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "login fail"})
		return
	}
}

func Register(c *gin.Context) {
	if len(service.CreateUser(c).InsertedID.(primitive.ObjectID).Hex()) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "registation successs", "success": true})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "registation failed", "success": false})
}

// GetUsers is
func GetUsers(c *gin.Context) {
	userList := service.Broadcaster.GetUserList()

	c.IndentedJSON(http.StatusOK, userList)
}

// GetUsrCookies is
func GetUsrCookies(c *gin.Context) {
	// for _, cookie := range req.Cookies() {
	// 	fmt.Println("Found a cookie named:", cookie.Name)
	// 	fmt.Println("Found a cookie expired:", cookie.Expires)
	// }
	// r, err := json.Marshal(req.Cookies())
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Fprint(w, string(r))
}
