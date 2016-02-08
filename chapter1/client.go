package main

import "github.com/gorilla/websocket"

type client struct {
	// クライアントのためのWebsocket
	socket *websocket.Conn
	// メッセージを送るチャネル
	send chan []byte
	// クライアントが参加しているチャットルーム
	room *room
}

// websocketからデータを読み込む
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// websocketに書き出す
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
