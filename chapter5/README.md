### 5.2.1 NSQ

インストール
参考 http://tleyden.github.io/blog/2014/11/12/an-example-of-using-nsq-from-go/


```bash
$ wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-0.3.6.linux-amd64.go1.5.1.tar.gz
$ tar xvfz nsq-0.3.6.linux-amd64.go1.5.1.tar.gz 
$ sudo mv nsq-0.3.6.linux-amd64.go1.5.1/bin/* /usr/local/bin
```

起動確認

```bash
# nsqlookupd
[nsqlookupd] 2016/02/11 15:15:03.295271 nsqlookupd v0.3.6 (built w/go1.5.1)
[nsqlookupd] 2016/02/11 15:15:03.295438 HTTP: listening on [::]:4161
[nsqlookupd] 2016/02/11 15:15:03.295915 TCP: listening on [::]:4160
```

go-nsqをインストール

```bash
$ go get -v github.com/bitly/go-nsq
```

### 5.2.2 Mongo DB

```
$ wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian71-3.2.1.tgz
$ tar xvzf mongodb-linux-x86_64-debian71-3.2.1.tgz 
$ mv mongodb-linux-x86_64-debian71-3.2.1/bin/* /usr/local/bin/
```

起動確認

```bash
$ mongod --dbpath /tmp/
```

MongoDBドライバーインストール

```bash
$ go get gopkg.in/mgo.v2
```

###  異なるターミナルでそれぞれ実行する

* nsqlooupdを起動しnsqdインスタンスを発見できるようにする

```bash
$ nsqlookupd
```

* nsqdを起動し、どのnsqlookupdを利用するか指定する

```bash
$ nsqd --lookupd-tcp-address=localhost:4160
```

* mongodを起動してデータ関連のサービスを実行する

```bash
$ mongod --dbpath ./db
```

--dbpathはデータの保管先を指定する。存在しないディレクトリは指定できない。

またdockerなどを使っている場合は、ホストとvolumeを共有している場合は、失敗する。共有してないディレクトリを指定すると良さそう
参考:MongoDB not working in Vagrant Centos Box - Stack Overflow http://stackoverflow.com/questions/35005560/mongodb-not-working-in-vagrant-centos-box

```
2016-02-11T16:00:43.055+0000 E STORAGE  [initandlisten] WiredTiger (22) [1455206443:55491][1063:0x7f9dbe501cc0], connection: ./db/WiredTiger.wt: fsync: Invalid argument
2016-02-11T16:00:43.057+0000 I -        [initandlisten] Fatal Assertion 28561
2016-02-11T16:00:43.057+0000 I -        [initandlisten] 

***aborting after fassert() failure
```


### 5.3.1 Twitterを使った認証

https://apps.twitter.com



