package association

import (
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/infrastructure/events"
)

type CreateAssociationCommand struct {
	OrganisationId    uuid.UUID `sql:"organisationid"`
	ServiceUserNumber string
}

type DeleteAssociationCommand struct {
	Id uuid.UUID
}

func HandleCreateAssociationCommand(command CreateAssociationCommand) Association {
	logrus.Info("HandleCreateAssociationCommand:", command)
	id, _ := uuid.NewV4()
	association := Association{
		Id:             id,
		OrganisationId: command.OrganisationId,
		Record: AssociationRecord{
			ServiceUserNumber: command.ServiceUserNumber,
		},
		Version:   1,
		IsDeleted: false,
		IsLocked:  false,
	}

	Save(&association)
	events.Publish(AssociationCreatedEvent{
		CreatedAssociation: association,
	})

	return association
}

func HandleDeleteAssociationCommand(command DeleteAssociationCommand) {
	Delete(command.Id)
}
