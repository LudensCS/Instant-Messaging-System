package main

import "net"

// 用户类
type User struct {
	Name string      //用户名称
	Addr string      //用户IP地址
	Ch   chan string //用户绑定的消息管道
	Conn net.Conn    //用户对应的客户端链接
	Svr  *Server     //用户对应的服务端
}

// 创建一个用户API
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		Ch:   make(chan string),
		Conn: conn,
		Svr:  server,
	}
	//创建一个go程来监听当前User的Channel
	go user.ListenMessage()
	return user
}

// 用户上线业务
func (this *User) Online() {
	//用户上线,将用户加入在线用户列表
	this.Svr.MapLock.Lock()
	this.Svr.OnlineMap[this.Name] = this
	this.Svr.MapLock.Unlock()
	//广播当前用户上线消息
	this.Svr.BroadCast(this, "已上线")
}

// 用户下线业务
func (this *User) Offline() {
	//用户下线,将用户从在线用户列表中删除
	this.Svr.MapLock.Lock()
	delete(this.Svr.OnlineMap, this.Name)
	this.Svr.MapLock.Unlock()
	//广播当前用户下线消息
	this.Svr.BroadCast(this, "下线")
}

// 处理用户消息业务
func (this *User) DoMessage(msg string) {
	//广播消息
	this.Svr.BroadCast(this, msg)
}

// 监听当前User管道的方法,一旦有Message,立即发送给对应客户端
func (this *User) ListenMessage() {
	for msg := range this.Ch {
		this.Conn.Write([]byte(msg + "\n"))
	}
}
