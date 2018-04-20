package bacs_gateway

import (
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/configuration"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/infrastructure/database"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/infrastructure/httpserver"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/infrastructure/events"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/association"
)

func Start(stopServer <-chan bool, serverReady chan<- bool) {
	configuration.Start()
	database.StartDatabase()
	events.StartEvents()
	events.Subscribe(association.HandleAssociationCreatedEvent)
	httpserver.StartServer(stopServer, serverReady)
}
