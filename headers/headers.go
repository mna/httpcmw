// Package headers defines a middleware that adds static headers to
// the requests.
package headers

import (
	"net/http"

	"github.com/PuerkitoBio/httpcmw"
)

// Headers is an http.Header map that implements the httpcmw.Wrapper
// interface so that the headers are added to each request using the
// middleware. By default the header values are set (replace any existing
// value), but the behaviour can be controlled by prepending a "+" to
// the header name (add the value) or a "-" (remove this header).
type Headers http.Header

// Add adds the value v to the header k.
func (hd Headers) Add(k, v string) {
	http.Header(hd).Add(k, v)
}

// Set sets the value v to the header k, replacing any existing value.
func (hd Headers) Set(k, v string) {
	http.Header(hd).Set(k, v)
}

// Get returns the first value of the header k.
func (hd Headers) Get(k string) string {
	return http.Header(hd).Get(k)
}

// Del removes the header k.
func (hd Headers) Del(k string) {
	http.Header(hd).Del(k)
}

// Wrap returns a handler that adds the headers to the request's Header.
func (hd Headers) Wrap(d httpcmw.Doer) httpcmw.Doer {
	return httpcmw.DoerFunc(func(r *http.Request) (*http.Response, error) {
		reqHd := r.Header
		for k, v := range hd {
			start := byte(' ')
			if len(k) > 0 {
				start = k[0]
			}
			switch start {
			case '+':
				k := k[1:]
				for _, vv := range v {
					reqHd.Add(k, vv)
				}
			case '-':
				reqHd.Del(k[1:])
			default:
				reqHd[k] = v
			}
		}
		return d.Do(r)
	})
}
