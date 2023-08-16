package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UID is
type UID struct {
	UID string
}

// User is
type User struct {
	UserInfo       *UserInfo     `json:"user_info"`
	MessageChannel chan *Message `json:"-"`

	Conn *websocket.Conn
}

// UserInfo is
type UserInfo struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func DbContext() ConnectionInfo {
	return ConnectionInfo{
		DBName:         "chat_db",
		CollectionName: "chat_acc",
	}
}

// NewUser is
func (u *User) NewUser(userID string) *User {
	if len(userID) == 0 {
		fmt.Println("lost clientID")
	}
	uid := UID{UID: userID}
	u.UserInfo = uid.GetUser()
	u.MessageChannel = make(chan *Message)
	return u
}

// SendMessage is
func (u *User) SendMessage() {
	for msg := range u.MessageChannel {
		if msg.ToUser != nil {
			// 密語
			if msg.ToUser.UID == u.UserInfo.UID {
				tmp := Message{
					Content: "From <font color='red'>" + msg.User.UserInfo.Name + "</font> says : " + msg.Content,
					ToUser:  msg.ToUser,
					User:    msg.User,
					Type:    msg.Type,
				}

				r, err := json.Marshal(tmp)
				if err != nil {
					log.Fatal(err)
				}
				u.Conn.WriteMessage(websocket.TextMessage, r)
			} else if msg.User.UserInfo.UID == u.UserInfo.UID {
				tmp := Message{
					Content: "To <font color='red'>" + msg.ToUser.Name + "</font> says : " + msg.Content,
					ToUser:  msg.ToUser,
					User:    msg.User,
					Type:    msg.Type,
				}
				r, err := json.Marshal(tmp)
				if err != nil {
					log.Fatal(err)
				}
				u.Conn.WriteMessage(websocket.TextMessage, r)
			}
		} else {
			// 一般公頻
			r, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err)
			}
			u.Conn.WriteMessage(websocket.TextMessage, r)
		}
	}
}

// ReceiveMessage is
func (u *User) ReceiveMessage() error {
	for {
		// var tmp string
		reply := Message{}
		_, p, err := u.Conn.ReadMessage()

		if err != nil {
			return err
		}

		// 解析json
		err = json.Unmarshal([]byte(p), &reply)
		reply.User = u

		if err != nil {
			return err
		}

		// 内容发送到聊天室
		sendMsg := reply.NewMessage()
		Hub.Broadcast(sendMsg)
	}
}

// GetUser is
func (u UID) GetUser() *UserInfo {
	chatAcc := Acc{}
	collection := DbContext().MongoDBcontext()
	objID, err := primitive.ObjectIDFromHex(u.UID)
	if err != nil {
		fmt.Println(err)
	}
	filter := bson.M{"_id": objID}
	err = collection.FindOne(context.Background(), filter).Decode(&chatAcc)
	if err == nil {
		return &UserInfo{
			UID:  chatAcc.ID,
			Name: chatAcc.Name,
		}
	}
	return &UserInfo{}
}

// ValidUser is checkout login user exist in mongodb
func ValidUser(user *Acc) {
	collection := DbContext().MongoDBcontext()
	filter := bson.M{"acc": user.Acc, "pswd": user.Pswd}
	collection.Find(context.Background(), filter)
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Print(err)
	}
}

// CreateUser is
func CreateUser(c *gin.Context) *mongo.InsertOneResult {
	collection := DbContext().MongoDBcontext()
	res, err := collection.InsertOne(context.Background(), Acc{
		Acc:    c.PostForm("acc"),
		Pswd:   c.PostForm("pswd"),
		Name:   c.PostForm("name"),
		Email:  c.PostForm("email"),
		Gender: c.PostForm("gender"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return res
}
