package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/log"
)

func Run(r http.Handler, svrConfig config.ServerConfig) {
	// create channel to gracefully stop server
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	port, found := os.LookupEnv(constant.SERVER_PORT_ENV)
	if !found || port == "" {
		port = constant.DEFAULT_SERVER_PORT
	}
	if svrConfig.Port != "" {
		port = svrConfig.Port
	}
	host, found := os.LookupEnv(constant.SERVER_HOST_ENV)
	if !found || port == "" {
		host = constant.DEFAULT_SERVER_HOST
	}
	if svrConfig.Host != "" {
		host = svrConfig.Host
	}

	// set log data
	log.SetLogLevel(svrConfig.LogLevel)
	logSetter := log.GetStaticDataSetter()
	if svrConfig.Name != "" {
		logSetter.Add("service", svrConfig.Name)
	}
	if svrConfig.Version != "" {
		logSetter.Add("version", svrConfig.Version)
	}
	logSetter.Set()

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}
	var srvStop sync.WaitGroup
	defer srvStop.Wait()

	go func() {
		defer srvStop.Done()

		// wait for termination signal
		<-sc

		log.WithNoCaller().Info("Closing Server")
		err := s.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		log.WithNoCaller().Info("Server Closed!!")
	}()

	srvStop.Add(1)
	log.WithNoCaller().Info(fmt.Sprintf("Starting server(%s)", s.Addr))
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithNoCaller().Error("Server error::", err.Error())
	}
}
