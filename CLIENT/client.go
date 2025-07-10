package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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
		flag:       4869,
	}
	return client
}

// 客户端菜单
func (this *Client) Menu() bool {
	var flag int = 4869
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询当前在线用户")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 4 {
		this.flag = flag
		return true
	} else {
		fmt.Println("<<<<<<<<请输入合法范围内的数字>>>>>>>>")
		return false
	}
}

// 更新用户名
func (this *Client) UpdateName() bool {
	fmt.Println("<<<<<<<<请输入您的新用户名>>>>>>>>")
	//读取包含空格的字符串方法
	reader := bufio.NewReader(os.Stdin)
	msg, ERR := reader.ReadString('\n')
	if ERR != nil {
		fmt.Println("reader.ReadString error:", ERR)
		return false
	}
	msg = strings.TrimSuffix(msg, "\r\n")
	msg = strings.TrimSuffix(msg, "\n")
	this.Name = msg
	msg = "更改用户名@群助手@" + this.Name
	_, err := this.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Conn.Write error:", err)
		return false
	}
	return true
}

// 公聊模式
func (this *Client) PublicChat() {
	for {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("<<<<<<<<请输入聊天内容,exit表示退出当前模式>>>>>>>>")
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\r\n")
		msg = strings.TrimSuffix(msg, "\n")
		if err != nil {
			fmt.Println("reader.ReadString error:", err)
			return
		}
		if msg == "exit" {
			return
		}
		if len([]rune(msg)) > 0 {
			_, ERR := this.Conn.Write([]byte(msg))
			if ERR != nil {
				fmt.Println("Conn.Write error:", err)
				return
			}
		} else {
			fmt.Println("<<<<<<<<聊天内容不能为空>>>>>>>>")
		}
	}
}

// 查询当前在线用户
func (this *Client) OnlineUsers() {
	_, err := this.Conn.Write([]byte("查询当前在线用户@群助手\n"))
	if err != nil {
		fmt.Println("Conn.Write error:", err)
		return
	}
	time.Sleep(100 * time.Millisecond)
}

// 私聊模式
func (this *Client) PrivateChat() {
	for {
		this.OnlineUsers()
		fmt.Println("<<<<<<<<请输入聊天对象的用户名,exit表示退出当前模式>>>>>>>>")
		reader := bufio.NewReader(os.Stdin)
		receiver, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("reader.ReadString error:", err)
			return
		}
		receiver = strings.TrimSuffix(receiver, "\r\n")
		receiver = strings.TrimSuffix(receiver, "\n")
		if receiver == "exit" {
			return
		}
		if len([]rune(receiver)) > 0 {
			fmt.Println("<<<<<<<<请输入聊天内容>>>>>>>>")
			var msg string
			msg, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("reader.ReadString error:", err)
				return
			}
			this.Conn.Write([]byte("@" + receiver + "@" + msg))
		}
	}
}

// 打印服务端的回应消息,读写分离的思路
func (this *Client) DoResponse() {
	//这行代码相当于输出重定向,一旦连接内有数据就立刻输出,跟自己实现缓冲区读取后输出是等价的
	io.Copy(os.Stdout, this.Conn)
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
			this.PublicChat()
		case 2:
			this.PrivateChat()
		case 3:
			for this.UpdateName() == false {
			}
		case 4:
			this.OnlineUsers()
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
	go client.DoResponse()
	//启动客户端业务
	client.Run()
}
