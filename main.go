package main

import (
	"container/list"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"golang.org/x/net/websocket"
)

// ChatData is
type ChatData struct {
	ClientID string `json:"client_id"`
	Name     string `json:"usr_name"`
	Msg      string `json:"msg"`
	ToID     string `json:"to_id"`
}

// ChatAcc is
type ChatAcc struct {
	Acc    string `bson:"acc"`
	Email  string `bson:"email"`
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
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
		if data := ValidUser(r); len(data.ClientID) > 0 {
			tmpl := template.Must(template.ParseFiles("./views/chatroom.html"))
			tmpl.Execute(w, data)
		} else {
			tmpl := template.Must(template.ParseFiles("./views/login.html"))
			tmpl.Execute(w, nil)
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

// ValidUser is checkout login user exist in mongodb
func ValidUser(r *http.Request) ChatData {
	chatAcc := ChatAcc{}
	r.ParseForm()
	acc := r.FormValue("acc")
	pswd := r.FormValue("pswd")
	collection := MongoDBcontext("chat_db", "chat_acc")
	filter := bson.M{"acc": acc, "pswd": pswd}
	err := collection.FindOne(context.Background(), filter).Decode(&chatAcc)

	if err == nil {
		res := ChatData{
			ClientID: GetUUID(),
			Name:     chatAcc.Name,
		}
		return res
	}
	return ChatData{}
}

// Register is user
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./views/register.html"))
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

var conns = list.New()

// Echo is
func Echo(ws *websocket.Conn) {
	pool := conns.PushBack(ws)
	defer ws.Close()
	defer conns.Remove(pool)

	for {
		var tmp string
		reply := ChatData{}
		if err := websocket.Message.Receive(ws, &tmp); err != nil {
			fmt.Println("Can't receive, reason : " + err.Error())
			break
		}
		json.Unmarshal([]byte(tmp), &reply)

		message := reply.Name + " say : " + reply.Msg

		for item := conns.Front(); item != nil; item = item.Next() {
			ws, ok := item.Value.(*websocket.Conn)
			if !ok {
				panic("item not *websocket.Conn")
			}
			if len(reply.ToID) > 0 && reply.ToID != "all" {
			} else {
				io.WriteString(ws, message)
			}
		}
	}
}

func main() {
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
