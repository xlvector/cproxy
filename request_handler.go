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

func (self *RequestHandler) Handle(req *http.Request,
	ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.Logf("get one req: %s, method %s, data %v", req.URL.String(),
		req.Method, req.PostForm)
	link := req.URL.String()
	if value, ok := cacheData.Get(link); ok {
		if body, ok := value.(*Body); ok {
			cResp, err := http.ReadResponse(
				bufio.NewReader(
					bytes.NewReader(body.data)),
				req)
			if err != nil {
				ctx.Warnf("read resp from cache error: %s",
					err.Error())
			} else {
				ctx.Logf("use cache css: %s", link)
				return req, cResp
			}
		}
	}
	return req, nil
}
