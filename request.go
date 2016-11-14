package gwc

import (
	"context"
	"net/http"

	"go.delic.rs/cliware"

	"go.delic.rs/cliware-middlewares/cookies"
	"go.delic.rs/cliware-middlewares/headers"
	"go.delic.rs/cliware-middlewares/query"
	cwurl "go.delic.rs/cliware-middlewares/url"
)

// Request is struct used to hold information (mostly middlewares) used
// to construct HTTP request.
type Request struct {
	//Middleware *cliware.Chain
	before  *cliware.Chain
	after   *cliware.Chain
	Client  *Client
	context context.Context
}

// NewRequest creates new instance of request for provided client and with
// initial middleware chain.
func NewRequest(client *Client, before *cliware.Chain, after *cliware.Chain) *Request {
	return &Request{
		Client:  client,
		before:  before,
		after:   after,
		context: nil,
	}
}

// Context returns instance of context.Context used for this request.
func (r *Request) Context() context.Context {
	return r.context
}

// SetContext sets instance of context.Context used for this request.
func (r *Request) SetContext(ctx context.Context) *Request {
	r.context = ctx
	return r
}

// Use adds provided middleware to this request middleware chain.
func (r *Request) Use(m cliware.Middleware) *Request {
	r.before.Use(m)
	return r
}

// UseFunc adds provided function to this request middleware chain.
func (r *Request) UseFunc(m func(cliware.Handler) cliware.Handler) *Request {
	r.before.UseFunc(m)
	return r
}

// Utility methods - shortcuts to using middlewares.

// Method sets HTTP method (verb) for this request.
func (r *Request) Method(method string) *Request {
	r.Use(headers.Method(method))
	return r
}

// URL parses and sets URL for this request.
func (r *Request) URL(rawURL string) *Request {
	r.Use(cwurl.URL(rawURL))
	return r
}

// BaseURL sets schema and host from provided URL to this request.
func (r *Request) BaseURL(rawURL string) *Request {
	r.Use(cwurl.BaseURL(rawURL))
	return r
}

// Path sets path to URL for this request.
func (r *Request) Path(path string) *Request {
	r.Use(cwurl.Path(path))
	return r
}

// AddPath appends path segment to current path for this request.
func (r *Request) AddPath(path string) *Request {
	r.Use(cwurl.AddPath(path))
	return r
}

// Param replaces key in caramelized URL with given value for this request.
func (r *Request) Param(key, value string) *Request {
	r.Use(cwurl.Param(key, value))
	return r
}

// Params replaces all keys in URL with key-value pairs provided in map for this request.
func (r *Request) Params(values map[string]string) *Request {
	r.Use(cwurl.Params(values))
	return r
}

// AddQuery adds query parameter to URL for this request.
// If parameter with same name already exist, new one will be appended. To
// replace it, use SetQuery.
func (r *Request) AddQuery(key, value string) *Request {
	r.Use(query.Add(key, value))
	return r
}

// SetQuery sets query parameter to URL for this request.
// If parameter already exists it will be replaced.
func (r *Request) SetQuery(key, value string) *Request {
	r.Use(query.Set(key, value))
	return r
}

// SetQueryParams sets query parameters for this request provided in the map.
func (r *Request) SetQueryParams(values map[string]string) *Request {
	r.Use(query.SetMap(values))
	return r
}

// SetHeader sets provided header value to current request.
// If header with same name already exists, it will be overridden. To add same
// header again, use AddHeader.
func (r *Request) SetHeader(key, value string) *Request {
	r.Use(headers.Set(key, value))
	return r
}

// AddHeader adds header with provided name and value to current request.
// It does not override existing headers.
func (r *Request) AddHeader(key, value string) *Request {
	r.Use(headers.Add(key, value))
	return r
}

// SetHeaders sets headers for this request provided in map.
func (r *Request) SetHeaders(headerMap map[string]string) *Request {
	r.Use(headers.SetMap(headerMap))
	return r
}

// AddCookie adds provided cookie to current request.
func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	r.Use(cookies.Add(cookie))
	return r
}

// SetCookie set cookie with provided name and value to current request.
func (r *Request) SetCookie(key, value string) *Request {
	r.Use(cookies.Set(key, value))
	return r
}

// sendRequest is private method that does actual request dispatching.
func (r *Request) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	reqWithContext := req.WithContext(ctx)
	return r.Client.client.Do(reqWithContext)
}

// Send constructs and sends HTTP request.
// This method uses all defined middlewares and client defined in requests
// to construct HTTP request.
func (r *Request) Send() (*Response, error) {
	r.before.Use(r.after)
	sender := r.before.Exec(cliware.HandlerFunc(r.sendRequest))

	// sender := r.Middleware.Exec(cliware.HandlerFunc(r.sendRequest))
	if r.context == nil {
		r.context = context.Background()
	}
	r.context = SetClient(r.context, r.Client.client)
	resp, err := sender.Handle(r.context, cliware.EmptyRequest())
	return BuildResponse(resp, err), err
}
