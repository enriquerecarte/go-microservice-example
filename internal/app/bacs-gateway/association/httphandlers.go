package association

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"fmt"
)

func HandleCreateAssociation(c *gin.Context) {
	createAssociationCommand := CreateAssociationCommand{}
	c.BindJSON(&createAssociationCommand)

	createdAssociation := HandleCreateAssociationCommand(createAssociationCommand)

	c.JSON(200, createdAssociation)
}

func HandleGetAllAssociations(c *gin.Context) {
	var associations []Association

	sun := c.Query("sun")
	org := c.Query("org")
	if sun != "" {
		fmt.Println("Filtering out sun", sun)
		associations = GetByServiceUserNumber(sun)
	} else if org != "" {
		fmt.Println("Filtering out org", org)
		organisationId, _ := uuid.FromString(org)
		associations = GetByOrganisationId(organisationId)
	} else {
		associations = GetAll()
	}

	c.JSON(200, associations)
}

func HandleGetAssociation(c *gin.Context) {
	associationId, _ := uuid.FromString(c.Param("id"))

	associations := Get(associationId)

	c.JSON(200, associations)
}

func HandleDelete(c *gin.Context) {
	associationId, _ := uuid.FromString(c.Param("id"))
	command := DeleteAssociationCommand{
		Id: associationId,
	}

	HandleDeleteAssociationCommand(command)

	c.Status(201)
}

func HandleDeleteAll(c *gin.Context) {
	DeleteAll()

	c.Status(201)
}
