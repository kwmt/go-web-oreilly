package main
import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	// 受け取ったメッセージをすべてのクライアントに転送するために使われる
	forward chan []byte
	// チャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// チャットリームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// 在室しているすべてのクライアントが保持される
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := r.forward:
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージ送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}

	}
}


const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:socketBufferSize,
	WriteBufferSize:socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter ,req *http.Request) {
	// Websocketコネクションを取得
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket:socket,
		send:make(chan []byte, messageBufferSize),
		room:r,
	}
	// 現在のチャットルームに参加
	r.join <- client

	// クライアントの終了時、退室する
	defer func() { r.leave <- client}()
	// 書き込み待ちのgoroutine
	go client.write()
	// 読み込み待ちのgoroutine(main)
	client.read()
}

