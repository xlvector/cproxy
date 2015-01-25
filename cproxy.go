package main

import (
	"bufio"
	"bytes"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/http/httputil"
)

func time33(b []byte) int64 {
	ret := int64(0)
	for _, e := range b {
		ret *= 33
		ret += int64(e)
	}
	return ret
}

const (
	IMMUTABLE_UNKNOWN = 0
	IMMUTABLE_YES     = 1
	IMMUTABLE_NO      = 2
)

type Body struct {
	data      []byte
	sn        int64
	immutable int
	hit       int
}

var cache map[string]*Body

type RespHandler struct {
}

func (r *RespHandler) Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return resp
	}
	link := resp.Request.URL.String()
	ctx.Logf("resp req url: %s", link)
	b, _ := httputil.DumpResponse(resp, true)
	body, ok := cache[link]
	if !ok {
		body = &Body{
			data:      b,
			sn:        time33(b),
			immutable: IMMUTABLE_UNKNOWN,
			hit:       1,
		}
		cache[link] = body
	} else {
		if body.immutable != IMMUTABLE_NO {
			if body.sn != time33(b) {
				body.immutable = IMMUTABLE_NO
			}
			if body.immutable == IMMUTABLE_UNKNOWN && body.hit > 2 && body.sn == time33(b) {
				body.immutable = IMMUTABLE_YES
			}
		}
		body.hit += 1
	}
	return resp
}

func main() {
	cache = make(map[string]*Body)
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			link := r.URL.String()
			body, ok := cache[link]
			if ok && body.immutable == IMMUTABLE_YES {
				cResp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(body.data)), r)
				if err != nil {
					ctx.Warnf("read resp from cache error: %s", err.Error())
				} else {
					ctx.Logf("use cache css: %s", link)
					return r, cResp
				}
			}
			return r, nil
		})

	proxy.OnResponse().Do(&RespHandler{})
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
