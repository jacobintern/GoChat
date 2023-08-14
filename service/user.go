package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"
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
func (u *User) NewUser() *User {
	uid := UID{UID: u.Conn.Request().URL.Query().Get("clientId")}
	u.UserInfo = uid.GetUser()
	u.MessageChannel = make(chan *Message)
	return u
}

// SendMessage is
func (u *User) SendMessage() {
	for msg := range u.MessageChannel {
		if msg.ToUser != nil {
			if msg.ToUser.UID == u.UserInfo.UID {
				tmp := Message{
					Content: "This secret message comes from <font color='red'>" + msg.User.UserInfo.Name + "</font> says : " + msg.Content,
					ToUser:  msg.ToUser,
					User:    msg.User,
					Type:    msg.Type,
				}

				r, err := json.Marshal(tmp)
				if err != nil {
					fmt.Println(err)
					log.Fatal(err)
				}
				websocket.Message.Send(u.Conn, string(r))
			} else if msg.User.UserInfo.UID == u.UserInfo.UID {
				tmp := Message{
					Content: "You sent a secret message to <font color='red'>" + msg.ToUser.Name + "</font> says : " + msg.Content,
					ToUser:  msg.ToUser,
					User:    msg.User,
					Type:    msg.Type,
				}
				r, err := json.Marshal(tmp)
				if err != nil {
					fmt.Println(err)
					log.Fatal(err)
				}
				websocket.Message.Send(u.Conn, string(r))
			}
		} else {
			r, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}
			websocket.Message.Send(u.Conn, string(r))
		}
	}
}

// ReceiveMessage is
func (u *User) ReceiveMessage() error {
	for {
		var tmp string
		reply := Message{}
		if err := websocket.Message.Receive(u.Conn, &tmp); err != nil {
			return err
		}
		reply.User = u
		// 解析json
		json.Unmarshal([]byte(tmp), &reply)

		// 内容发送到聊天室
		sendMsg := reply.NewMessage()
		Broadcaster.Broadcast(sendMsg)
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
