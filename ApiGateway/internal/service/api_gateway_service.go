package service

import (
	"net/http"
	"net/url"
	"strings"
)

func StripSuffix(suffix string, h http.Handler) http.Handler {
	if suffix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimSuffix(r.URL.Path, suffix)
		rp := strings.TrimSuffix(r.URL.RawPath, suffix)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}
