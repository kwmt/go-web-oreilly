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