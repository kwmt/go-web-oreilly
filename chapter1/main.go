package main
import (
	"net/http"
	"log"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<html>
			<head>
				<title>チャット</title>
			</head>
			<body>
				チャットしましょう！
			</body>
		</html>
		`))
	})


	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
