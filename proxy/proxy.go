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
	host := flag.String("host", "", "host of current machine")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().Do(&cproxy.RequestHandler{})
	proxy.OnResponse().Do(&cproxy.RespHandler{})

	if len(*managerHost) == 0 {
		http.HandleFunc("/register", manager.HandleRegister)
		http.HandleFunc("/select", manager.HandleSelect)
		http.HandleFunc("/check", manager.HandleCheck)
		http.HandleFunc("/list", manager.HandleList)
		s := &http.Server{
			Addr:           ":" + *managePort,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   40 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go s.ListenAndServe()
	} else {
		go func() {
			c := http.Client{}
			for i := 0; i < 8; i++ {
				time.Sleep(time.Second * 20)
				_, err := c.Get(*managerHost + "/register?proxy=" + *host)
				if err == nil {
					break
				}
			}
		}()
	}
	log.Fatal(http.ListenAndServe(":"+*proxyPort, proxy))
}
