package cg

import (
	"ipc"
	"sync"
	"encoding/json"
	"errors"
)

type Message struct {
	From    string "From"
	To      string "To"
	Content string "Content"
}

type Room struct {

}

type CenterServer struct {
	servers map[string]ipc.Server
	players []*Player
	rooms   []*Room
	mutex   sync.RWMutex
}

func NewCenterServer() *CenterServer {
	servers := make(map[string]ipc.Server)
	players := make([]*Player, 0)
	return &CenterServer{servers:servers, players:players}
}

/**
添加游戏选手
 */
func (server *CenterServer) addPlayer(param string) error {
	player := NewPlayer()
	err := json.Unmarshal([]byte(param), &player)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()
	server.players = append(server.players, player)
	return nil
}

/**
移除选手
 */
func (server *CenterServer) removePlayer(param string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	for i, v := range server.players {
		if v.Name == param {
			if len(server.players) == 1 {
				server.players = make([]*Player, 0)
			} else if i == len(server.players) - 1 {
				server.players = server.players[:i - 1]
			} else if i == 0 {
				server.players = server.players[1:]
			} else {
				server.players = append(server.players[:i - 1], server.players[i + 1:]...)
			}
			return nil
		}
		return errors.New("Player is not found")
	}
	return nil
}

/**
返回列表
 */
func (server *CenterServer) listPlayer(param string) (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		b, _ := json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("No player online")
	}
	return
}

/**
广播
 */
func (server *CenterServer) broadcast(param string) error {
	var message Message
	err := json.Unmarshal([]byte(param), &message)
	if err != nil {
		return nil
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()
	if len(server.players) > 0 {
		for _, player := range server.players {
			player.mq <- &message
		}
	} else {
		err = errors.New("No player online")
	}
	return err
}

/**
处理
 */
func (server *CenterServer) Handle(method,param string) *ipc.Response{
	switch method {
	case "addPlayer":
		err:=server.addPlayer(param)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
	case "removePlayer":
		err:=server.removePlayer(param)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
	case "listPlayer":
		player,err:=server.listPlayer(param)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{"200",player}
	case "broadcast":
		err:=server.broadcast(param)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{Code:"200"}
	default:
		return &ipc.Response{Code:"404",Body:method+":"+param}
	}
	return &ipc.Response{Code:"200"}
}

func (server *CenterServer) Name() string{
	return "CenterServer"
}
