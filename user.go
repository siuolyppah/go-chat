package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	// user will be listening for incoming messages on this channel.
	ReadChan chan string

	// in server view
	conn   net.Conn
	server *Server
}

func newUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:     userAddr,
		Addr:     userAddr,
		ReadChan: make(chan string),
		conn:     conn,
		server:   server,
	}

	go user.ListenComingMessage()

	return user
}

func (this *User) ListenComingMessage() {
	for {
		msg := <-this.ReadChan

		_, err := this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Conn.Write err: ", err)
		}
	}
}

func (this *User) Login() {
	this.server.mapLock.Lock()
	this.server.UserOnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "Login")
}

func (this *User) Logout() {
	this.server.mapLock.Lock()
	delete(this.server.UserOnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "Logout")
}

// SendMsg send msg to `this`
func (this *User) SendMsg(msg string) {
	_, err := this.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("Conn.Write err: ", err)
	}
}

func (this *User) MsgHandle(msg string) {

	if msg == "/cmd online-users" {
		this.server.mapLock.RLock()
		for _, u := range this.server.UserOnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s : Online", u.Addr, u.Name)
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.RUnlock()
	} else {
		this.server.Broadcast(this, msg)
	}

}
