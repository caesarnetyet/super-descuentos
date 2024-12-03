package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"super-descuentos/errs"
	"super-descuentos/model"

	"github.com/google/uuid"
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

func HandleErrorResponse(w http.ResponseWriter, err error) {
	var validationError model.ValidationError
	var validationErrors model.ValidationErrors

	switch {
	case errors.As(err, &validationError):
		// Single validation error
		SendErrorMessage(w, err, http.StatusBadRequest)

	case errors.As(err, &validationErrors):
		// Multiple validation errors
		JsonWithErrors(w, map[string]interface{}{"errors": err}, http.StatusBadRequest)

	default:
		// Generic server error
		http.Error(w, "Hubo un problema en tu petición, inténtalo de nuevo más tarde.", http.StatusInternalServerError)
	}
}
