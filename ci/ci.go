package ci

import (
	"github.com/elazarl/goproxy"
	"github.com/xlvector/cproxy"
	"github.com/xlvector/cproxy/manager"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/register", manager.HandleRegister)
	http.HandleFunc("/select", manager.HandleSelect)
	http.HandleFunc("/check", manager.HandleCheck)
	http.HandleFunc("/list", manager.HandleList)
	s := &http.Server{
		Addr:           ":10086",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   40 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().Do(&cproxy.RequestHandler{})
	proxy.OnResponse().Do(&cproxy.RespHandler{})
	go http.ListenAndServe(":10087", proxy)
}
