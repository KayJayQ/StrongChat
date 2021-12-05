/*
	This file includes user module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import "net"

// User class datatypes
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// User constructor
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// start listening user's message
	go user.ListenMessage()

	return user
}

// Listen user's channel and forward to user connection
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
