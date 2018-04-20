package association

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/go-pg/pg"
	"github.com/enriquerecarte/microservices-example/internal/app/bacs-gateway/infrastructure/database"
)

func Save(association *Association) {
	err := database.Connection().Insert(association)
	if err != nil {
		panic(err)
	}
}

func Get(associationId uuid.UUID) Association {
	association := Association{
		Id: associationId,
	}

	err := database.Connection().Select(&association)
	if err != nil {
		panic(err)
	}

	return association
}

func GetAll() []Association {
	var associations = make([]Association, 0)

	err := database.Connection().Model(&associations).Select()
	if err != nil {
		panic(err)
	}

	return associations
}

func Delete(associationId uuid.UUID) {
	association := Association{
		Id: associationId,
	}

	result, err := database.Connection().Model(&association).Delete()

	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted:", result.RowsAffected())

}

func DeleteAll() {
	allAssociations := GetAll()

	var idsToDelete = make([]uuid.UUID, 0)
	for _, association := range allAssociations {
		idsToDelete = append(idsToDelete, association.Id)
	}

	result, err := database.Connection().Model(&Association{}).Where("id IN (?)", pg.In(idsToDelete)).Delete()

	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted:", result.RowsAffected())
}
