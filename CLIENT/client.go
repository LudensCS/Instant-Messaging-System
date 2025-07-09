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
	flag       int //用户模式
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
		flag:       1,
	}
	return client
}

// 客户端菜单
func (this *Client) Menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scan(&flag)
	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("<<<<<<<<请输入合法范围内的数字>>>>>>>>")
		return false
	}
}

// 客户端业务
func (this *Client) Run() {
	for {
		for {
			if this.Menu() {
				break
			}
		}
		switch this.flag {
		case 0:
			return
		case 1:
			fmt.Println("公聊模式...")
		case 2:
			fmt.Println("私聊模式...")
		case 3:
			fmt.Println("更新用户名...")
		}
	}
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
	//启动客户端业务
	client.Run()
}
