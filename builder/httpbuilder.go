package builder

import (
	"net/http"

	"github.com/prasannavl/mchain"
)

type HttpChainBuilder struct {
	chain   mchain.HttpChain
	handler http.Handler
}

func CreateHttp(middlewares ...mchain.HttpMiddleware) HttpChainBuilder {
	return HttpChainBuilder{mchain.HttpChain{middlewares}, nil}
}

func (b *HttpChainBuilder) Add(m ...mchain.HttpMiddleware) *HttpChainBuilder {
	b.chain.Middlewares = append(b.chain.Middlewares, m...)
	return b
}

func (b *HttpChainBuilder) Handler(finalHandler http.Handler) *HttpChainBuilder {
	b.handler = finalHandler
	return b
}

func (b *HttpChainBuilder) Build() http.Handler {
	h := b.handler
	if h == nil {
		h = http.DefaultServeMux
	}
	c := b.chain
	mx := c.Middlewares
	for i := len(mx) - 1; i >= 0; i-- {
		h = mx[i](h)
	}
	return h
}
