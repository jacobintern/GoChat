package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/jacobintern/GoChat/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"
)

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
		if uid := service.ValidUser(r); len(uid) > 0 {
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
			data := service.GetUser(uid)
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

// CreateUser is
func CreateUser(r *http.Request) *mongo.InsertOneResult {
	r.ParseForm()
	collection := service.MongoDBcontext("chat_db", "chat_acc")
	res, err := collection.InsertOne(context.Background(), service.Acc{
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

// Echo is
func Echo(conn *websocket.Conn) {
	// 建立使用者
	user := service.NewUser(conn)
	// 建立傳送訊息通道 goroutine監聽
	go user.SendMessage()

	// 使用者進入
	msg := service.NewUserEnterMessage(user)
	service.Broadcaster.UserEntering(user)
	service.Broadcaster.Broadcast(msg)

	// 訊息接收並傳送給其他使用者
	err := user.ReceiveMessage()

	// 使用者離開
	msg = service.NewUserLeaveMessage(user)
	service.Broadcaster.UserLeaving(user)
	service.Broadcaster.Broadcast(msg)

	if err == nil {
		conn.Close()
	} else {
		log.Println("read from client error:", err)
		conn.Close()
	}

}

// RegisterchatHandler is
func RegisterchatHandler() {
	go service.Broadcaster.Start()

	http.Handle("/ws", websocket.Handler(Echo))
}

// GetUsers is
func GetUsers(w http.ResponseWriter, req *http.Request) {
	userList := service.Broadcaster.GetUserList()
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
