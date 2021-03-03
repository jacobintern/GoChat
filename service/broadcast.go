package service

// BroadcasterModel is
type BroadcasterModel struct {
	users map[string]*User

	enterChannel        chan *User
	leaveChannel        chan *User
	messageChannel      chan *Message
	requestUsersChannel chan struct{}
	usersChannel        chan []*UserInfo
}

// Broadcaster is
var Broadcaster = &BroadcasterModel{
	users:               make(map[string]*User),
	enterChannel:        make(chan *User),
	leaveChannel:        make(chan *User),
	messageChannel:      make(chan *Message),
	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*UserInfo),
}

// Start is
func (b *BroadcasterModel) Start() {
	for {
		select {
		case user := <-b.enterChannel:
			b.users[user.UserInfo.Name] = user
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				user.MessageChannel <- msg
			}
		case user := <-b.leaveChannel:
			delete(b.users, user.UserInfo.Name)
			close(user.MessageChannel)
		case <-b.requestUsersChannel:
			userList := make([]*UserInfo, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user.UserInfo)
			}
			b.usersChannel <- userList
		}
	}
}

// UserEntering is
func (b *BroadcasterModel) UserEntering(u *User) {
	b.enterChannel <- u
}

// UserLeaving is
func (b *BroadcasterModel) UserLeaving(u *User) {
	b.leaveChannel <- u
}

// Broadcast is
func (b *BroadcasterModel) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

// GetUserList is
func (b *BroadcasterModel) GetUserList() []*UserInfo {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
