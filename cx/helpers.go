package cx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ErrLoginRequired = errors.New("login required")

func GetID(r *http.Request) (int64, error) {
	id, err := StrToID(mux.Vars(r)["id"])
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func StrToID(id string) (int64, error) {
	i, err := strconv.Atoi(id)
	return int64(i), err
}

// JSON renders json
func JSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

// Decode a json payload in the request
func Decode(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
