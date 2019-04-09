package protectron

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type ProxyAuth = proxy.Auth

var NoAuth = ProxyAuth{}

// use NoAuth as authData in case if no authentication required
// In that case proxy.SOCKS5 expects nil, not &proxy.Auth{},
// or else emits error "invalid username/password"
func SOCKS5(proxyAddr string, authData ProxyAuth) (*http.Client, error) {
	const timeout = 5 * time.Second
	var auth *proxy.Auth
	if authData != NoAuth {
		auth = &authData
	}
	var tcpDialer = &net.Dialer{
		Timeout: timeout,
	}
	var proxyDialer, errSOCKS5 = proxy.SOCKS5("tcp", proxyAddr, auth, tcpDialer)
	if errSOCKS5 != nil {
		return nil, errSOCKS5
	}
	return &http.Client{
		Timeout: 2 * timeout,
		Transport: &http.Transport{
			Dial: proxyDialer.Dial,
		},
	}, nil
}
