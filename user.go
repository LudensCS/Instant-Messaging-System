package main

import "net"

// 用户类
type User struct {
	Name string      //用户名称
	Addr string      //用户IP地址
	Ch   chan string //用户绑定的消息管道
	Conn net.Conn    //用户对应的客户端链接
}

// 创建一个用户API
func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		Ch:   make(chan string),
		Conn: conn,
	}
	//创建一个go程来监听当前User的Channel
	go user.ListenMessage()
	return user
}

// 监听当前User管道的方法,一旦有Message,立即发送给对应客户端
func (this *User) ListenMessage() {
	for msg := range this.Ch {
		this.Conn.Write([]byte(msg + "\n"))
	}
}
