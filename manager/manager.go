//manager a cluster of proxy
package manager

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	Link          string
	LastCheckOKTm time.Time
	LastCheckOK   bool
	TotalSecs     float64
	TotalWeight   float64
	Ticker        *time.Ticker
}

func (p *Proxy) AveSecs() float64 {
	return p.TotalSecs / p.TotalWeight
}

var proxies map[string]*Proxy
var proxyList []*Proxy

func init() {
	proxies = make(map[string]*Proxy)
	proxyList = make([]*Proxy, 0, 10)
}

func checkProxy(link string) (bool, float64) {
	log.Println("begin check:", link)
	start := time.Now()
	proxy, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(5 * time.Second)
				c, err := net.DialTimeout(network, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: 5 * time.Second,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(proxy),
		},
	}
	resp, err := client.Get("http://54.223.171.0:7183/check")
	if err != nil {
		log.Println(err)
		return false, 0
	}
	if resp == nil {
		log.Println("resp is nil")
		return false, 0
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	t := strings.Trim(string(b), " \n\t\r")
	log.Println(t)
	return strings.Contains(link, t), time.Now().Sub(start).Seconds()
}

func Register(link string) {
	if _, ok := proxies[link]; ok {
		return
	}
	alive, secs := checkProxy(link)
	if !alive {
		continue
	}
	p := &Proxy{
		Link:          link,
		LastCheckOKTm: time.Now(),
		LastCheckOK:   true,
		TotalSecs:     secs,
		TotalWeight:   1.0,
		Ticker:        time.NewTicker(time.Minute),
	}
	proxies[link] = p
	proxyList = append(proxyList, p)
	go func() {
		for _ = range p.Ticker.C {
			alive, secs := checkProxy(link)
			if alive {
				p.LastCheckOK = true
				p.LastCheckOKTm = time.Now()
				p.TotalSecs = p.TotalSecs*0.7 + secs
				p.TotalWeight = p.TotalWeight*0.7 + 1.0
			} else {
				p.LastCheckOK = false
				p.TotalSecs = p.TotalSecs*0.7 + 5.0
				p.TotalWeight = p.TotalWeight*0.7 + 1.0
			}
		}
	}()
}

func Select() *Proxy {
	for i := 0; i < len(proxyList) && i < 3; i++ {
		k := rand.Intn(len(proxyList))
		if time.Now().Sub(proxyList[k].LastCheckOKTm).Minutes() < 5 && proxyList[k].LastCheckOK && proxyList[k].AveSecs() < 2.0 {
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
