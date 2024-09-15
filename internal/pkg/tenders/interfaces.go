package tenders

import (
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
)

type TenderRepoPostgres interface {
	SelectTendersList(limit, offset int32, serviceType []string) ([]*models.TendersResponse, error)
	CheckUsernameOrganization(creatorUsername string, organizationId uuid.UUID) (bool, error)
	CheckUsernameTender(username string, tenderId uuid.UUID) (bool, error)
	CreateTender(tenderData *models.TendersRequest) (*models.TendersResponse, error)
	SelectUserTenders(limit, offset int32, username string) ([]*models.TendersResponse, error)
	SelectTenderStatus(tenderId uuid.UUID) (string, error)
	EditStatusTender(tenderId uuid.UUID, status string) (*models.TendersResponse, error)
	EditTender(tenderId uuid.UUID, editedData *models.TenderEditRequest) (*models.TendersResponse, error)
}

type TenderUsecase interface {
	GetTendersList(limit, offset int32, serviceType []string) ([]*models.TendersResponse, error)
	CreateNewTender(tenderData *models.TendersRequest) (*models.TendersResponse, error)
	GetUserTender(limit, offset int32, username string) ([]*models.TendersResponse, error)
	GetTenderStatus(tenderId uuid.UUID, username string) (string, error)
	EditTenderStatus(tenderId uuid.UUID, username, status string) (*models.TendersResponse, error)
	EditTender(tenderId uuid.UUID, username string, editedData *models.TenderEditRequest) (*models.TendersResponse, error)
}
