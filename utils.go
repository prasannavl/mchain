package mchain

import (
	"errors"
	"net/http"
)

func HandlerFuncToHttp(h HandlerFunc, errorHandler func(error)) http.HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)
		if err != nil && errorHandler != nil {
			errorHandler(err)
		}
	}
	return http.HandlerFunc(hh)
}

func HandlerToHttp(h Handler, errorHandler func(error)) http.Handler {
	hf := HandlerFunc(h.ServeHTTP)
	stdHf := HandlerFuncToHttp(hf, errorHandler)
	return http.Handler(stdHf)
}

func HandlerFuncFromHttp(h http.HandlerFunc) HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	}
	return HandlerFunc(hh)
}

func HandlerFuncFromHttpRecoverable(h http.HandlerFunc) HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer RecoverIntoError(&err)
		h.ServeHTTP(w, r)
		return err
	}
	return HandlerFunc(hh)
}

func HandlerFromHttp(h http.Handler) Handler {
	hh := HandlerFuncFromHttp(http.HandlerFunc(h.ServeHTTP))
	return Handler(hh)
}

func HandlerFromHttpRecoverable(h http.Handler) Handler {
	hh := HandlerFuncFromHttpRecoverable(http.HandlerFunc(h.ServeHTTP))
	return Handler(hh)
}

func MiddlewareToHttp(m Middleware, errorHandler func(error)) func(http.Handler) http.Handler {
	hh := func(hx http.Handler) http.Handler {
		handler := HandlerFromHttp(hx)
		return HandlerToHttp(m(handler), errorHandler)
	}
	return hh
}

func MiddlewareFromHttp(h func(http.Handler) http.Handler, innerErrorHandler func(error)) Middleware {
	hh := func(hx Handler) Handler {
		handler := HandlerToHttp(hx, innerErrorHandler)
		return HandlerFromHttp(handler)
	}
	return Middleware(hh)
}

func MiddlewareFromHttpRecoverable(h func(http.Handler) http.Handler, innerErrorHandler func(error)) Middleware {
	hh := func(hx Handler) Handler {
		handler := HandlerToHttp(hx, innerErrorHandler)
		return HandlerFromHttpRecoverable(handler)
	}
	return Middleware(hh)
}

func RecoverIntoError(errPointer *error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case error:
			*errPointer = x
		case string:
			*errPointer = errors.New(x)
		default:
			*errPointer = errors.New("unknown panic")
		}
	}
}
