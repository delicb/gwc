package gwc

import (
	"context"

	"go.delic.rs/cliware"
)

// Group is container of any additional middlewares that group of endpoints use.
type Group struct {
	Next  Doer
	Chain *cliware.Chain
}

// NewGroup creates and returns new instance of Group that will apply all provided
// middlewares and use next Doer to do actual work.
func NewGroup(next Doer, middlewares ...cliware.Middleware) *Group {
	return &Group{
		Next:  next,
		Chain: cliware.NewChain(middlewares...),
	}
}

// Exec is implementation of cliware.Middleware interface.
func (s *Group) Exec(handler cliware.Handler) cliware.Handler {
	return s.Chain.Exec(handler)
}

// Use adds provided middlewares to this group's chain.
func (s *Group) Use(middleware ...cliware.Middleware) *Group {
	s.Chain.Use(middleware...)
	return s
}

// Do applies all middlewares from this layer and provided middlewares and
// calls next Doer to do actual work.
func (s *Group) Do(middlewares ...cliware.Middleware) (*Response, error) {
	return s.DoCtx(context.Background(), middlewares...)
}

// DoCtx applies all middlewares from this layer and provided middlewares and
// calls next Doer to do actual work.
func (s *Group) DoCtx(ctx context.Context, middlewares ...cliware.Middleware) (*Response, error) {
	// insert service itself to first place in middleware chain
	middlewares = append(middlewares, nil)
	copy(middlewares[1:], middlewares[0:])
	middlewares[0] = s
	return s.Next.DoCtx(ctx, middlewares...)
}
