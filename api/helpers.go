package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/erraroo/erraroo/models"
	"github.com/gorilla/mux"
)

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

func QueryToID(r *http.Request, name string) (int64, error) {
	return StrToID(r.URL.Query().Get(name))
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

func Limit(r *http.Request) int {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 50
	}
	return limit
}

func Page(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	return page
}

func QueryOptions(r *http.Request) models.QueryOptions {
	return models.QueryOptions{
		Page:    Page(r),
		PerPage: Limit(r),
	}
}
