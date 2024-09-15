package usecase

import (
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
	"zadanie-6105/internal/pkg/tenders"
)

type TenderUsecase struct {
	r tenders.TenderRepoPostgres
}

func NewUsecase(r tenders.TenderRepoPostgres) *TenderUsecase {
	return &TenderUsecase{
		r: r,
	}
}

func (u *TenderUsecase) GetTendersList(limit, offset int32, serviceType []string) ([]*models.TendersResponse, error) {
	tendersList, err := u.r.SelectTendersList(limit, offset, serviceType)
	if err != nil {
		return nil, err
	}
	return tendersList, nil
}
func (u *TenderUsecase) CreateNewTender(tenderData *models.TendersRequest) (*models.TendersResponse, error) {
	ok, err := u.r.CheckUsernameOrganization(tenderData.CreatorUsername, tenderData.OrganizationId)
	if err != nil {
		return nil, myErrors.ErrBadRequest
	}
	if !ok {
		return nil, myErrors.ErrForbidden
	}
	newTender, err := u.r.CreateTender(tenderData)
	if err != nil {
		return nil, err
	}
	return newTender, nil
}
func (u *TenderUsecase) GetUserTender(limit, offset int32, username string) ([]*models.TendersResponse, error) {
	userTenders, err := u.r.SelectUserTenders(limit, offset, username)
	if err != nil {
		return nil, err
	}
	return userTenders, nil
}
func (u *TenderUsecase) GetTenderStatus(tenderId uuid.UUID, username string) (string, error) {
	ok, err := u.r.CheckUsernameTender(username, tenderId)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", myErrors.ErrUserNotFound
	}
	status, err := u.r.SelectTenderStatus(tenderId)
	if err != nil {
		return "", err
	}
	return status, nil
}
func (u *TenderUsecase) EditTenderStatus(tenderId uuid.UUID, username, status string) (*models.TendersResponse, error) {
	ok, err := u.r.CheckUsernameTender(username, tenderId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, myErrors.ErrUserNotFound
	}
	editedTender, err := u.r.EditStatusTender(tenderId, status)
	if err != nil {
		return nil, err
	}
	return editedTender, nil
}
func (u *TenderUsecase) EditTender(tenderId uuid.UUID, username string, editedData *models.TenderEditRequest) (*models.TendersResponse, error) {
	ok, err := u.r.CheckUsernameTender(username, tenderId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, myErrors.ErrUserNotFound
	}
	editedTender, err := u.r.EditTender(tenderId, editedData)
	if err != nil {
		return nil, err
	}
	return editedTender, nil
}
