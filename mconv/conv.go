package mconv

import (
	"net/http"

	"github.com/prasannavl/mchain"
	"github.com/prasannavl/mchain/hconv"
)

func FromSimple(fn func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error) mchain.Middleware {
	m := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			return fn(w, r, next)
		}
		return mchain.HandlerFunc(f)
	}
	return m
}

func ToSimple(m mchain.Middleware) (fn func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error) {
	h := func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
		return m(next).ServeHTTP(w, r)
	}
	return h
}

func HttpFromSimple(fn func(w http.ResponseWriter, r *http.Request, next http.Handler)) mchain.HttpMiddleware {
	m := func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, next)
		}
		return http.HandlerFunc(f)
	}
	return m
}

func ToHttpSimple(m mchain.HttpMiddleware) (fn func(w http.ResponseWriter, r *http.Request, next http.Handler)) {
	h := func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		m(next).ServeHTTP(w, r)
	}
	return h
}

func ToHttp(m mchain.Middleware, errorHandler func(error)) func(http.Handler) http.Handler {
	hh := func(hx http.Handler) http.Handler {
		handler := hconv.FromHttp(hx)
		return hconv.ToHttp(m(handler), errorHandler)
	}
	return hh
}

func FromHttp(h func(http.Handler) http.Handler, innerErrorHandler func(error)) mchain.Middleware {
	hh := func(hx mchain.Handler) mchain.Handler {
		httpHandler := hconv.ToHttp(hx, innerErrorHandler)
		nextHttpHandler := h(httpHandler)
		return hconv.FromHttp(nextHttpHandler)
	}
	return mchain.Middleware(hh)
}

func FromHttpRecoverable(h func(http.Handler) http.Handler, innerErrorHandler func(error)) mchain.Middleware {
	hh := func(hx mchain.Handler) mchain.Handler {
		httpHandler := hconv.ToHttp(hx, innerErrorHandler)
		nextHttpHandler := h(httpHandler)
		return hconv.FromHttpRecoverable(nextHttpHandler)
	}
	return mchain.Middleware(hh)
}
