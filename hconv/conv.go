package hconv

import (
	"net/http"

	"github.com/prasannavl/mchain"
)

func ToHttpFunc(h mchain.HandlerFunc, errorHandler mchain.ErrorHandler) http.HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)
		if err != nil && errorHandler != nil {
			errorHandler(err, w, r)
		}
	}
	return http.HandlerFunc(hh)
}

func ToHttp(h mchain.Handler, errorHandler mchain.ErrorHandler) http.Handler {
	hf := mchain.HandlerFunc(h.ServeHTTP)
	stdHf := ToHttpFunc(hf, errorHandler)
	return http.Handler(stdHf)
}

func FromHttpFunc(h http.HandlerFunc, recoverPanic bool) mchain.Handler {
	hh := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		if recoverPanic {
			defer mchain.RecoverIntoError(&err)
		}
		h.ServeHTTP(w, r)
		return err
	}
	return mchain.HandlerFunc(hh)
}

func FromHttp(h http.Handler, recoverPanic bool) mchain.Handler {
	hh := FromHttpFunc(http.HandlerFunc(h.ServeHTTP), recoverPanic)
	return mchain.Handler(hh)
}
