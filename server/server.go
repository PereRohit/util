package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/log"
)

func Run(r http.Handler) {
	// create channel to gracefully stop server
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	port, found := os.LookupEnv(constant.SERVER_PORT_ENV)
	if !found || port == "" {
		port = constant.DEFAULT_SERVER_PORT
	}
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "", port),
		Handler: r,
	}

	go func() {
		// wait for termination signal
		<-sc

		log.Info("Closing Server")
		err := s.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		log.Info("Server Closed!!")
	}()

	log.Info(fmt.Sprintf("Starting server(%s)", s.Addr))
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Server error::", err.Error())
	}
}

