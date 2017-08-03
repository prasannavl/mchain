package mchain

import "net/http"

type ChainBuilder struct {
	chain   Chain
	handler Handler
}

func NewBuilder(middlewares ...Middleware) ChainBuilder {
	return ChainBuilder{Chain{middlewares}, nil}
}

func (b ChainBuilder) Add(m ...Middleware) ChainBuilder {
	c := b.chain
	c.Middlewares = append(c.Middlewares, m...)
	return b
}

func (b ChainBuilder) Handler(finalHandler Handler) ChainBuilder {
	b.handler = finalHandler
	return b
}

func (b ChainBuilder) Build() Handler {
	h := b.handler
	if h == nil {
		h = HandlerFromHttp(http.DefaultServeMux)
	}
	c := b.chain
	mx := c.Middlewares
	mLen := len(mx)
	for i := range mx {
		h = c.Middlewares[mLen-1-i](h)
	}
	return h
}

func (b ChainBuilder) BuildHttp(errorHandler func(error)) http.Handler {
	return HandlerToHttp(b.Build(), errorHandler)
}
