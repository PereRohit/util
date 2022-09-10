package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/log"
)

type respWriterWithStatus struct {
	status   int
	response string
	http.ResponseWriter
}

func (w *respWriterWithStatus) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *respWriterWithStatus) Write(d []byte) (int, error) {
	w.response = string(d)
	return w.ResponseWriter.Write(d)
}

func RequestHijacker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rT := *r
		hijackedWriter := &respWriterWithStatus{-1, "", w}

		start := time.Now()
		next.ServeHTTP(hijackedWriter, r)
		w.Header().Set("user-agent", constant.UserAgentSvc)
		end := time.Now().Sub(start)

		log.Info(fmt.Sprintf("%20s | %5s | %20s | %d | %10s | %s", rT.RemoteAddr, rT.Method, rT.URL.String(), hijackedWriter.status, end.String(), hijackedWriter.response))
	})
}
