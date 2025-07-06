package main

import (
	"fmt"
	"net"
)

// 服务器
type Server struct {
	Ip   string //IP地址
	Port int    //端口
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// 处理当前链接的业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
}

// 启动Server的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	//close listen socket
	defer listener.Close()
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
