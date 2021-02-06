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

// ChatData is
type ChatData struct {
	ClientID string `json:"client_id"`
	Name     string `json:"usr_name"`
	Msg      string `json:"msg"`
	ToID     string `json:"to_id"`
}

// ChatAcc is
type ChatAcc struct {
	ID     string `bson:"_id,omitempty"`
	Acc    string `bson:"acc"`
	Pswd   string `bson:"pswd"`
	Email  string `bson:"email"`
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
}

// UsrConn is
type UsrConn struct {
	ConID string
	conn  *websocket.Conn
}

var conns = list.New()

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
			ClientID: chatAcc.ID,
			Name:     chatAcc.Name,
		}
		return res
	}
	return ChatData{}
}

// CreateUser is
func CreateUser(r *http.Request) *mongo.InsertOneResult {
	r.ParseForm()
	collection := MongoDBcontext("chat_db", "chat_acc")
	res, err := collection.InsertOne(context.Background(), ChatAcc{
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
	pool := conns.PushBack(UsrConn{
		ConID: ws.Request().URL.Query().Get("clientId"),
		conn:  ws})
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
			usr, ok := item.Value.(UsrConn)
			if !ok {
				panic("item not *websocket.Conn")
			}

			if len(reply.ToID) > 0 && reply.ToID != "all" && reply.ToID == usr.ConID {
			} else {
				websocket.Message.Send(usr.conn, message)
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
