package stage

import (
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/monitoring"
	"github.com/smartystreets/assertions"
	"github.com/go-resty/resty"
	"testing"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/association"
	"github.com/satori/go.uuid"
	"fmt"
)

var createdAssociation association.Association
var organisationId uuid.UUID
var serviceUserNumber string
var restClient *resty.Client

func The_Application_Is_Running(t *testing.T) {
	healthStatus := monitoring.HealthStatus{}

	restClient.NewRequest().
		SetResult(&healthStatus).
		Get("/v1/health")

	assertion := assertions.New(t)

	assertion.So(healthStatus.Status, assertions.ShouldEqual, "UP")
}

func An_Association_For_Bacs_Is_Created(t *testing.T) {
	organisationId = uuid.NewV4()
	serviceUserNumber = "123456"
	createAssociationCommand := association.CreateAssociationCommand{
		OrganisationId:    organisationId,
		ServiceUserNumber: serviceUserNumber,
	}

	restClient.NewRequest().
		SetResult(&createdAssociation).
		SetBody(createAssociationCommand).
		Post("/v1/association")
}

func The_Association_Can_Be_Retrieved_Correctly(t *testing.T) {
	var associations []association.Association
	restClient.NewRequest().
		SetResult(&associations).
		Get("/v1/association")

	assertion := assertions.New(t)

	fmt.Println(createdAssociation.Id)

	assertion.So(associations, assertions.ShouldHaveLength, 1)
	assertion.So(associations[0].Id, assertions.ShouldNotBeNil)
	assertion.So(associations[0].OrganisationId, assertions.ShouldEqual, organisationId)
	assertion.So(associations[0].Record.ServiceUserNumber, assertions.ShouldEqual, serviceUserNumber)

}

func ConfigureAssociationsStage(address string) {
	restClient = resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHostURL(address)
}

func ResetAssociationsStage() {
	createdAssociation = association.Association{}
}
