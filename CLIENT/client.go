package main

import (
	"fmt"
	"net"
)

// 客户端
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

// 创建一个Client接口
func NewClient(serverIp string, serverPort int) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Conn:       conn,
	}
	return client
}
func (this *Client) Handler() {

}
func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("<<<<<<<<链接服务器失败>>>>>>>>")
		return
	}
	fmt.Println("<<<<<<<<链接服务器成功>>>>>>>>")
	go client.Handler()
	for {
	}
}
