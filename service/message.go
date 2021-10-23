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
func (m *Message) NewMessage() *Message {
	if m.ToUser == nil {
		return &Message{
			User:    m.User,
			Type:    MsgNormal,
			Content: m.User.UserInfo.Name + " says : " + m.Content,
			ToUser:  m.ToUser,
		}
	}
	return &Message{
		User:    m.User,
		Type:    MsgNormal,
		Content: m.Content,
		ToUser:  m.ToUser,
	}
}

// NewUserEnterMessage is
func (user *User) NewUserEnterMessage() *Message {
	return &Message{
		User:    user,
		Type:    MsgSentUserList,
		Content: user.UserInfo.Name + " 加入了聊天室",
	}
}

// NewUserLeaveMessage is
func (user *User) NewUserLeaveMessage() *Message {
	return &Message{
		User:    user,
		Type:    MsgSentUserList,
		Content: user.UserInfo.Name + " 離開了聊天室",
	}
}
