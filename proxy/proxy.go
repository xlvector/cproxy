package main

import (
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/xlvector/cproxy"
	"github.com/xlvector/cproxy/manager"
	"log"
	"net/http"
	"time"
)

func main() {
	proxyPort := flag.String("proxyPort", "7182", "port of cproxy")
	managePort := flag.String("managePort", "7183", "manage port")
	managerHost := flag.String("managerHost", "", "manager host")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().Do(&cproxy.RequestHandler{})
	proxy.OnResponse().Do(&cproxy.RespHandler{})

	if len(*managerHost) == 0 {
		http.HandleFunc("/register", manager.HandleRegister)
		http.HandleFunc("/heartbeat", manager.HandleHeartBeat)
		http.HandleFunc("/select", manager.HandleSelect)
		http.HandleFunc("/check", manager.HandleCheck)
		s := &http.Server{
			Addr:           ":" + *managePort,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   40 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go s.ListenAndServe()
	} else {
		ticker := time.NewTicker(time.Minute)
		c := http.Client{}
		go func() {
			for _ = range ticker.C {
				c.Get(*managerHost)
			}
		}()
	}
	log.Fatal(http.ListenAndServe(":"+*proxyPort, proxy))
}
