package middleware

import (
	"net/http"

	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/response"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response.ToJson(w, http.StatusInternalServerError, "Oops! Something went wrong.", nil)
				log.WithNoCaller().Warn("Panic occurred:", err)
				if e, ok := err.(error); ok {
					log.WithNoCaller().Warn("Panic occurred:", e.Error())
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
