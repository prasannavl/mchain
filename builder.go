package mchain

import (
	"net/http"

	"github.com/prasannavl/mchain/hconv"
	"github.com/prasannavl/mchain/mconv"
)

type ChainBuilder struct {
	chain   Chain
	handler Handler
}

func CreateBuilder(middlewares ...Middleware) ChainBuilder {
	return ChainBuilder{Chain{middlewares}, nil}
}

func (b ChainBuilder) Add(m ...Middleware) ChainBuilder {
	b.chain.Middlewares = append(b.chain.Middlewares, m...)
	return b
}

func (b ChainBuilder) AddSimple(m ...SimpleMiddleware) ChainBuilder {
	s := make([]Middleware, 0, len(m))
	for _, x := range m {
		s = append(s, mconv.From(x))
	}
	b.chain.Middlewares = append(b.chain.Middlewares, s...)
	return b
}

func (b ChainBuilder) Handler(finalHandler Handler) ChainBuilder {
	b.handler = finalHandler
	return b
}

func (b ChainBuilder) Build() Handler {
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

func (b ChainBuilder) BuildHttp(errorHandler func(error)) http.Handler {
	return hconv.ToHttp(b.Build(), errorHandler)
}
