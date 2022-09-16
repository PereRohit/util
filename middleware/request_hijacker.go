package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

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
		reqId := r.Header.Get(constant.RequestIdHeader)
		if reqId == "" {
			reqId = uuid.NewString()
			r.Header.Set(constant.RequestIdHeader, reqId)
		}
		rT := *r
		hijackedWriter := &respWriterWithStatus{-1, "", w}

		start := time.Now()
		next.ServeHTTP(hijackedWriter, r)
		w.Header().Set("user-agent", constant.UserAgentSvc)
		end := time.Now().Sub(start)

		log.WithNoCaller().Info(fmt.Sprintf("%20s | %-6s | %-25s | %d | %10s | %s:%s | %s",
			rT.RemoteAddr, rT.Method, rT.URL.String(), hijackedWriter.status, end.String(), constant.RequestIdHeader, reqId, hijackedWriter.response))
	})
}
