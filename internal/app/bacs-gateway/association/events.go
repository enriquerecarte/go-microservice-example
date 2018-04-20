package association

import "github.com/sirupsen/logrus"

type AssociationCreatedEvent struct {
	CreatedAssociation Association
}

func HandleAssociationCreatedEvent(event AssociationCreatedEvent) {
	logrus.Info("AssociationCreatedEvent:", event)
}
