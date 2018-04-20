package main

import (
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway"
	"os"
	"os/signal"
)

func main() {
	stopServerSender := make(chan bool, 1)
	serverReadyListener := make(chan bool, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		stopServerSender <- <-stop == os.Interrupt
	}()

	bacs_gateway.Start(stopServerSender, serverReadyListener)
}
