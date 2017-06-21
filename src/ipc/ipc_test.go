package ipc

import (
	"testing"
	"fmt"
)

type EchoServer struct {

}

func (server *EchoServer) Handle(method,request string) *Response{
      return &Response{}
}

func (server *EchoServer) Name() string{
	return "EchoServer"
}

func TestIpc(t *testing.T) {
	server:=NewIpcServer(&EchoServer{})

	client1:=NewIpcClient(server)
	client2:=NewIpcClient(server)

	resp1,_ :=client1.Call("","From Client1")
	resp2,_:=client2.Call("","From Client2")

	fmt.Println(resp1)
	fmt.Println(resp2)

	client1.Close()
	client2.Close()

}
