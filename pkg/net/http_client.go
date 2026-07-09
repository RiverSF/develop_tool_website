package net

import (
	"net"
	"net/http"
	"time"
)

var (
	sharedTransport = &http.Transport{
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     5 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	HttpClient1000   = createHTTPClient(1000)
	HttpClient2000   = createHTTPClient(2000)
	HttpClient5000   = createHTTPClient(5000)
	HttpClient120000 = createHTTPClient(120000)
)

func createHTTPClient(requestTimeout int) *http.Client {
	return &http.Client{
		Transport: sharedTransport,
		Timeout:   time.Duration(requestTimeout) * time.Millisecond,
	}
}
