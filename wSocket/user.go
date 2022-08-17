package wSocket

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dragno99/VartalapApi/database"
	"github.com/dragno99/VartalapApi/model"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	// UserId
	UserId string

	// room in which user resides
	Room *Room

	// websocket connection for the user
	Conn *websocket.Conn

	// send channel for the user to recieve messages from the websocket
	Send chan Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// readPump : pumps message from the websocket connection to the room
func (u *User) readPump() {

	for {
		var msg Message
		err := u.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			u.Room.Unregister <- u
			break
		}
		u.Room.Broadcast <- msg
	}
}

// writePump : pumps message from the hub to the websocket connetion
func (u *User) writePump() {

	for {
		msg, ok := <-u.Send
		if !ok {
			// the room closed the channel
			u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			u.Room.Unregister <- u
			break
		}
		// write message to websocket connection

		var wg sync.WaitGroup

		wg.Add(1)

		go func() {
			for !pushMessageIntoDatabse(u.Room.ChatId, msg) {
				break
			}
			wg.Done()
		}()

		// wait for pushing message into database
		wg.Wait()
		err := u.Conn.WriteJSON(msg)

		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			u.Room.Unregister <- u
			break
		}
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

// serveWS handles websocket request from the peer
func ServeWS(room *Room, userId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	user := &User{UserId: userId, Room: room, Conn: conn, Send: make(chan Message)}

	room.Register <- user

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go user.writePump()
	go user.readPump()

}

func pushMessageIntoDatabse(chatId string, msg Message) bool {
	currMessage := model.Message{
		Text:      msg.Text,
		Image:     msg.Image,
		SenderId:  msg.SenderId,
		CreatedAt: msg.CreatedAt,
	}

	update := bson.M{
		"$push": bson.M{
			"messages": currMessage,
		},
	}
	chatIdObj, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		return false
	}
	result, err := database.ChatCollection.UpdateByID(context.TODO(), chatIdObj, update)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if result.MatchedCount > 0 {
		return true
	} else {
		return false
	}
}
