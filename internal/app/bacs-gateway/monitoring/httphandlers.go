package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/configuration"
)

type HealthStatus struct {
	Status string `json:"status"`
}

func HandleHealthCheck(c *gin.Context) {
	c.JSON(200, HealthStatus{Status: "UP"})
}

func HandleConfigurationCheck(c *gin.Context) {
	allKeys := viper.AllKeys()
	config := make(map[string]string)

	for _, key := range allKeys {
		config[key] = viper.GetString(key)
	}
	c.JSON(200, config)
}

func HandleSecretsCheck(c *gin.Context) {
	c.JSON(200, configuration.GetAllSecrets())
}
