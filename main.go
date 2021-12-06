/*
	This file is the entrance of the server program

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import "fmt"

func main() {
	fmt.Println("Program Started...")
	ApiInitialize()
	server := NewServer("127.0.0.1", 8080)
	server.Start()
}
