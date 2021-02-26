package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/websocket"
)

// Acc is
type Acc struct {
	ID     string `bson:"_id,omitempty"`
	Acc    string `bson:"acc"`
	Pswd   string `bson:"pswd"`
	Email  string `bson:"email"`
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
}

// UserInfo is
type UserInfo struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

// Message is
type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	ToUser  *UserInfo `json:"to_user"`
}

// User is
type User struct {
	UserInfo       *UserInfo     `json:"user_info"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

// Broadcaster is
type Broadcaster struct {
	users map[string]*User

	enterChannel        chan *User
	leaveChannel        chan *User
	messageChannel      chan *Message
	requestUsersChannel chan struct{}
	usersChannel        chan []*UserInfo
}

// Message type
const (
	MsgNormal = iota
	MsgSystem
	MsgSentUserList
)

var broadcaster = &Broadcaster{
	users:               make(map[string]*User),
	enterChannel:        make(chan *User),
	leaveChannel:        make(chan *User),
	messageChannel:      make(chan *Message),
	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*UserInfo),
}

// NewMessage is
func NewMessage(user *User, content string, toUser *UserInfo) *Message {
	return &Message{
		User:    user,
		Type:    MsgNormal,
		Content: content,
		ToUser:  toUser,
	}
}

// NewUserEnterMessage is
func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgSentUserList,
		Content: user.UserInfo.Name + " 加入了聊天室",
	}
}

// NewUserLeaveMessage is
func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgSentUserList,
		Content: user.UserInfo.Name + " 離開了聊天室",
	}
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

// UserEntering is
func (b *Broadcaster) UserEntering(u *User) {
	b.enterChannel <- u
}

// UserLeaving is
func (b *Broadcaster) UserLeaving(u *User) {
	b.leaveChannel <- u
}

// Broadcast is
func (b *Broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
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
		broadcaster.Broadcast(sendMsg)
	}
}

// GetUserList is
func (b *Broadcaster) GetUserList() []*UserInfo {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}

// Home is home page
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello jacob")
}

// LoginPage is login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./views/login.html"))
		tmpl.Execute(w, nil)
		break
	case "POST":
		if uid := ValidUser(r); len(uid) > 0 {
			http.Redirect(w, r, "/chatroom?uid="+uid, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		break
	case "PUT":
		fmt.Fprintf(w, "put")
		break
	case "Delete":
		fmt.Fprintf(w, "delete")
		break
	default:
		return
	}
	return
}

// Register is user
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./views/register.html"))
		tmpl.Execute(w, nil)
		break
	case "POST":
		if len(CreateUser(r).InsertedID.(primitive.ObjectID).Hex()) > 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}
		break
	case "PUT":
		fmt.Fprintf(w, "put")
		break
	case "Delete":
		fmt.Fprintf(w, "delete")
		break
	default:
		return
	}
}

// ChatRoom is chat room
func ChatRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if uid := r.URL.Query().Get("uid"); len(uid) > 0 {
			data := GetUser(uid)
			tmpl := template.Must(template.ParseFiles("./views/chatroom.html"))
			tmpl.Execute(w, data)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		break
	case "POST":
		fmt.Fprintf(w, "post")
		break
	case "PUT":
		fmt.Fprintf(w, "put")
		break
	case "Delete":
		fmt.Fprintf(w, "delete")
		break
	default:
		return
	}
}

// MongoDBcontext is connect setting
func MongoDBcontext(dbName string, collectionName string) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://j_dev:zHYJQ2jc7UAqHThV@jdev.y4x5s.gcp.mongodb.net/"+dbName+"?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(dbName).Collection(collectionName)
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

// Start is
func (b *Broadcaster) Start() {
	for {
		select {
		case user := <-b.enterChannel:
			b.users[user.UserInfo.Name] = user
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				user.MessageChannel <- msg
			}
		case user := <-b.leaveChannel:
			delete(b.users, user.UserInfo.Name)
			close(user.MessageChannel)
		case <-b.requestUsersChannel:
			userList := make([]*UserInfo, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user.UserInfo)
			}
			b.usersChannel <- userList
		}
	}
}

// Echo is
func Echo(conn *websocket.Conn) {
	// 建立使用者
	user := NewUser(conn)
	// 建立傳送訊息通道 goroutine監聽
	go user.SendMessage()

	// 使用者進入
	msg := NewUserEnterMessage(user)
	broadcaster.UserEntering(user)
	broadcaster.Broadcast(msg)

	// 訊息接收並傳送給其他使用者
	err := user.ReceiveMessage()

	// 使用者離開
	msg = NewUserLeaveMessage(user)
	broadcaster.UserLeaving(user)
	broadcaster.Broadcast(msg)

	if err == nil {
		conn.Close()
	} else {
		log.Println("read from client error:", err)
		conn.Close()
	}

}

// RegisterchatHandler is
func RegisterchatHandler() {
	go broadcaster.Start()

	http.Handle("/ws", websocket.Handler(Echo))
}

// GetUsers is
func GetUsers(w http.ResponseWriter, req *http.Request) {
	userList := broadcaster.GetUserList()
	r, err := json.Marshal(userList)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, string(r))
	}
}

func main() {
	// page
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/chatroom", ChatRoom)
	http.HandleFunc("/GetUserList", GetUsers)
	RegisterchatHandler()

	// static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// if any err log
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
