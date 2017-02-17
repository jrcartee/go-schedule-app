package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jrcartee/scheduling/httputil"
)

func PanicRecovery() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					var err error
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = fmt.Errorf("PanicRecovery -- Unexpected Type:\n%+v\n", t)
					}
					httputil.InternalError(w, err)
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

func RequestLogger(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}
