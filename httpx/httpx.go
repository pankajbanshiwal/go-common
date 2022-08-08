package httpx

import (
	"net/http"
	"github.com/okcredit/go-common/errors"
	"github.com/okcredit/go-common/encoding/json"
)

func WriteJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error) error {
	e := errors.FromError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	return json.NewEncoder(w).Encode(e)
}
