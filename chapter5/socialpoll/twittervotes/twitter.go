package main
import (
	"net"
	"time"
	"io"
	"github.com/garyburd/go-oauth/oauth"
	"go-web-oreilly/chapter5/socialpoll/vendor/github.com/joeshaw/envdecode"
	"log"
	"go-web-oreilly/chapter5/socialpoll/vendor/github.com/garyburd/go-oauth/oauth"
	"sync"
	"net/http"
	"net/url"
	"strconv"
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
