// デフォルトのレスポンスヘッダ

package main

import "net/http"

var defaultHeaders = map[string]string{
	"Server":                 "GoCORS/" + VERSION,
	"X-Content-Type-Options": "nosniff",
	"X-Frame-Options":        "DENY",
	"X-XSS-Protection":       "1; mode=block",
}

// DefaultHeaders : デフォルトのレスポンスヘッダ
type DefaultHeaders struct {
	Handler http.Handler
}

func (dh DefaultHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	for k, v := range defaultHeaders {
		headers.Set(k, v)
	}
	dh.Handler.ServeHTTP(w, r)
}
