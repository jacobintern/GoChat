package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"
)

// User is
type User struct {
	UserInfo       *UserInfo     `json:"user_info"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

// UserInfo is
type UserInfo struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

// NewUser is
func NewUser(conn *websocket.Conn) *User {
	userInfo := GetUser(conn.Request().URL.Query().Get("clientId"))
	user := &User{
		UserInfo:       &userInfo,
		MessageChannel: make(chan *Message),
		conn:           conn,
	}
	return user
}

// SendMessage is
func (u *User) SendMessage() {
	for msg := range u.MessageChannel {
		if msg.ToUser != nil {
			if msg.ToUser.UID == u.UserInfo.UID {
				tmp := Message{
					Content: "This secret message comes from <font color='red'>" + u.UserInfo.Name + "</font> says : " + msg.Content,
					ToUser:  msg.ToUser,
					User:    msg.User,
					Type:    msg.Type,
				}

				r, err := json.Marshal(tmp)
				if err != nil {
					fmt.Println(err)
					log.Fatal(err)
				}
				websocket.Message.Send(u.conn, string(r))
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
				websocket.Message.Send(u.conn, string(r))
			}
		} else {
			r, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}
			websocket.Message.Send(u.conn, string(r))
		}
	}
}

// ReceiveMessage is
func (u *User) ReceiveMessage() error {
	for {
		var tmp string
		reply := Message{}
		if err := websocket.Message.Receive(u.conn, &tmp); err != nil {
			return err
		}
		// 解析json
		json.Unmarshal([]byte(tmp), &reply)

		// 内容发送到聊天室
		sendMsg := NewMessage(u, reply.Content, reply.ToUser)
		Broadcaster.Broadcast(sendMsg)
	}
}

// GetUser is
func GetUser(uid string) UserInfo {
	chatAcc := Acc{}
	collection := MongoDBcontext("chat_db", "chat_acc")
	objID, err := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": objID}
	err = collection.FindOne(context.Background(), filter).Decode(&chatAcc)
	if err == nil {
		return UserInfo{
			UID:  chatAcc.ID,
			Name: chatAcc.Name,
		}
	}
	return UserInfo{}
}

// ValidUser is checkout login user exist in mongodb
func ValidUser(r *http.Request) string {
	chatAcc := Acc{}
	r.ParseForm()
	acc := r.FormValue("acc")
	pswd := r.FormValue("pswd")
	collection := MongoDBcontext("chat_db", "chat_acc")
	filter := bson.M{"acc": acc, "pswd": pswd}
	collection.Find(context.Background(), filter)
	err := collection.FindOne(context.Background(), filter).Decode(&chatAcc)
	if err == nil {
		return chatAcc.ID
	}
	return ""
}

// CreateUser is
func CreateUser(r *http.Request) *mongo.InsertOneResult {
	r.ParseForm()
	collection := MongoDBcontext("chat_db", "chat_acc")
	res, err := collection.InsertOne(context.Background(), Acc{
		Acc:    r.FormValue("acc"),
		Pswd:   r.FormValue("pswd"),
		Name:   r.FormValue("name"),
		Email:  r.FormValue("email"),
		Gender: r.FormValue("gender"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return res
}
