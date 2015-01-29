package cproxy

import (
	"bufio"
	"bytes"
	"github.com/elazarl/goproxy"
	"github.com/pmylund/go-cache"
	"net/http"
	"strings"
	"time"
)

var cacheData = cache.New(12*time.Hour, 1*time.Hour)

type RequestHandler struct {
}

func (self *RequestHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.Logf("get one req:%s", req.URL.String())
	if randcode := req.FormValue("randcode"); randcode == "true" {
		ctx.Logf("get request uri:%s", req.URL.String())
		req.RequestURI = strings.TrimRight(req.RequestURI, "&randcode=true")
		id := req.URL.String()
		imageCache.Set(id, "", 0)
		return req, nil
	}

	link := req.URL.String()
	if value, ok := cacheData.Get(link); ok {
		if body, ok := value.(*Body); ok && body.immutable == IMMUTABLE_YES {
			cResp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(body.data)), req)
			if err == nil {
				ctx.Logf("use cache css: %s", link)
				return req, cResp
			}
			ctx.Warnf("read resp from cache error: %s", err.Error())
		}
	}
	return req, nil
}
