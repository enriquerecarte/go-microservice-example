package microservices_example

import (
	"testing"
	. "github.com/enriquerecarte/go-bdd"
	. "github.com/enriquerecarte/microservices-example/test/stage"
)

func Test_GetPayment(t *testing.T) {
	TestThat(t).
		When(I_request_a_payment).
		Then(The_payment_is_returned_successfully)
}
