/*
	This file includes user module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import (
	"net"
	"strings"
)

// User class datatypes
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// User constructor
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// start listening user's message
	go user.ListenMessage()

	return user
}

// user online function
func (this *User) Online() {
	// user checkin, join in online map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// Broad cast user login message
	this.server.BroadCast(this, "Logged in")
}

// user offline function
func (this *User) Offline() {
	// user checkout, join in online map
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// Broad cast user logout message
	this.server.BroadCast(this, "Logged out")
}

// Send message to certain user
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// User message handler
func (this *User) HandleMessage(msg string) {
	// parse message to handle user cmds
	params := strings.Split(msg, "?")
	if len(params) == 0 {
		return
	}

	// if found cmd in api map, execute, or broadcast message
	if f, found := api_map[params[0]]; found {
		f(this, this.server, params)
	} else {
		this.server.BroadCast(this, msg)
	}
}

// Listen user's channel and forward to user connection
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
