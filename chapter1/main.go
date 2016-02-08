package main

import (
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
	t.templ.Execute(w, nil)
}

func main() {

	http.Handle("/", &templateHandler{filename: "chat.html"})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
