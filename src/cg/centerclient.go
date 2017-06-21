package cg

import (
	"ipc"
	"encoding/json"
	"errors"
)

type CenterClient struct {
	*ipc.IpcClient
}

/**
添加选手
 */
func (client *CenterClient) AddPlayer(player *Player) error {
	b, err := json.Marshal(*player)
	if err != nil {
		return err
	}
	resp, err := client.Call("addPlayer", string(b))
	if err == nil&&resp.Code == "200" {
		return nil
	}
	return err
}

/**
移除选手
 */
func (client *CenterClient)  RemovePlayer(name string) error {
	ret, _ := client.Call("removePlayer", name)
	if ret.Code == "200" {
		return nil
	}
	return errors.New(ret.Code)
}

/**
列表
 */
func (client *CenterClient) ListPlayer(param string) (ps []*Player, err error) {
	resp, _ := client.Call("listPlayer", param)
	if resp.Code != "200" {
		err = errors.New(resp.Code)
		return
	}
	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

/**
广播
 */
func (client *CenterClient) BroadCast(message string) error{
	m:=&Message{Content:message}
	b,err:=json.Marshal(m);
	if err!=nil{
		return err;
	}
	resp,_:= client.Call("broadcast",string(b))
	if resp.Code=="200"{
		return nil
	}
	return errors.New(resp.Code)
}
