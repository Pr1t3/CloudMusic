package service

import (
	"errors"
	"io"
	"net/http"
)

type ProxyRequestStruct struct {
	ForbiddenError error
}

func NewProxyRequestStruct() *ProxyRequestStruct {
	return &ProxyRequestStruct{ForbiddenError: errors.New("status forbidden")}
}

func (p *ProxyRequestStruct) ProxyRequest(r *http.Request, target string, reqBody io.Reader, method string) ([]byte, *http.Header, error) {
	proxyReq, err := http.NewRequest(method, target, reqBody)
	if err != nil {
		return nil, nil, err
	}
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	for _, cookie := range r.Cookies() {
		proxyReq.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, p.ForbiddenError
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, &resp.Header, err
}
