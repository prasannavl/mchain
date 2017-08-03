package mchain

import "net/http"

type HttpMiddleware func(http.Handler) http.Handler

type HttpChain struct {
	Middlewares []HttpMiddleware
}

type HttpChainBuilder struct {
	chain   HttpChain
	handler http.Handler
}

func NewHttpBuilder(middlewares ...HttpMiddleware) HttpChainBuilder {
	return HttpChainBuilder{HttpChain{middlewares}, nil}
}

func (b HttpChainBuilder) Add(m ...HttpMiddleware) HttpChainBuilder {
	c := b.chain
	c.Middlewares = append(c.Middlewares, m...)
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
