package gwc

import (
	"net/http"

	"github.com/delicb/cliware"
)

// Client is main point of contact with this this library. It is used to set
// up basic configuration that will be used for all requests (except if it is
// overridden on per-request basis). For better performance (reuse of connections)
// only one instance of client should be created.
type Client struct {
	Middleware *cliware.Chain
	client     *http.Client
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
		client:     client,
		Middleware: cliware.NewChain(middlewares...),
	}
}

// Use adds provided middleware to this clients middleware chain.
func (c *Client) Use(middleware cliware.Middleware) *Client {
	c.Middleware.Use(middleware)
	return c
}

// UseFunc adds provided function to this clients middleware chain.
func (c *Client) UseFunc(m func(cliware.Handler) cliware.Handler) *Client {
	c.Middleware.UseFunc(m)
	return c
}

// Request creates and returns new request that uses this client to perform
// HTTP request and uses its defined middlewares.
func (c *Client) Request() *Request {
	return NewRequest(c, c.Middleware.ChildChain())
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
