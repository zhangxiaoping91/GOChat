package main

import (
	"ipc"
	"fmt"
	"cg"
	"strconv"
	"strings"
	"bufio"
	"os"
)

var centerClient *cg.CenterClient
/**
初始化服务
 */
func startCenterService() error {
	server := ipc.NewIpcServer(&cg.CenterServer{})
	client := ipc.NewIpcClient(server);
	centerClient = &cg.CenterClient{client}
	return nil
}

func Help(args []string) int {
	fmt.Println(`
	Commands:
		login <username><level><exp>
		logout <username>
		send <message>
		listplayer
		quit(q)
		help(h)
	`)
	return 1
}

/**
退出
 */
func Quit(args []string) int {
	return 2
}
/**
登出
 */
func Logout(args []string) int {
	if len(args) != 2 {
		fmt.Println("USAGE: logout <username>")
		return 0
	}
	centerClient.RemovePlayer(args[1])
	return 1
}

/**
登陆
 */
func Login(args []string) int {
	if len(args) != 4 {
		fmt.Println("USAGE:login <username><level><exp>")
		return 0
	}
	level, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid paramter:<level> should be an integer")
		return 0
	}
	exp, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("Invalid paramter:<exp> should be an integer")
		return 0
	}
	player := cg.NewPlayer()
	player.Name = args[1]
	player.Level = level
	player.Exp = exp
	err = centerClient.AddPlayer(player)
	if err != nil {
		fmt.Println("Failed adding player", err)
		return 0
	}
	return 1
}

/**
列表
 */
func ListPlayer(args []string) int {
	ps, err := centerClient.ListPlayer("")
	if err != nil {
		fmt.Println("Failed find listplayers", err)
		return 0
	} else {
		for i, v := range ps {
			fmt.Println(i + 1, ":", v)
		}
		return 1
	}
}

/**
发送消息
 */
func Send(args[] string) int {
	message := strings.Join(args[1:], " ")
	err := centerClient.BroadCast(message)
	if err != nil {
		fmt.Println("Failed BroadCast", err)
		return 0
	}
	return 1
}

func GetCommandHandel() map[string]func(args []string) int {
	return map[string]func(args []string) int{"help":Help, "h":Help, "quit":Quit, "q":Quit, "login":Login, "logout":Logout, "listplayer":ListPlayer, "send":Send, }
}

func main() {
	fmt.Println("Casual Game Server Solution")
	startCenterService()
	Help(nil)

	r := bufio.NewReader(os.Stdin)

	handles := GetCommandHandel()

	for {
		fmt.Println("Command> ")
		b, _, _ := r.ReadLine();
		line := string(b)
		tokens := strings.Split(line, " ")
		if handler, ok := handles[tokens[0]]; ok {
			ret := handler(tokens)
			if ret == 0 {
				fmt.Println(line,"命令有误")
				continue
			}else if ret ==1 {
				 fmt.Println(tokens,"执行成功")
				continue
			}else{
				break;
			}
		} else {
			fmt.Println("Unknown command", tokens[0])
		}
	}
}

func test() {
	server := ipc.NewIpcServer(&ipc.IpcServer{})

	client1 := ipc.NewIpcClient(server)
	client2 := ipc.NewIpcClient(server)

	resp1, _ := client1.Call("", "From Client1")
	resp2, _ := client2.Call("", "From Client2")
	fmt.Println(resp1)
	fmt.Println(resp2)

	client1.Close()
	client2.Close()

}
