package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートのコンパイルは1回だけで良いため、
	// 複数のgoroutineがServeHTTPを呼び出しても、引数として渡した関数が1回しか実行されないこと保証する。
	// また、本当に必要になるまで処理を後回しにする遅延初期化の役割ももつ。
	// ただし、エラーが発生しうる処理の場合、エラーに気づきにくくなる問題もあるので、
	// `template.Must`を使ってグローバル変数に初期化時にセットする方が好まれる場合もある。
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// HTMLを出力する際にhttp.Requestに含まれるデータを参照できるようにrを渡す
	log.Println("Host:", r.Host)
	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈

	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// チャットルームを開始
	go r.run()

	// Webサーバーを起動
	log.Println("Webサーバーを開始します。port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
