package mconv

import (
	"net/http"

	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/hconv"
)

func FromSimple(fn mchain.SimpleMiddleware) mchain.Middleware {
	m := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			return fn(w, r, next)
		}
		return mchain.HandlerFunc(f)
	}
	return m
}

func ToSimple(m mchain.Middleware) mchain.SimpleMiddleware {
	h := func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
		return m(next).ServeHTTP(w, r)
	}
	return h
}

func HttpFromSimple(fn mchain.HttpSimpleMiddleware) mchain.HttpMiddleware {
	m := func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, next)
		}
		return http.HandlerFunc(f)
	}
	return m
}

func HttpToSimple(m mchain.HttpMiddleware) mchain.HttpSimpleMiddleware {
	h := func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		m(next).ServeHTTP(w, r)
	}
	return h
}

func ToHttp(m mchain.Middleware, errorHandler mchain.ErrorHandler, recoverPanic bool) mchain.HttpMiddleware {
	hh := func(hx http.Handler) http.Handler {
		handler := hconv.FromHttp(hx, recoverPanic)
		return hconv.ToHttp(m(handler), errorHandler)
	}
	return hh
}

func FromHttp(h mchain.HttpMiddleware, innerErrorHandler mchain.ErrorHandler, recoverPanic bool) mchain.Middleware {
	hh := func(hx mchain.Handler) mchain.Handler {
		httpHandler := hconv.ToHttp(hx, innerErrorHandler)
		nextHttpHandler := h(httpHandler)
		return hconv.FromHttp(nextHttpHandler, recoverPanic)
	}
	return hh
}
