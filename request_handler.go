package cproxy

import (
	"bufio"
	"bytes"
	"github.com/elazarl/goproxy"
	"github.com/pmylund/go-cache"
	"net/http"
	"time"
)

var cacheData = cache.New(12*time.Hour, 1*time.Hour)

type RequestHandler struct {
}

func (self *RequestHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.Logf("get one req:%s", req.URL.String())
	params := req.URL.Query()
	if randcode := params.Get("randcode"); randcode == "true" {
		params.Del("randcode")
		req.URL.RawQuery = params.Encode()
		id := req.URL.String()
		ctx.Logf("get randoce id:%s", id)
		imageCache.Set(id, "", 0)
		return req, nil
	}

	link := req.URL.String()
	if value, ok := cacheData.Get(link); ok {
		if body, ok := value.(*Body); ok && body.immutable == IMMUTABLE_YES {
			cResp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(body.data)), req)
			if err != nil {
				ctx.Warnf("read resp from cache error: %s", err.Error())
			} else {
				ctx.Logf("use cache css: %s", link)
				return req, cResp
			}
		}
	}
	return req, nil
}
