package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"super-descuentos/errs"
)

func ValidateUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidID
	}
	return id, nil
}

func GetOffsetAndLimit(r *http.Request) (int, int, error) {
	offset := 0
	limit := 10

	if r.URL.Query().Get("offset") != "" {
		offsetInt, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			return 0, 0, errs.ErrInvalidQueryParams
		}
		offset = offsetInt
	}

	if r.URL.Query().Get("limit") != "" {
		limitInt, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			return 0, 0, errs.ErrInvalidQueryParams
		}
		limit = limitInt
	}

	return offset, limit, nil
}

func JsonWithErrors(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func SendErrorMessage(w http.ResponseWriter, err error, code int) {
	JsonWithErrors(w, map[string]string{"message": err.Error()}, code)
}
