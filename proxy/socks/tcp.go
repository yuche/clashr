package socks

import (
	"io"
	"io/ioutil"
	"net"

	adapters "github.com/yuche/clashr/adapters/inbound"
	"github.com/yuche/clashr/component/socks5"
	C "github.com/yuche/clashr/constant"
	"github.com/yuche/clashr/log"
	authStore "github.com/yuche/clashr/proxy/auth"
	"github.com/yuche/clashr/tunnel"
)

var (
	tun = tunnel.Instance()
)

type SockListener struct {
	net.Listener
	address string
	closed  bool
}

func NewSocksProxy(addr string) (*SockListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	sl := &SockListener{l, addr, false}
	go func() {
		log.Infoln("SOCKS proxy listening at: %s", addr)
		for {
			c, err := l.Accept()
			if err != nil {
				if sl.closed {
					break
				}
				continue
			}
			go handleSocks(c)
		}
	}()

	return sl, nil
}

func (l *SockListener) Close() {
	l.closed = true
	_ = l.Listener.Close()
}

func (l *SockListener) Address() string {
	return l.address
}

func handleSocks(conn net.Conn) {
	target, command, err := socks5.ServerHandshake(conn, authStore.Authenticator())
	if err != nil {
		_ = conn.Close()
		return
	}
	_ = conn.(*net.TCPConn).SetKeepAlive(true)
	if command == socks5.CmdUDPAssociate {
		defer conn.Close()
		_, _ = io.Copy(ioutil.Discard, conn)
		return
	}
	tun.Add(adapters.NewSocket(target, conn, C.SOCKS, C.TCP))
}
