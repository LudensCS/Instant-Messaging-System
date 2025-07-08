package main

import (
	"net"
	"time"
)

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
	defer this.Svr.MapLock.Unlock()
	if _, ok := this.Svr.OnlineMap[this.Name]; !ok {
		return
	}
	delete(this.Svr.OnlineMap, this.Name)
	//广播当前用户下线消息
	this.Svr.BroadCast(this, "下线")
	time.Sleep(100 * time.Millisecond)
	close(this.Ch)
	this.Conn.Close()
}

// 给当前user对应客户端发送消息
func (this *User) SendMessage(msg string) {
	if this.Ch != nil {
		this.Ch <- msg
	}
}

// 处理用户消息业务
func (this *User) DoMessage(msg string) {
	if msg == "查询当前在线用户@群助手" {
		this.Svr.MapLock.RLock()
		for _, user := range this.Svr.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "在线..."
			this.SendMessage(onlineMsg)
		}
		this.Svr.MapLock.RUnlock()
	} else if len([]rune(msg)) > 10 && string([]rune(msg)[0:10]) == "更改用户名@群助手@" {
		//带中文的字符串一定要先转成[]rune类型再做各类操作!!!
		NewName := string([]rune(msg)[10:])
		this.Svr.MapLock.Lock()
		defer this.Svr.MapLock.Unlock()
		_, ok := this.Svr.OnlineMap[NewName]
		if ok {
			this.SendMessage("您更改的用户名已存在")
			return
		}
		delete(this.Svr.OnlineMap, this.Name)
		this.Svr.OnlineMap[NewName] = this
		this.Name = NewName
		this.SendMessage("用户名更改成功!")
	} else {
		//广播消息
		this.Svr.BroadCast(this, msg)
	}
}

// 监听当前User管道的方法,一旦有Message,立即发送给对应客户端
func (this *User) ListenMessage() {
	for msg := range this.Ch {
		this.Conn.Write([]byte(msg + "\n"))
	}
}
