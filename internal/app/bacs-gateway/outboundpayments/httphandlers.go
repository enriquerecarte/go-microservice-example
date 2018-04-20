package outboundpayments

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
)

type PaymentSubmission struct {
	Data PaymentSubmissionData `json:"data"`
}

type PaymentSubmissionData struct {
	Id string `json:"id"`
}

func HandleGetPayment(c *gin.Context) {
	paymentSubmission := PaymentSubmission{}

	host := fmt.Sprintf("http://localhost:%s", os.Getenv("WIREMOCK_PORT"))

	logrus.Info("Host", host)

	resty.New().R().SetResult(&paymentSubmission).
		SetHeader("Accept", "application/json").
		Get(host + "/v1/payments/9add65d3-272a-43ab-9693-27244c3b433d/submissions/a3f7b511-29ea-4e90-9753-a6afd757844f")

	logrus.Info("Sending response:", paymentSubmission)

	c.JSON(200, paymentSubmission)
}
