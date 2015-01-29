package main

import (
	"flag"
	"github.com/BigTong/cproxy"
	"github.com/elazarl/goproxy"
	"github.com/xlvector/cproxy"
	"log"
	"net/http"
	"time"
)

func main() {
	proxyPort := flag.String("proxyPort", "7182", "port of cproxy")
	randcodePort := flag.String("randcodePort", "9100", "port of randcoe server")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().Do(&cproxy.RequestHandler{})
	proxy.OnResponse().Do(&cproxy.RespHandler{})

	randCodeServer := &cproxy.RandCodeServer{}
	http.Handle("/randcode", randCodeServer)
	s := &http.Server{
		Addr:           ":" + *randcodePort,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   40 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	log.Fatal(http.ListenAndServe(":"+*proxyPort, proxy))
}
