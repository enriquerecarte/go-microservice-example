package association

import "github.com/satori/go.uuid"

type Association struct {
	Id             uuid.UUID
	OrganisationId uuid.UUID `sql:"organisationid"`
	Record         AssociationRecord
	Version        int
	IsDeleted      bool      `sql:"isdeleted,notnull"`
	IsLocked       bool      `sql:"islocked,notnull"`
}

type AssociationRecord struct {
	ServiceUserNumber string `json:"serviceUserNumber"`
}

