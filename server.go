/*
	This file includes server module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Server class datatypes
type Server struct {
	Ip   string //Server IPv4 address
	Port int    //Server Listening Port

	OnlineMap map[string]*User // Current user status map
	mapLock   sync.RWMutex     // RW lock for user status

	Message chan string // Channel for broadcast
}

// Constructor
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// Listening Message and forward to all users' channel
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast message
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + time.Now().Format("15:04:05") + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// Socket handler
func (this *Server) Handler(conn net.Conn) {
	// create user assests
	user := NewUser(conn, this)

	user.Online()

	// setup activity counter
	keepAlive := make(chan bool)

	// Accept user messages
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
				return
			}

			// rstrip message
			msg := string(buf[:n-1])
			// handle current user message
			user.HandleMessage(msg)
			// reset activity counter
			keepAlive <- true
		}
	}()

	// temp block
	for {
		select {
		case <-keepAlive:
			// NOOP
		case <-time.After(time.Minute * 15):
			// No activity timeout
			// force logout
			user.SendMsg("No activity for 15min, you are logged out")
			close(user.C)
			conn.Close()
			return
		}
	}
}

// Server startup
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// socket close
	defer listener.Close()
	// start listening message
	go this.ListenMessage()
	// loop begin
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// handler
		go this.Handler(conn)
	}
}
