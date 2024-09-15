package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type MessageResponse struct {
	Message string `json:"reason"`
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	resp, err := json.Marshal(MessageResponse{Message: err.Error()})
	if err != nil {
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, _ = w.Write(resp)
}

func ReadRequestData(r *http.Request, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}
	return nil
}

func ReadLimitOffset(r *http.Request) (int32, int32, error) {
	var limit int32 = 5
	var offset int32 = 0
	limitQuery := r.URL.Query().Get("limit")
	if limitQuery != "" {
		limitInt, err := strconv.ParseInt(limitQuery, 10, 32)
		if err != nil {
			return 0, 0, err
		}
		if limitInt < 0 || limitInt > 50 {
			return 0, 0, errors.New("limit must be between 0 and 50")
		}
		limit = int32(limitInt)
	}
	offsetQuery := r.URL.Query().Get("offset")
	if offsetQuery != "" {
		offsetInt, err := strconv.ParseInt(offsetQuery, 10, 32)
		if err != nil {
			return 0, 0, err
		}
		if offsetInt <= 0 {
			return 0, 0, errors.New("offset must be greater than 0")
		}
		offset = int32(offsetInt)
	}
	return limit, offset, nil
}
