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

func ToSimple(middleware mchain.Middleware) (fn func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error) {
	h := func(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
		hx := mchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			return next.ServeHTTP(w, r)
		})
		return middleware(hx).ServeHTTP(w, r)
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

func ToHttpSimple(middleware mchain.HttpMiddleware) (fn func(w http.ResponseWriter, r *http.Request, next http.Handler)) {
	h := func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		hx := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
		middleware(hx).ServeHTTP(w, r)
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
		handler := hconv.ToHttp(hx, innerErrorHandler)
		return hconv.FromHttp(handler)
	}
	return mchain.Middleware(hh)
}

func FromHttpRecoverable(h func(http.Handler) http.Handler, innerErrorHandler func(error)) mchain.Middleware {
	hh := func(hx mchain.Handler) mchain.Handler {
		handler := hconv.ToHttp(hx, innerErrorHandler)
		return hconv.FromHttpRecoverable(handler)
	}
	return mchain.Middleware(hh)
}
