package http

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/PereRohit/util/crypt"
)

var (
	client    *http.Client
	singleton sync.Once
)

func GetClient() (*http.Client, error) {
	var singletonErr error

	if client == nil {
		singleton.Do(func() {
			localAddr, err := net.ResolveIPAddr("ip", "localhost")
			if err != nil {
				singletonErr = err
				return
			}
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						LocalAddr: &net.TCPAddr{
							IP: localAddr.IP,
						},
						Timeout: 5 * time.Second,
					}).DialContext,
					TLSHandshakeTimeout: 5 * time.Second,
					IdleConnTimeout:     10 * time.Second,
				},
				Timeout: 10 * time.Second,
			}
		})
	}
	return client, singletonErr
}

func EncryptAndCall(httpMethod, endpointUrl, message, encryptionKey string, shouldEncrypt bool, header map[string]string) (*http.Response, error) {
	var err error
	if shouldEncrypt && encryptionKey == "" {
		return nil, fmt.Errorf("encryption key required")
	} else if shouldEncrypt {
		message, err = crypt.RsaOaepEncrypt(message, encryptionKey)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(httpMethod, endpointUrl, bytes.NewReader([]byte(message)))
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}
	cli, err := GetClient()
	if err != nil {
		return nil, err
	}
	return cli.Do(req)
}

func Call(httpMethod, endpointUrl, message string, header map[string]string) (*http.Response, error) {
	return EncryptAndCall(httpMethod, endpointUrl, message, "", false, header)
}
