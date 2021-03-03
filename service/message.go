package service

// Message type
const (
	MsgNormal = iota
	MsgSystem
	MsgSentUserList
)

// Message is
type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	ToUser  *UserInfo `json:"to_user"`
}

// NewMessage is
func NewMessage(user *User, content string, toUser *UserInfo) *Message {
	if toUser == nil {
		return &Message{
			User:    user,
			Type:    MsgNormal,
			Content: user.UserInfo.Name + " says : " + content,
			ToUser:  toUser,
		}
	}
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
