package microservices_example

import (
	"github.com/enriquerecarte/microservices-example/pkg/dockercompose"
	"os"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway"
	"testing"
	. "github.com/enriquerecarte/microservices-example/test/stage"
	"fmt"
	. "github.com/enriquerecarte/microservices-example/pkg/go-bdd"
)

func TestMain(m *testing.M) {
	startDockerCompose()

	//pact := test.StartPact()

	testAddress, stopServer := startServer()

	ConfigureAssociationsStage(testAddress)
	ConfigureOutboundPaymentsStage(testAddress)

	result := m.Run()

	consoleReporter := &ConsoleReporter{}
	ReportFeatureTo(consoleReporter)
	consoleReporter.Flush()

	//pact.Teardown()
	dockercompose.StopDockerCompose()
	stopServer <- true

	os.Exit(result)
}

func startServer() (string, chan bool) {
	port, _ := freeport.GetFreePort()
	testAddress := fmt.Sprintf("http://localhost:%d", port)
	viper.Set("server.port", port)
	serverReady := make(chan bool, 1)
	stopServer := make(chan bool, 1)
	go bacs_gateway.Start(stopServer, serverReady)
	<-serverReady
	return testAddress, stopServer
}

func startDockerCompose() {
	dockercompose.StartDockerCompose("docker-compose.yml", "bacs-gateway", []string{
		"VAULT_PORT",
		"POSTGRES_PORT",
		"SQS_PORT",
		"FAKE_ETS_PORT",
		"PACT_GO_PORT",
		"WIREMOCK_PORT",
	})
	viper.Set("db.port", os.Getenv("POSTGRES_PORT"))
	viper.Set("vault.address", fmt.Sprintf("http://localhost:%s", os.Getenv("VAULT_PORT")))
}
