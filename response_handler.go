package cproxy

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"net/http/httputil"
)

type RespHandler struct {
}

func (r *RespHandler) Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return resp
	}
	link := resp.Request.URL.String()
	ctx.Logf("resp req url: %s", link)
	b, _ := httputil.DumpResponse(resp, true)

	if len(resp.Request.URL.Query()) == 0 {
		body := &Body{
			data: b,
		}
		cacheData.Set(link, body, 0)
	}
	return resp
}
