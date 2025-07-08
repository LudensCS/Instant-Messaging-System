package main

import (
	"flag"
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

var serverIp string
var serverPort int

func init() {
	//命令行解析设置
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置连接的服务器地址")
	flag.IntVar(&serverPort, "port", 8888, "设置连接的服务器端口")
}
func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("<<<<<<<<链接服务器失败>>>>>>>>")
		return
	}
	fmt.Println("<<<<<<<<链接服务器成功>>>>>>>>")
	go client.Handler()
	for {
	}
}
