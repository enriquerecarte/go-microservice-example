package test

import (
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/spf13/viper"
)

var pact *dsl.Pact

func StartPact() *dsl.Pact {
	getenv := viper.GetInt("PACT_GO_PORT")
	pact = &dsl.Pact{
		Port:     getenv, // Ensure this port matches the daemon port!
		Consumer: "bacsgateway",
		Provider: "paymentapi",
		Host:     "localhost",
		PactDir:  "build/pacts/",
	}

	pact.Setup(true)

	return pact
}
