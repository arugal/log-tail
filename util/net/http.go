package net

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type HttpAuthWarper struct {
	h    http.Handler
	user string
	pwd  string
}

type HttpAuthMiddleware struct {
	user string
	pwd  string
}

func NewHttpAuthMiddleware(user, pwd string) *HttpAuthMiddleware {
	return &HttpAuthMiddleware{
		user: user,
		pwd:  pwd,
	}
}

func (auth *HttpAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUser, reqPwd, hasAuth := r.BasicAuth()
		if (auth.user == "" && auth.pwd == "") ||
			(hasAuth && reqUser == auth.user && reqPwd == auth.pwd) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}

type CrossDomainMiddleware struct {
}

func (cd *CrossDomainMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Request-Method", "*")
		next.ServeHTTP(w, r)
	})
}

func NewCrossDomainMiddleware() *CrossDomainMiddleware {
	return &CrossDomainMiddleware{}
}

type HttpGzipWraper struct {
	h http.Handler
}

func (gw *HttpGzipWraper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gizp") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		gw.h.ServeHTTP(gzr, r)
	}
}

func MakeHttpGzipHandler(h http.Handler) http.Handler {
	return &HttpGzipWraper{
		h: h,
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
