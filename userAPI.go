/*
	This file includes user interfaces towards server.
	API list:
		Get current online users (name)
		Set current user name
		Send direct message to another user

	Author: Kay Qiang
	Email : qiangkj@cmu.edu
*/
package main

import "time"

var api_map map[string]func(user *User, server *Server, param []string)

// Put all API functions to api map
func ApiInitialize() {
	api_map = make(map[string]func(user *User, server *Server, param []string))
	api_map["!LISTUSERS"] = GetOnlineUsers
	api_map["!CHANGEUSERNAME"] = ChangeUserName
	api_map["!TO"] = DirectMessage
}

// Receive all online users' name from server
func GetOnlineUsers(user *User, server *Server, param []string) {
	server.mapLock.RLock()
	for _, userClass := range server.OnlineMap {
		onlineMsg := "[" + time.Now().Format("15:04:05") + "]" + userClass.Name + ":" + "is online...\n"
		user.SendMsg(onlineMsg)
	}
	server.mapLock.RUnlock()
}

// Change current user name
func ChangeUserName(user *User, server *Server, param []string) {
	result := false
	newName := param[1]
	server.mapLock.Lock()
	if _, ok := server.OnlineMap[newName]; ok {
		result = false
	} else {
		delete(server.OnlineMap, user.Name)
		server.OnlineMap[newName] = user
		user.Name = newName
		result = true
	}
	server.mapLock.Unlock()
	if result {
		user.SendMsg("[" + time.Now().Format("15:04:05") + "]" + user.Name + ":" + "User name changed")
	} else {
		user.SendMsg("[" + time.Now().Format("15:04:05") + "]" + user.Name + ":" + "User name change failed: existing user name")
	}
}

// Send direct message
func DirectMessage(user *User, server *Server, param []string) {
	if len(param) < 3 {
		return
	}
	tgtName := param[1]
	msg := param[2]

	result := false

	server.mapLock.Lock()
	tgtUser, ok := server.OnlineMap[tgtName]
	if ok {
		result = true
	} else {
		result = false
	}
	server.mapLock.Unlock()
	if result {
		tgtUser.SendMsg("[" + time.Now().Format("15:04:05") + "]" + user.Name + "|DM:" + msg)
		user.SendMsg("RET: Message sent to " + tgtName)
	} else {
		user.SendMsg("ERR: User Not Found")
	}
}
