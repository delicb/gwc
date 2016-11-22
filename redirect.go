package gwc

import (
	"errors"
	"net/http"
)

// CopyHeadersRedirect is CheckRedirect function that fixes issue
// github.com/golang/go/issues/4800. Namely, when redirect response is
// received, header from original request are not copied to new request.
// This function fixes that by copying all headers to new request except
// Authorization header (to avoid credentials leak). Authorization header
// is copied only if Host header from previous request matches with Host
// header of new request.
// This function implements same behaviour as standard lib does by default (error
// when more then 10 redirects happen), but it adds header copy feature.
func CopyHeadersRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}
	lastRequest := via[len(via)-1]

	for attr, val := range lastRequest.Header {
		// if hosts do not match do not copy Authorization header
		if attr == "Authorization" && req.Host != lastRequest.Host {
			continue
		}
		if _, ok := req.Header[attr]; !ok {
			req.Header[attr] = val
		}
	}
	return nil
}
