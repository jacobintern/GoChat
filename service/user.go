package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	mongoDB := ConnectionInfo{
		DBName:         "chat_db",
		CollectionName: "chat_acc",
	}
	collection := mongoDB.MongoDBcontext()
	objID, err := primitive.ObjectIDFromHex(u.UID)
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
func ValidUser(r *http.Request) *Acc {
	chatAcc := Acc{}
	mongoDB := ConnectionInfo{
		DBName:         "chat_db",
		CollectionName: "chat_acc",
	}
	r.ParseForm()
	acc := r.FormValue("acc")
	pswd := r.FormValue("pswd")
	collection := mongoDB.MongoDBcontext()
	filter := bson.M{"acc": acc, "pswd": pswd}
	collection.Find(context.Background(), filter)
	err := collection.FindOne(context.Background(), filter).Decode(&chatAcc)
	if err == nil {
		return &chatAcc
	}
	return &Acc{}
}

// CreateUser is
func CreateUser(r *http.Request) *mongo.InsertOneResult {
	mongoDB := ConnectionInfo{
		DBName:         "chat_db",
		CollectionName: "chat_acc",
	}
	r.ParseForm()
	collection := mongoDB.MongoDBcontext()
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

// SetUsrCookie is
func (acc *Acc) SetUsrCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    uuid.New().String(),
		Value:   acc.ID,
		Expires: time.Now().Add(time.Minute * time.Duration(5)),
	}
	http.SetCookie(w, &cookie)
}
