package main

import (
	"container/list"
	"context"
	"encoding/hex"
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
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
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

// Broadcaster is
// type broadcaster struct {
// 	users map[string]*User

// 	// channel
// 	enterChannel   chan *User
// 	leaveChannel   chan *User
// 	messageChannel chan *Message
// }

// User is
type User struct {
	UID     string  `json:"client_id"`
	Name    string  `json:"usr_name"`
	Message Message `json:"msg"`
	//MessageChan chan *Message `json:"-"`

	conn *websocket.Conn
}

// Message is
type Message struct {
	ToID    string `json:"to_id"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// Message is
// type Message struct {
// 	User  *User            `json:"user"`
// 	Type  int              `json:"type"`
// 	Msg   string           `json:"msg"`
// 	Users map[string]*User `json:"users"`
// }

var conns = list.New()

// Broadcaster is
// var Broadcaster = &broadcaster{
// 	users: make(map[string]*User),

// 	enterChannel:   make(chan *User),
// 	leaveChannel:   make(chan *User),
// 	messageChannel: make(chan *Message),
// }

// func (b *broadcaster) Start() {
// 	for {
// 		select {
// 		case user := <-b.enterChannel:
// 			b.users[user.Name] = user
// 		case user := <-b.leaveChannel:
// 			delete(b.users, user.Name)
// 		case msg := <-b.messageChannel:
// 			for _, user := range b.users {
// 				user.MessageChan <- msg
// 			}
// 		}
// 	}
// }

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
		if data := ValidUser(r); len(data.UID) > 0 {
			tmpl := template.Must(template.ParseFiles("./views/chatroom.html"))
			tmpl.Execute(w, data)
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
		tmpl := template.Must(template.ParseFiles("./views/chatroom.html"))
		tmpl.Execute(w, nil)
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
func ValidUser(r *http.Request) User {
	chatAcc := Acc{}
	r.ParseForm()
	acc := r.FormValue("acc")
	pswd := r.FormValue("pswd")
	collection := MongoDBcontext("chat_db", "chat_acc")
	filter := bson.M{"acc": acc, "pswd": pswd}
	err := collection.FindOne(context.Background(), filter).Decode(&chatAcc)
	if err == nil {
		res := User{
			UID:  chatAcc.ID,
			Name: chatAcc.Name,
		}
		return res
	}
	return User{}
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

// GetUUID is
func GetUUID() string {
	uuid, err := uuid.New()

	if err != nil {
		log.Fatal(err)
	}
	var buf [36]byte
	hex.Encode(buf[0:8], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])

	return string(buf[:])
}

// Echo is
func Echo(ws *websocket.Conn) {
	ws.Request().ParseForm()
	pool := conns.PushBack(User{
		UID:  ws.Request().URL.Query().Get("clientId"),
		conn: ws})
	defer ws.Close()
	defer conns.Remove(pool)
	for {

		fmt.Println(ws.Request())
		var tmp string
		reply := User{}
		if err := websocket.Message.Receive(ws, &tmp); err != nil {
			fmt.Println("Can't receive, reason : " + err.Error())
			break
		}
		json.Unmarshal([]byte(tmp), &reply)

		switch reply.Message.Type {
		//enter
		case 0:
			Wellcome(conns, &reply)
			break
		//leave
		case 1:
			Leaving(conns, &reply)
			break
		//normal
		case 2:
			SendMessage(conns, &reply)
			break
		}
	}
}

// Wellcome is
func Wellcome(conns *list.List, reply *User) {
	for item := conns.Front(); item != nil; item = item.Next() {
		usr, ok := item.Value.(User)
		if !ok {
			panic("item not *websocket.Conn")
		}
		reply.Message.Content = "-----     wellcome " + reply.Name + " come in.     -----"
		if reply.UID == usr.UID {
			continue
		} else {
			websocket.Message.Send(usr.conn, reply)
		}
	}
}

// Leaving is
func Leaving(conns *list.List, reply *User) {
	for item := conns.Front(); item != nil; item = item.Next() {
		usr, ok := item.Value.(User)
		if !ok {
			panic("item not *websocket.Conn")
		}
		reply.Message.Content = "-----     " + reply.Name + " is leaved.     -----"
		if reply.UID == usr.UID {
			continue
		} else {
			websocket.Message.Send(usr.conn, reply)
		}
	}
}

// SendMessage is
func SendMessage(conns *list.List, reply *User) {
	for item := conns.Front(); item != nil; item = item.Next() {
		usr, ok := item.Value.(User)
		if !ok {
			panic("item not *websocket.Conn")
		}

		if reply.Message.ToID != "all" && reply.Message.ToID == usr.UID && len(reply.Message.ToID) > 0 {
			websocket.Message.Send(usr.conn, "This secret message from <font style='red'>"+reply.Name+"</font> say : "+reply.Message.Content)
		} else {
			websocket.Message.Send(usr.conn, reply.Name+" say : "+reply.Message.Content)
		}
	}
}

func main() {

	//go Broadcaster.Start()
	// page
	http.HandleFunc("/", Home)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/chatroom", ChatRoom)
	http.Handle("/ws", websocket.Handler(Echo))

	// static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// if any err log
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
