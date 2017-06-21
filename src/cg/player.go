package cg

import "fmt"

type Player struct {
	Name  string "name"
	Level int "Level"
	Exp   int "Exp"
	Room  string "Room"

	mq    chan *Message
}

func NewPlayer() *Player {
	m := make(chan *Message, 1024)
	player := &Player{"", 0, 0,"", m}

	go func(p *Player) {
		for {
			msg := <-p.mq
			fmt.Println(p.Name, "Received message:", msg.Content)
		}
	}(player)
	return player
}
