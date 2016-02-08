package main

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	// 受け取ったメッセージをすべてのクライアントに転送するために使われる
	forward chan []byte
}
