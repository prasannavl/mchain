package builder

import (
	"net/http"

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
		h = hconv.FromHttp(http.DefaultServeMux)
	}
	c := b.chain
	mx := c.Middlewares
	mLen := len(mx)
	for i := range mx {
		h = c.Middlewares[mLen-1-i](h)
	}
	return h
}

func (b *ChainBuilder) BuildHttp(errorHandler func(error)) http.Handler {
	return hconv.ToHttp(b.Build(), errorHandler)
}
