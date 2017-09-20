package builder

import (
	"net/http"

	"github.com/prasannavl/goerror/httperror"

	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/hconv"
)

type ChainBuilder struct {
	chain   mchain.Chain
	handler mchain.Handler
}

func Create(middlewares ...mchain.Middleware) ChainBuilder {
	return ChainBuilder{mchain.Chain{middlewares}, nil}
}

func (b *ChainBuilder) Add(m ...mchain.Middleware) *ChainBuilder {
	b.chain.Middlewares = append(b.chain.Middlewares, m...)
	return b
}

func (b *ChainBuilder) Handler(finalHandler mchain.Handler) *ChainBuilder {
	b.handler = finalHandler
	return b
}

func (b *ChainBuilder) Build() mchain.Handler {
	h := b.handler
	if h == nil {
		h = mchain.HandlerFunc(defaultHandler)
	}
	c := b.chain
	mx := c.Middlewares
	mLen := len(mx)
	for i := range mx {
		h = c.Middlewares[mLen-1-i](h)
	}
	return h
}

func (b *ChainBuilder) BuildHttp(errorHandler mchain.ErrorHandler) http.Handler {
	return hconv.ToHttp(b.Build(), errorHandler)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) error {
	// Send with an empty message, so the default pipeline doesn't
	// even touch the content. Just the status code.
	// And set end=false, so that other middlewares that can take the hint
	// optionally can and continue instead of halting entirely.
	return httperror.New(http.StatusNotFound, "handler not found", false)
}
