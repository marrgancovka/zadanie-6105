package usecase

import (
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
	"zadanie-6105/internal/pkg/bids"
)

type BidUsecase struct {
	r bids.BidRepository
}

func NewBidUsecase(r bids.BidRepository) *BidUsecase {
	return &BidUsecase{r: r}
}

func (u *BidUsecase) CreateNewBid(bidData *models.BidRequest) (*models.BidResponse, error) {
	bid, err := u.r.CreateBid(bidData)
	if err != nil {
		return nil, err
	}
	return bid, nil
}
func (u *BidUsecase) GetUserBids(limit, offset int32, username string) ([]*models.BidResponse, error) {
	userBids, err := u.r.SelectUserBids(limit, offset, username)
	if err != nil {
		return nil, err
	}
	return userBids, nil
}
func (u *BidUsecase) GetTenderBids(limit, offset int32, tenderId uuid.UUID, username string) ([]*models.BidResponse, error) {
	tenderBids, err := u.r.SelectTenderBids(limit, offset, tenderId, username)
	if err != nil {
		return nil, err
	}
	return tenderBids, nil
}
func (u *BidUsecase) GetBidStatus(bidId uuid.UUID, username string) (string, error) {
	ok, err := u.r.CheckBidAuthor(bidId, username)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", myErrors.ErrForbidden
	}
	status, err := u.r.SelectBidStatus(bidId)
	if err != nil {
		return "", err
	}
	return status, nil
}
func (u *BidUsecase) EditBidStatus(bidId uuid.UUID, username, status string) (*models.BidResponse, error) {
	ok, err := u.r.CheckBidAuthor(bidId, username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, myErrors.ErrForbidden
	}
	bid, err := u.r.UpdateBidStatus(bidId, status)
	if err != nil {
		return nil, err
	}
	return bid, nil
}
func (u *BidUsecase) EditBid(bidId uuid.UUID, username string, editedData *models.BidEditRequest) (*models.BidResponse, error) {
	ok, err := u.r.CheckBidAuthor(bidId, username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, myErrors.ErrForbidden
	}
	bid, err := u.r.UpdateBid(bidId, editedData)
	if err != nil {
		return nil, err
	}
	return bid, nil
}
func (u *BidUsecase) SubmitDecision(bidId uuid.UUID, username string, decision string) (*models.BidResponse, error) {
	bid, err := u.r.SubmitDecision(bidId, decision, username)
	if err != nil {
		return nil, err
	}
	return bid, nil
}
