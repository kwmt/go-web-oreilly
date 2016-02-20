package main
import (
	"net"
	"time"
	"io"
	"log"
	"sync"
	"net/http"
	"net/url"
	"strconv"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
)

var conn net.Conn


// ネットワーク(netw)上のアドレス(addr)に接続する。
// すでに接続済みなら接続を閉じて再接続する。
func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}


var (
	authClient *oauth.Client
	creds *oauth.Credentials
)

func setupTwitterAuth() {
	var ts struct{
		ConsumerKey string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken string `env:"SP_TWITTER_ACCESSTOKEN,required"`
		AccessSecret string `env:"SP_TWITTER_ACCESSSECRET,required"`
	}

	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth.Credentials{
		Token:ts.AccessToken,
		Secret:ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials:oauth.Credentials{
			Token:ts.ConsumerKey,
			Secret:ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	// makeRequest関数が何回呼び出されても初期化のコードは1回しか実行されないようにしている
	authSetupOnce.Do(func(){
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport:&http.Transport{
				Dial: dial,
			},
		}
	})

	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	authClient.SetAuthorizationHeader(req, creds, "POST", req.URL, params)
	return httpClient.Do(req)


}

// stopchan:受信専用 シグナルのチャネル goroutineの終了を指示(
// votes:投票内容が送信されるチャネル
// 戻り:シグナルのチャネル goroutineの完了を伝える
func startTwitterStream(stopchan <-chan struct{}, votes chan <- string) <-chan struct{} {
	stoppedchan := make(chan struct{}, 1)
	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()
		for {
			select {
			case <-stopchan:
				log.Println("Twitterへの問い合わせを終了します...")
				return
			default:
				log.Println("Twitterに問い合わせします...")
				readFromTwitter(votes)
				log.Println(" (待機中)")
				time.Sleep(10 * time.Second) //待機してから再接続
			}
		}
	}()
	return stoppedchan
}

