package redir

import (
	"errors"
	"net"

	"github.com/yuche/clashr/component/socks5"
)

func parserPacket(conn net.Conn) (socks5.Addr, error) {
	return nil, errors.New("Windows not support yet")
}
