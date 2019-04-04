package bot

import "strconv"

type SOCKS5ProxyConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (socks5 SOCKS5ProxyConfig) Address() string {
	return socks5.Host + ":" + strconv.Itoa(socks5.Port)
}

func (proxy *SOCKS5ProxyConfig) NoSecurity() bool {
	return proxy.Password == "" && proxy.Username == ""
}
