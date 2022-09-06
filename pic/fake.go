package pic

import (
	"net/http"
	"net/url"
)

func RandomUA() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"
}

func DefaultHeader() *http.Header {
	header := &http.Header{}
	header.Add("user-agent", RandomUA())
	return header
}

func RandomProxy() *http.Transport {
	proxyUrl, _ := url.Parse("http://127.0.0.1:18081")
	return &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func RandomClient() *http.Client {
	return &http.Client{Transport: RandomProxy()}
}
