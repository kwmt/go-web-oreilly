package main

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
