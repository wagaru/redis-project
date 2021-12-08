package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wagaru/redis-project/domain"
)

func (d *delivery) routeRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		d.registerUser(w, r)
	default:
		FailureResponse(w, errors.New("invalid request"))
	}
}

func (d *delivery) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		FailureResponse(w, err)
		return
	}
	if r.FormValue("name") == "" || r.FormValue("password") == "" {
		FailureResponse(w, errors.New("invalid arguments"))
		return
	}
	user := &domain.User{
		Name:     r.Form.Get("name"),
		Password: r.Form.Get("password"),
	}
	err = d.userusecase.StoreUser(context.Background(), user)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(map[string]interface{}{
		"token": user.Token,
	})
	w.Write([]byte(response))
}
