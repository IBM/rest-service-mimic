package handlers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rest-service-mimic/routes"
	"time"
)

type ProxyHandler struct {
	Config routes.Route
}

func (handler ProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	proxyUrl, err := url.Parse(handler.Config.Response.Proxy.Host)
	if err != nil {
		log.Printf("Error reading proxy url: %v", err)
		http.Error(w, "can't read proxy url", http.StatusBadRequest)
		return
	}

	r.URL.Scheme = proxyUrl.Scheme
	r.URL.Host = proxyUrl.Host
	r.Host = proxyUrl.Host

	fmt.Printf("Received request for proxying [%s] %s -H %s \n\n", r.Method, r.URL, r.Header)

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

	if handler.Config.Response.Proxy.InsecureSkipVerify == true {
		handler.skipSslVerify(proxy)
	}

	proxy.ServeHTTP(w, r)
}

func (handler ProxyHandler) skipSslVerify(proxy *httputil.ReverseProxy) {
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
}
