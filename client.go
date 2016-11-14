package gwc

import (
	"context"
	"net/http"

	"go.delic.rs/cliware"
)

// Client is main point of contact with this this library. It is used to set
// up basic configuration that will be used for all requests (except if it is
// overridden on per-request basis). For better performance (reuse of connections)
// only one instance of client should be created.
type Client struct {
	// Middleware     *cliware.Chain
	Before *cliware.Chain
	After  *cliware.Chain
	client *http.Client
	//PostMiddleware *cliware.Chain
}

// New creates and returns instance of a client.
// First parameter is http.Client instance that will be used to send actual
// requests. This can be used to modify default behavior (like TLS or
// timeouts). If nil is provided, http.DefaultClient will be used.
//
// Also, any number of middlewares can be provided. These will be used for
// each request. Additional middlewares can be added by using "Use*" methods
// on a client.
func New(client *http.Client, middlewares ...cliware.Middleware) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		client: client,
		Before: cliware.NewChain(middlewares...),
		After:  cliware.NewChain(),
		// Middleware:     cliware.NewChain(middlewares...),
		//PostMiddleware: cliware.NewChain(),
	}
}

// Use adds provided middleware to this clients middleware chain.
func (c *Client) Use(m cliware.Middleware) *Client {
	c.Before.Use(m)
	return c
}

// UseFunc adds provided function to this clients middleware chain.
func (c *Client) UseFunc(m func(cliware.Handler) cliware.Handler) *Client {
	c.Before.UseFunc(m)
	return c
}

// UsePost adds middleware that will be added to all requests sent by
// this client AFTER middlewares from request itself are executed.
func (c *Client) UsePost(m cliware.Middleware) *Client {
	c.After.Use(m)
	return c
}

// UsePOstFunc adds middleware that will be added to all requests sent by this
// client AFTER middlewares from request itself are executed.
func (c *Client) UsePostFunc(m func(cliware.Handler) cliware.Handler) *Client {
	c.After.UseFunc(m)
	return c
}

// Request creates and returns new request that uses this client to perform
// HTTP request and uses its defined middlewares.
func (c *Client) Request() *Request {
	return NewRequest(c, c.Before.Copy(), c.After.Copy())
}

// Get creates and returns new GET request.
func (c *Client) Get() *Request {
	r := c.Request()
	r.Method("GET")
	return r
}

// Post creates and returns new POST request.
func (c *Client) Post() *Request {
	r := c.Request()
	r.Method("POST")
	return r
}

// Put creates and returns new PUT request.
func (c *Client) Put() *Request {
	r := c.Request()
	r.Method("PUT")
	return r
}

// Delete creates and returns new DELETE request.
func (c *Client) Delete() *Request {
	r := c.Request()
	r.Method("DELETE")
	return r
}

// Patch creates and returns new PATCH request.
func (c *Client) Patch() *Request {
	r := c.Request()
	r.Method("PATCH")
	return r
}

// Head creates and returns new HEAD request.
func (c *Client) Head() *Request {
	r := c.Request()
	r.Method("HEAD")
	return r
}

// Options creates and returns new OPTIONS request.
func (c *Client) Options() *Request {
	r := c.Request()
	r.Method("OPTIONS")
	return r
}

// Do creates new request, applies all provided middlewares to it and sends request.
// Context that is used is context.Background()
func (c *Client) Do(middlewares ...cliware.Middleware) (*Response, error) {
	return c.DoCtx(context.Background(), middlewares...)
}

// DoCtx creates new request, applies all provided middlewares to it and sends request
// with provided context.
func (c *Client) DoCtx(ctx context.Context, middlewares ...cliware.Middleware) (*Response, error) {
	req := c.Request()
	req.SetContext(ctx)
	for _, m := range middlewares {
		req.Use(m)
	}
	//req.Use(c.PostMiddleware)
	return req.Send()
}
