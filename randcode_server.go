package cproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pmylund/go-cache"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var imageCache = cache.New(10*time.Minute, 5*time.Minute)

type RandCodeServer struct {
}

func (server *RandCodeServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("randcode server crashed")
		}
	}()

	if id := req.FormValue("id"); len(id) != 0 {
		if value, ok := imageCache.Get(id); ok {
			if data, ok := value.([]byte); ok {
				resp, _ := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), req)
				if resp.Body != nil {
					defer resp.Body.Close()
				}
				htmlByte, _ := ioutil.ReadAll(resp.Body)
				fmt.Fprint(w, string(htmlByte))
				return
			}
		}
		fmt.Fprint(w, "time out")
		return
	}
	fmt.Fprint(w, "args wrong")
	return
}
