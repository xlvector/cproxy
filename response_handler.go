package cproxy

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type RespHandler struct {
}

func (r *RespHandler) Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	ctx.Logf("finish one req:", resp.Request.URL.String())
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return resp
	}
	link := resp.Request.URL.String()
	ctx.Logf("resp req url: %s", link)
	b, _ := httputil.DumpResponse(resp, true)

	if _, ok := imageCache.Get(link); ok {
		//id := md5.Sum([]byte(link + strconv.FormatInt(time.Now().UnixNano(), 10)))
		id := strconv.FormatInt(time.Now().UnixNano(), 10)
		imageCache.Set(id, b, 0)
		resp.Header.Set("randcode_url:", "http://127.0.0.1:9100/randcode?id="+id)
		return resp
	}

	if data, ok := cacheData.Get(link); ok {
		if body, ok := data.(*Body); ok {
			body.hit += 1
			if body.immutable == IMMUTABLE_NO {
				return resp
			}

			if body.sn != time33(b) {
				body.immutable = IMMUTABLE_NO
			}

			if body.immutable == IMMUTABLE_UNKNOWN && body.hit > 2 && body.sn == time33(b) {
				body.immutable = IMMUTABLE_YES
			}

			return resp
		}
	}

	body := &Body{
		data:      b,
		sn:        time33(b),
		immutable: IMMUTABLE_UNKNOWN,
		hit:       1,
	}
	cacheData.Set(link, body, 0)
	return resp
}
