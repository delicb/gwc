package gwc

import (
	"context"

	"go.delic.rs/cliware"
)

// Layer is container of any additional middlewares that group of endpoints use.
type Layer struct {
	Next  Doer
	Chain *cliware.Chain
}

// NewLayer creates and returns new instance of Layer that will apply all provided
// middlewares and use next Doer to do actual work.
func NewLayer(next Doer, middlewares ...cliware.Middleware) *Layer {
	return &Layer{
		Next:  next,
		Chain: cliware.NewChain(middlewares...),
	}
}

// Exec is implementation of cliware.Middleware interface.
func (s *Layer) Exec(handler cliware.Handler) cliware.Handler {
	return s.Chain.Exec(handler)
}

// Use adds provided middlewares to this layers chain.
func (s *Layer) Use(middleware ...cliware.Middleware) *Layer {
	s.Chain.Use(middleware...)
	return s
}

// Do applies all middlewares from this layer and provided middlewares and
// calls next Doer to do actual work.
func (s *Layer) Do(middlewares ...cliware.Middleware) (*Response, error) {
	return s.DoCtx(context.Background(), middlewares...)
}

// DoCtx applies all middlewares from this layer and provided middlewares and
// calls next Doer to do actual work.
func (s *Layer) DoCtx(ctx context.Context, middlewares ...cliware.Middleware) (*Response, error) {
	// insert service itself to first place in middleware chain
	middlewares = append(middlewares, nil)
	copy(middlewares[1:], middlewares[0:])
	middlewares[0] = s
	return s.Next.DoCtx(ctx, middlewares...)
}
