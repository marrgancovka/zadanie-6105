package bids

import (
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
)

type BidRepository interface {
	CreateBid(bidData *models.BidRequest) (*models.BidResponse, error)
	SelectUserBids(limit, offset int32, username string) ([]*models.BidResponse, error)
	SelectTenderBids(limit, offset int32, tenderId uuid.UUID, username string) ([]*models.BidResponse, error)
	SelectBidStatus(bidId uuid.UUID) (string, error)
	CheckBidAuthor(bidId uuid.UUID, username string) (bool, error)
	UpdateBidStatus(bidId uuid.UUID, status string) (*models.BidResponse, error)
	UpdateBid(bidId uuid.UUID, editedData *models.BidEditRequest) (*models.BidResponse, error)
	SubmitDecision(bidId uuid.UUID, decision string, username string) (*models.BidResponse, error)
}

type BidUsecase interface {
	CreateNewBid(bidData *models.BidRequest) (*models.BidResponse, error)
	GetUserBids(limit, offset int32, username string) ([]*models.BidResponse, error)
	GetTenderBids(limit, offset int32, tenderId uuid.UUID, username string) ([]*models.BidResponse, error)
	GetBidStatus(bidId uuid.UUID, username string) (string, error)
	EditBidStatus(bidId uuid.UUID, username, status string) (*models.BidResponse, error)
	EditBid(bidId uuid.UUID, username string, editedData *models.BidEditRequest) (*models.BidResponse, error)
	SubmitDecision(bidId uuid.UUID, username string, decision string) (*models.BidResponse, error)
}
