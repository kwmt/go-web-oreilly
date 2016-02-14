package main
import (
	"net"
	"time"
	"io"
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
