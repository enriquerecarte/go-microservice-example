package microservices_example

import "testing"
import (
	. "github.com/enriquerecarte/go-bdd"
	. "github.com/enriquerecarte/microservices-example/test/stage"
)

func Test_CanAddAnAssociationCorrectly(t *testing.T) {
	TestThat(t).
		Given(The_Application_Is_Running).
		When(An_Association_For_Bacs_Is_Created).
		Then(The_Association_Can_Be_Retrieved_Correctly)
}
