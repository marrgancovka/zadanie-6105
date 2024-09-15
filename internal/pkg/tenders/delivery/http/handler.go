package http

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
	"zadanie-6105/internal/pkg/tenders"
	"zadanie-6105/internal/pkg/utils"
)

type TenderHandler struct {
	u tenders.TenderUsecase
}

func NewHandler(u tenders.TenderUsecase) *TenderHandler {
	return &TenderHandler{u: u}
}

func (h *TenderHandler) GetTendersList(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ReadLimitOffset(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	serviceType := r.URL.Query()["service_type"]

	tendersList, err := h.u.GetTendersList(limit, offset, serviceType)
	if err != nil {
		if !errors.Is(err, myErrors.ErrBadRequest) {
			utils.WriteError(w, http.StatusInternalServerError, myErrors.ErrInternal)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tendersList)
}

func (h *TenderHandler) CreateNewTender(w http.ResponseWriter, r *http.Request) {
	var tenderData *models.TendersRequest
	if err := utils.ReadRequestData(r, &tenderData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	newTender, err := h.u.CreateNewTender(tenderData)
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

	utils.WriteJSON(w, http.StatusOK, newTender)

}

func (h *TenderHandler) GetUserTenders(w http.ResponseWriter, r *http.Request) {
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
	tendersList, err := h.u.GetUserTender(limit, offset, username)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tendersList)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	tenderId, err := uuid.FromString(tenderIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
		return
	}

	status, err := h.u.GetTenderStatus(tenderId, username)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, status)
}

func (h *TenderHandler) EditTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	tenderId, err := uuid.FromString(tenderIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	status := r.URL.Query().Get("status")
	username := r.URL.Query().Get("username")
	if username == "" || status == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	tender, err := h.u.EditTenderStatus(tenderId, username, status)
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
	utils.WriteJSON(w, http.StatusOK, tender)
}

func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	var editedData *models.TenderEditRequest
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	tenderId, err := uuid.FromString(tenderIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, myErrors.ErrBadRequest)
		return
	}
	if err = utils.ReadRequestData(r, &editedData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	tender, err := h.u.EditTender(tenderId, username, editedData)
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
	utils.WriteJSON(w, http.StatusOK, tender)
}
