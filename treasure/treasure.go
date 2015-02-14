package treasure

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var proxies []string
var lock *sync.Mutex

func refreshProxy() {
	defer lock.Unlock()
	lock.Lock()
	proxies = []string{}
	c := HttpClient(5 * time.Second)
	out, status, err := HttpGet(c, "http://www.kuaidaili.com/api/getproxy/?orderid=952370142723238&num=100&area=%E4%B8%AD%E5%9B%BD&browser=1&protocol=1&method=1&an_ha=1&sp2=1&sort=2&sep=2")
	out = "120.27.42.68:7182"
	if err == nil && status == http.StatusOK {
		tks := strings.Split(out, "\n")
		for _, tk := range tks {
			addr := "http://" + tk
			if !checkProxy(addr) {
				log.Println("failed", addr)
				continue
			}
			log.Println("success", addr)
			proxies = append(proxies, addr)
		}
	}

}

func init() {
	ticker := time.NewTicker(time.Minute * 10)
	lock = &sync.Mutex{}
	refreshProxy()
	go func() {
		for _ = range ticker.C {
			refreshProxy()
		}
	}()
}

func getProxy() *url.URL {
	defer lock.Unlock()
	lock.Lock()
	if len(proxies) == 0 {
		return nil
	}
	k := rand.Intn(len(proxies))
	ret, err := url.Parse(proxies[k])
	if err != nil {
		return nil
	}
	return ret
}

func HttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(network, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: timeout,
			DisableCompression:    false,
		},
	}
}

func HttpProxyClient(timeout time.Duration, proxy *url.URL) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(network, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: timeout,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(proxy),
		},
	}
}

func HttpGet(c *http.Client, link string) (string, int, error) {
	reqest, _ := http.NewRequest("GET", link, nil)
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94")
	response, err := c.Do(reqest)
	if err != nil {
		return "", http.StatusRequestTimeout, err
	}
	if response.StatusCode != http.StatusOK {
		return "", response.StatusCode, errors.New("status not ok")
	}
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	if response.Body == nil {
		return "", http.StatusInternalServerError, errors.New("response body is nil")
	}
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	html, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return "", response.StatusCode, err
	}
	return strings.Trim(string(html), " \n\t"), http.StatusOK, nil
}

func HttpPost(c *http.Client, link string, params url.Values) (string, int, error) {
	postDataStr := params.Encode()
	postDataBytes := []byte(postDataStr)
	reqest, err := http.NewRequest("POST", link, bytes.NewReader(postDataBytes))
	if err != nil {
		return "", http.StatusOK, err
	}
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94")
	reqest.Header.Set("Accept-Encoding", "gzip")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	response, err := c.Do(reqest)

	if err != nil {
		return "", http.StatusRequestTimeout, err
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	if response.Body == nil {
		return "", http.StatusInternalServerError, errors.New("response body is nil")
	}
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	html, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", response.StatusCode, err
	}
	return strings.Trim(string(html), " \n\t"), response.StatusCode, nil
}

func checkProxy(proxyAddr string) bool {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return false
	}
	c := HttpProxyClient(time.Second, proxy)
	_, status, err := HttpGet(c, "http://www.baidu.com")
	if err != nil {
		log.Println(err)
		return false
	}
	return status == http.StatusOK
}
