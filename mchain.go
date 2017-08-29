// Package mchain provides a convenient way for middleware composition
package mchain

import (
	"net/http"
)

type Middleware func(Handler) Handler
type SimpleMiddleware func(http.ResponseWriter, *http.Request, Handler) error

type Chain struct {
	Middlewares []Middleware
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
