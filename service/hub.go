package service

// HubModel is
type HubModel struct {
	users map[string]*User

	enterChannel        chan *User
	leaveChannel        chan *User
	messageChannel      chan *Message
	requestUsersChannel chan struct{}
	usersChannel        chan []*UserInfo
}

// Hub is
var Hub = &HubModel{
	users:               make(map[string]*User),
	enterChannel:        make(chan *User),
	leaveChannel:        make(chan *User),
	messageChannel:      make(chan *Message),
	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*UserInfo),
}

// Run is
func (h *HubModel) Run() {
	for {
		select {
		case user := <-h.enterChannel:
			h.users[user.UserInfo.Name] = user
		case msg := <-h.messageChannel:
			for _, user := range h.users {
				user.MessageChannel <- msg
			}
		case user := <-h.leaveChannel:
			delete(h.users, user.UserInfo.Name)
			close(user.MessageChannel)
		case <-h.requestUsersChannel:
			userList := make([]*UserInfo, 0, len(h.users))
			for _, user := range h.users {
				userList = append(userList, user.UserInfo)
			}
			h.usersChannel <- userList
		}
	}
}

// UserEntering is
func (b *HubModel) UserEntering(u *User) {
	b.enterChannel <- u
}

// UserLeaving is
func (b *HubModel) UserLeaving(u *User) {
	b.leaveChannel <- u
}

// Broadcast is
func (b *HubModel) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

// GetUserList is
func (b *HubModel) GetUserList() []*UserInfo {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
