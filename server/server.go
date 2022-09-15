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
	logStaticData := ""
	if svrConfig.Name != "" {
		logStaticData += "name: " + svrConfig.Name
	}
	if svrConfig.Version != "" {
		logStaticData += " version: " + svrConfig.Version
	}
	log.SetStaticData(logStaticData)

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

		log.Info("Closing Server")
		err := s.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		log.Info("Server Closed!!")
	}()

	srvStop.Add(1)
	log.Info(fmt.Sprintf("Starting server(%s)", s.Addr))
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Server error::", err.Error())
	}
}
