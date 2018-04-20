package configuration

import (
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"fmt"
)

var vaultApi *vaultapi.Logical

func GetAllSecrets() map[string]interface{} {
	applicationName := viper.GetString("application.name")
	secret, err := vaultApi.Read(fmt.Sprintf("secret/%s/", applicationName))

	if err != nil {
		panic(err)
	}

	return secret.Data

}

func initVault() {
	config := vaultapi.DefaultConfig()
	config.Address = GetOrDefault("vault.address", "http://localhost:8200")
	client, err := vaultapi.NewClient(config)
	if err != nil {
		panic(err)
	}
	client.SetToken(viper.GetString("vault.token"))
	vaultApi = client.Logical()
}
