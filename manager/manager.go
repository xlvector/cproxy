//manager a cluster of proxy
package manager

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	Link          string
	LastHeartBeat time.Time
}

var proxies map[string]*Proxy
var proxyList []*Proxy

func init() {
	proxies = make(map[string]*Proxy)
	proxyList = make([]*Proxy, 0, 10)
}

func checkProxy(link string) bool {
	log.Println("begin check:", link)
	proxy, err := url.Parse(link)
	if err != nil {
		return false
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(5000)
				c, err := net.DialTimeout(network, addr, 5000)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: 5000,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(proxy),
		},
	}
	resp, err := client.Get("http://54.223.171.0:7183/check")
	if err != nil {
		return false
	}
	if resp == nil {
		return false
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	t := strings.Trim(string(b), " \n\t\r")
	log.Println(t)
	return strings.Contains(link, t)
}

func Register(link string) {
	if _, ok := proxies[link]; ok {
		return
	}
	if !checkProxy(link) {
		return
	}
	p := &Proxy{
		Link:          link,
		LastHeartBeat: time.Now(),
	}
	proxies[link] = p
	proxyList = append(proxyList, p)
}

func HeartBeat(link string) {
	if p, ok := proxies[link]; ok {
		p.LastHeartBeat = time.Now()
	}
}

func Select() *Proxy {
	for i := 0; i < len(proxyList) && i < 3; i++ {
		k := rand.Intn(len(proxyList))
		if time.Now().Sub(proxyList[k].LastHeartBeat).Minutes() < 5 {
			return proxyList[k]
		}
	}
	return nil
}

func HandleRegister(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	proxy := params.Get("proxy")
	Register(proxy)
	fmt.Fprint(rw, "ok")
}

func HandleHeartBeat(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	proxy := params.Get("proxy")
	HeartBeat(proxy)
	fmt.Fprint(rw, "ok")
}

func HandleSelect(rw http.ResponseWriter, req *http.Request) {
	p := Select()
	if p != nil {
		fmt.Fprint(rw, p.Link)
	}
}

func HandleCheck(rw http.ResponseWriter, req *http.Request) {
	tks := strings.Split(req.RemoteAddr, ":")
	fmt.Fprint(rw, tks[0])
}
