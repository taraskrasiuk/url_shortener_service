package httpdialer

import (
	"crypto/tls"
	"net"
	u "net/url"
	"strings"
)

type httpAddr struct {
	dialer *net.Dialer
	isTLS  bool
	addr   string
}

func proxyAddr(url *u.URL) (*httpAddr, error) {
	// by default set values default values as for http
	h := &httpAddr{
		dialer: &net.Dialer{},
		isTLS:  false,
		addr:   url.Host + ":80",
	}

	if strings.ToLower(url.Scheme) == "https" {
		h.isTLS = true
		h.addr = url.Host + ":443"
	}
	return h, nil
}

func HttpDialProxy(url *u.URL) (net.Conn, error) {
	h, err := proxyAddr(url)
	if err != nil {
		return nil, err
	}
	if h.isTLS {
		return tls.DialWithDialer(h.dialer, "tcp", h.addr, nil)
	}
	return net.Dial("tcp", h.addr)
}
