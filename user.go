/*
	This file includes user module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import "net"

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

// User message handler
func (this *User) HandleMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// Listen user's channel and forward to user connection
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
