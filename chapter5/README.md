NSQのインストール
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