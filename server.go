/*
	This file includes server module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import (
	"fmt"
	"net"
	"sync"
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
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// Socket handler
func (this *Server) Handler(conn net.Conn) {
	// user checkin, join in online map
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// Broad cast user login message
	this.BroadCast(user, "Logged in")

	// temp block
	select {}
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
