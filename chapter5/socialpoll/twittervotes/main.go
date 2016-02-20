/*
* mgoを使ってMongoDBデータベースからすべての投票結果を読み込み、
それぞれのドキュメントに含まれる配列optionsからすべての選択肢を取り出します。

* TwitterストリーミングAPIを使い、接続の開始や管理を行うとともに各選択肢に言及したツイートを検索します

* 検索にマッチしたツイートについて、選択肢の文字列をNSQに送信します。

*

 */


package main
import (
	"gopkg.in/mgo.v2"
	"log"
	"net/url"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/bitly/go-nsq"
)

func main() {

}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDBにダイアル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("データベースが閉じられました")
}


type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	// イテレータで順次アクセス。pollオブジェクトは1つしか使われないから、メモリ使用量を少なくできる
	// Allメソッドを使うと、投票の数に比例するのでとんでもないメモリの使用量になるかも。
	// optionsスライスも巨大になることがもちろんある。
	// このような場合にスケーラビリティを向上させるには、たとえば、複数のtwittervoteプログラムを起動して、処理を分担するといった対策が考えられる。
	for iter.Next(&p) {
		options = append(options, p.Options)
	}
	iter.Close()
	return options, iter.Err()
}

type tweet struct {
	Text string
}

// votesチャネルは送信専用
func readFromTwitter(votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("選択肢の読み込みに失敗しました:", err)
		return
	}

	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", "https://stream.twitter.com/1.1/statuses/filter.json",
		strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("検索リクエストの作成に失敗しました:", err)
		return
	}
	resp, err := makeRequest(req, query)
	if err != nil {
		log.Println("検索のリクエストに失敗しました:", err)
		return
	}
	reader = resp.Body
	decoder := json.NewDecoder(reader)

	for {
		var tweet tweet
		// 主に接続が閉じられたなどの理由でエラーがが発生したら、ループから抜けだして、呼び出し元に戻る
		if err := decoder.Decode(&tweet); err != nil {
			break
		}
		// すべての選択肢に対して
		for _, option := range options {
			// ツイートの中で言及されてる場合にはvotesチャネルにその選択肢を送信する。
			// 複数の選択肢を投票できる
			if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(option)) {
				log.Print("投票:", option)
				votes <- option
			}
		}
	}


}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub,_ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) //投稿内容をパブリッシュします
		}
		log.Println("Publisher: 停止中です")
		pub.Stop()
		log.Println("Publisher: 停止しました")
		stopchan <- struct{}{}
	}()
	return stopchan
}