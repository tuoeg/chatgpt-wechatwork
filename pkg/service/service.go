package service

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Service interface {
	Auth(msg_signature, timestamp, nonce, echostr string) (string, error)
	SendMsg(string, string, string, []byte) (string, error)
	sendMsgToWx(user, msg, token string)
	sendMsgToOpenAI(msg, user string) (string, error)
}

type service struct {
	httpClient *http.Client
}

func NewService() Service {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	return &service{httpClient: client}
}
