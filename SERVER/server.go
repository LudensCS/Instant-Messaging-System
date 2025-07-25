package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// 服务器
type Server struct {
	Ip   string //IP地址
	Port int    //端口
	//在线用户列表
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	//广播消息管道
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		MapLock:   sync.RWMutex{},
		Message:   make(chan string),
	}
	return server
}

// 监听Message管道,一旦有消息立马广播给所有user
func (this *Server) ListenMessage() {
	for msg := range this.Message {
		//给所有user广播msg
		this.MapLock.RLock()
		for _, user := range this.OnlineMap {
			user.Ch <- msg
		}
		this.MapLock.RUnlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	SendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- SendMsg
}

// 处理当前链接的业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
	//创建新用户
	user := NewUser(conn, this)
	user.Online()
	//可允许的最长不活跃时长
	const TIMEOUT = time.Minute * 5
	//计时器,监测用户活动状况
	tick := time.NewTimer(TIMEOUT)
	defer tick.Stop()
	//接受用户发送的消息并广播
	go func() {
		buf := make([]byte, 4096) //缓冲区
		for {
			//从conn中读取消息到buf,n是成功读取的字节数
			n, err := conn.Read(buf)
			//用户下线
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read error:", err)
				user.Offline()
				return
			}
			tick.Reset(TIMEOUT)
			//去除末尾换行
			msg := strings.TrimSuffix(string(buf[0:n]), "\r\n")
			msg = strings.TrimSuffix(msg, "\n")
			//用户针对message进行处理
			user.DoMessage(msg)
		}
	}()
	//超时强踢功能
WAIT:
	for {
		select {
		case <-tick.C:
			user.SendMessage("您长时间未活动,已被强制踢出")
			time.Sleep(100 * time.Millisecond)
			user.Offline()
			break WAIT
		}
	}
}

// 启动Server的接口
func (this *Server) Start() {
	//接收关服指令
	go func() {
		var command string
		for {
			fmt.Scanln(&command)
			if command == "exit" {
				os.Exit(0)
			}
		}
	}()
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	//close listen socket
	defer listener.Close()
	//启动一个go程监听Message管道
	go this.ListenMessage()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
}
