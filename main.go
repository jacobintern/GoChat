package main

import (
	"container/list"
	"context"
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
	Msg      string `json:"msg"`
	ToID     string `json:"to_id"`
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
		r.ParseForm()
		acc := r.FormValue("acc")
		pswd := r.FormValue("pswd")
		if ValidUser(acc, pswd) {
			http.Redirect(w, r, "/chatroom", http.StatusSeeOther)
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

// ValidUser is checkout login user exist in mongodb
func ValidUser(acc string, pswd string) bool {
	collection := MongoDBcontext("chat_db", "chat_acc")
	filter := bson.M{"acc": acc, "pswd": pswd}
	// data, err := collection.Find(context.Background(), filter)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	data := collection.FindOne(context.Background(), filter)
	return data.Err() != mongo.ErrNoDocuments
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./views/chatroom.html"))
		data := ChatData{
			ClientID: GetUUID(),
		}
		fmt.Printf(data.ClientID)
		tmpl.Execute(w, data)
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
	return string(uuid[:])
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
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			fmt.Println(err.Error())
			break
		}
		json.Unmarshal([]byte(tmp), &reply)

		fmt.Println(reply.ClientID)
		fmt.Println(reply.ToID)
		fmt.Println(reply.Msg)

		for item := conns.Front(); item != nil; item = item.Next() {
			ws, ok := item.Value.(*websocket.Conn)
			if !ok {
				panic("item not *websocket.Conn")
			}
			io.WriteString(ws, reply.Msg)
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
