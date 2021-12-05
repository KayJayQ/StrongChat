/*
	This file includes server module and related functions

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import (
	"fmt"
	"net"
)

// Server class datatypes
type Server struct {
	Ip   string //Server IPv4 address
	Port int    //Server Listening Port
}

// Constructor
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// Socket handler
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection established")
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
