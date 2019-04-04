package proxy

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type Auth = proxy.Auth

var NoAuth = Auth{}

func SOCKS5(proxyAddr string, authData Auth) (*http.Client, error) {
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