package mchain

import "net/http"

type HttpMiddleware func(http.Handler) http.Handler
type SimpleHttpMiddleware func(w http.ResponseWriter, r *http.Request, next http.Handler)

type HttpChain struct {
	Middlewares []HttpMiddleware
}

type HttpChainBuilder struct {
	chain   HttpChain
	handler http.Handler
}

func CreateHttpBuilder(middlewares ...HttpMiddleware) HttpChainBuilder {
	return HttpChainBuilder{HttpChain{middlewares}, nil}
}

func (b HttpChainBuilder) Add(m ...HttpMiddleware) HttpChainBuilder {
	b.chain.Middlewares = append(b.chain.Middlewares, m...)
	return b
}

func (b HttpChainBuilder) AddSimple(m ...SimpleHttpMiddleware) HttpChainBuilder {
	s := make([]HttpMiddleware, 0, len(m))
	for _, x := range m {
		s = append(s, CreateHttpMiddleware(x))
	}
	b.chain.Middlewares = append(b.chain.Middlewares, s...)
	return b
}

func (b HttpChainBuilder) Handler(finalHandler http.Handler) HttpChainBuilder {
	b.handler = finalHandler
	return b
}

func (b HttpChainBuilder) Build() http.Handler {
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
