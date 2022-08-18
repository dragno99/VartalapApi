package wSocket

type Room struct {
	// chatId
	ChatId string

	// collection of users(websocket connections) in this room
	Users map[*User]bool

	// broadcast channel to send message to all of the users in this room
	Broadcast chan Message

	// channel to register a new client
	Register chan *User

	// channel to register a new client
	Unregister chan *User

	// reason for defining register and unregister channel to avoid race condition
}

// create a new chat room
func NewRoom(chatId string) *Room {
	return &Room{
		Users:      make(map[*User]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *User),
		Unregister: make(chan *User),
		ChatId:     chatId,
	}
}

func (r *Room) Run() {
	for {
		select {

		case user := <-r.Register:
			{
				r.Users[user] = true
			}
		case user := <-r.Unregister:
			{
				user.Conn.Close()
				if _, ok := r.Users[user]; ok {
					delete(r.Users, user)
					close(user.Send)
				}
			}

		case msg := <-r.Broadcast:
			{
				for otherUser := range r.Users {
					// user should not send message to himself
					if otherUser.UserId == msg.SenderId.Hex() {
						continue
					}
					select {

					case otherUser.Send <- msg:

					default:
						{
							close(otherUser.Send)
							delete(r.Users, otherUser)
						}
					}
				}
			}
		}
	}
}
