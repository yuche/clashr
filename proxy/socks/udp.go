package socks

import (
	"bytes"
	"net"

	adapters "github.com/yuche/clashr/adapters/inbound"
	"github.com/yuche/clashr/common/pool"
	"github.com/yuche/clashr/component/socks5"
	C "github.com/yuche/clashr/constant"
)

type SockUDPListener struct {
	net.PacketConn
	address string
	closed  bool
}

func NewSocksUDPProxy(addr string) (*SockUDPListener, error) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}

	sl := &SockUDPListener{l, addr, false}
	go func() {
		for {
			buf := pool.BufPool.Get().([]byte)
			n, remoteAddr, err := l.ReadFrom(buf)
			if err != nil {
				pool.BufPool.Put(buf[:cap(buf)])
				if sl.closed {
					break
				}
				continue
			}
			handleSocksUDP(l, buf[:n], remoteAddr)
		}
	}()

	return sl, nil
}

func (l *SockUDPListener) Close() error {
	l.closed = true
	return l.PacketConn.Close()
}

func (l *SockUDPListener) Address() string {
	return l.address
}

func handleSocksUDP(pc net.PacketConn, buf []byte, addr net.Addr) {
	target, payload, err := socks5.DecodeUDPPacket(buf)
	if err != nil {
		// Unresolved UDP packet, return buffer to the pool
		pool.BufPool.Put(buf[:cap(buf)])
		return
	}
	conn := &fakeConn{
		PacketConn: pc,
		remoteAddr: addr,
		targetAddr: target,
		buffer:     bytes.NewBuffer(payload),
		bufRef:     buf,
	}
	tun.Add(adapters.NewSocket(target, conn, C.SOCKS, C.UDP))
}
