package ci

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestProxyCheck(t *testing.T) {
	proxy, _ := url.Parse("http://127.0.0.1:10087")
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
	{
		resp, err := client.Get("http://127.0.0.1:10086/check?aaa")
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if string(b) != "GET:127.0.0.1" {
			t.Error(string(b))
		}
	}
	{
		req, err := http.NewRequest("POST", "http://127.0.0.1:10086/check?bbb", nil)
		if err != nil {
			t.Error(err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if string(b) != "POST:127.0.0.1" {
			t.Error(string(b))
		}
	}
}
