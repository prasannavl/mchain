package hconv

import (
	"net/http"

	"github.com/prasannavl/mchain"
)

func FuncToHttp(h mchain.HandlerFunc, errorHandler mchain.ErrorHandler) http.HandlerFunc {
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
	stdHf := FuncToHttp(hf, errorHandler)
	return http.Handler(stdHf)
}

func FuncFromHttp(h http.HandlerFunc) mchain.HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	}
	return mchain.HandlerFunc(hh)
}

func FuncFromHttpRecoverable(h http.HandlerFunc) mchain.HandlerFunc {
	hh := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer mchain.RecoverIntoError(&err)
		h.ServeHTTP(w, r)
		return err
	}
	return mchain.HandlerFunc(hh)
}

func FromHttp(h http.Handler) mchain.Handler {
	hh := FuncFromHttp(http.HandlerFunc(h.ServeHTTP))
	return mchain.Handler(hh)
}

func FromHttpRecoverable(h http.Handler) mchain.Handler {
	hh := FuncFromHttpRecoverable(http.HandlerFunc(h.ServeHTTP))
	return mchain.Handler(hh)
}
