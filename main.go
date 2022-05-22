package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newReverseProxy(target string) *httputil.ReverseProxy {

	u, _ := url.Parse(target)
	rp := httputil.NewSingleHostReverseProxy(u)
	old := rp.Director
	rp.Director = func(r *http.Request) {
		old(r)
		modifyRequest(r)
	}

	rp.ModifyResponse = modifyResponse()

	return rp
}

func modifyResponse() func(*http.Response) error {
	return func(r *http.Response) error {
		r.Header.Add("resp", "成功响应")
		return nil
	}
}

func modifyRequest(r *http.Request) {
	r.Header.Set("name", "linuxea")
}

func reverProxyServer(reverProxyServer *httputil.ReverseProxy) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reverProxyServer.ServeHTTP(w, r)
	})
}

func main() {

	go func() {
		rp := newReverseProxy("http://127.0.0.1:8080")
		http.HandleFunc("/", reverProxyServer(rp))
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			panic(err.Error())
		}
	}()

	go func() {
		sm := &http.ServeMux{}
		sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("8080收到请求")
			resp := fmt.Sprintf("你好，已经收到 %s 的信息: %s", r.Header.Get("name"), r.Header.Get("msg"))
			w.Write([]byte(resp))
		})
		err2 := http.ListenAndServe(":8080", sm)
		if err2 != nil {
			panic(err2.Error())
		}
	}()

	<-make(chan struct{})

}
