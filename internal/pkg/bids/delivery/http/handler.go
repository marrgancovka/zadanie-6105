package http

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
	"zadanie-6105/internal/pkg/bids"
	"zadanie-6105/internal/pkg/utils"
)

type BidHandler struct {
	u bids.BidUsecase
}

func NewHandler(u bids.BidUsecase) *BidHandler {
	return &BidHandler{u: u}
}

func (h *BidHandler) CreateNewBid(w http.ResponseWriter, r *http.Request) {
	var bidData *models.BidRequest
	if err := utils.ReadRequestData(r, &bidData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	newBid, err := h.u.CreateNewBid(bidData)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrTenderNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrTenderNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, newBid)
}

func (h *BidHandler) GetUserBids(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ReadLimitOffset(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	bidsList, err := h.u.GetUserBids(limit, offset, username)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, bidsList)
}

func (h *BidHandler) GetTenderBids(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ReadLimitOffset(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	tenderId, err := uuid.FromString(tenderIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	bidsList, err := h.u.GetTenderBids(limit, offset, tenderId, username)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrTenderNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrTenderNotFound)
			return
		case errors.Is(err, myErrors.ErrBidNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrBidNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, bidsList)
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]
	bidId, err := uuid.FromString(bidIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	status, err := h.u.GetBidStatus(bidId, username)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrBidNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrBidNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, status)
}

func (h *BidHandler) EditBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]
	bidId, err := uuid.FromString(bidIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	status := r.URL.Query().Get("status")
	if username == "" || status == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	editedBid, err := h.u.EditBidStatus(bidId, username, status)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrBidNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrBidNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, editedBid)
}
func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]
	bidId, err := uuid.FromString(bidIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	var editedData *models.BidEditRequest
	if err = utils.ReadRequestData(r, &editedData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	editedBid, err := h.u.EditBid(bidId, username, editedData)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrBidNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrBidNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, editedBid)
}

func (h *BidHandler) SubmitDecision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]
	bidId, err := uuid.FromString(bidIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	decision := r.URL.Query().Get("decision")
	if username == "" || decision == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	bid, err := h.u.SubmitDecision(bidId, username, decision)
	if err != nil {
		switch {
		case errors.Is(err, myErrors.ErrBadRequest):
			utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
			return
		case errors.Is(err, myErrors.ErrForbidden):
			utils.WriteError(w, http.StatusForbidden, myErrors.ErrForbidden)
			return
		case errors.Is(err, myErrors.ErrUserNotFound):
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		case errors.Is(err, myErrors.ErrBidNotFound):
			utils.WriteError(w, http.StatusNotFound, myErrors.ErrBidNotFound)
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, bid)
}
