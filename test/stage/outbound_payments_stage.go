package stage

import (
	"github.com/go-resty/resty"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/outboundpayments"
	"github.com/smartystreets/assertions"
	"testing"
)

var outboundPaymentsRestClient *resty.Client
var paymentSubmission outboundpayments.PaymentSubmission

func I_request_a_payment(t *testing.T) {
	restClient.NewRequest().
		SetResult(&paymentSubmission).
		Get("/v1/payments")
}

func The_payment_is_returned_successfully(t *testing.T) {
	assertion := assertions.New(t)

	assertion.So(paymentSubmission.Data.Id, assertions.ShouldEqual, "")
}

func ConfigureOutboundPaymentsStage(address string) {
	outboundPaymentsRestClient = resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHostURL(address)
}
