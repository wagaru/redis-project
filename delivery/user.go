package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func (d *delivery) routeUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.getUsers(w, r)
	default:
		FailureResponse(w, errors.New("invalid request"))
	}
}

func (d *delivery) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := d.usecase.FetchUsers(context.Background())
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(users)
	w.Write([]byte(response))
}
