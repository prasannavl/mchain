package builder

import (
	"net/http"

	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/mconv"
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

func (b *HttpChainBuilder) AddSimple(m ...mchain.SimpleHttpMiddleware) *HttpChainBuilder {
	s := make([]mchain.HttpMiddleware, 0, len(m))
	for _, x := range m {
		s = append(s, mconv.HttpFrom(x))
	}
	b.chain.Middlewares = append(b.chain.Middlewares, s...)
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
	mLen := len(mx)
	for i := range mx {
		h = c.Middlewares[mLen-1-i](h)
	}
	return h
}
