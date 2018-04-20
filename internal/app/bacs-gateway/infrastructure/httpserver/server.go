package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/association"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/monitoring"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/configuration"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/outboundpayments"
	"fmt"
	log "github.com/sirupsen/logrus"
	"context"
	"net/http"
)

func StartServer(stopServerListener <-chan bool, serverReadySender chan<- bool) {
	router := gin.Default()

	associationHandlers := router.Group("/v1/association")
	{
		associationHandlers.GET("", association.HandleGetAllAssociations)
		associationHandlers.GET("/:id", association.HandleGetAssociation)
		associationHandlers.POST("", association.HandleCreateAssociation)
		associationHandlers.DELETE("", association.HandleDeleteAll)
		associationHandlers.DELETE("/:id", association.HandleDelete)
	}
	paymentHandlers := router.Group("/v1/payments")
	{
		paymentHandlers.GET("", outboundpayments.HandleGetPayment)
	}

	monitoringHandlers := router.Group("/v1/")
	{
		monitoringHandlers.GET("/health", monitoring.HandleHealthCheck)
		monitoringHandlers.GET("/env", monitoring.HandleConfigurationCheck)
		monitoringHandlers.GET("/secret", monitoring.HandleSecretsCheck)
	}

	port := configuration.Get("server.port")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	go func() {
		<-stopServerListener
		log.Info("Shutting down")
		server.Shutdown(context.Background())
	}()

	fmt.Println(fmt.Sprintf("listening on localhost:%s", port))
	serverReadySender <- true
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
