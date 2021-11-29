package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/wagaru/redis-project/usecase"
)

var decoder = schema.NewDecoder()

type delivery struct {
	mux     *http.ServeMux
	usecase usecase.Usecase
}

type httpDelivery interface {
	Run(string) error
	buildRoute()
}

func NewHttpDelivery(mux *http.ServeMux, usecase usecase.Usecase) httpDelivery {
	return &delivery{
		mux:     mux,
		usecase: usecase,
	}
}

func FailureResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json, _ := json.Marshal(map[string]interface{}{
		"err": err.Error(),
	})
	w.Write(json)
}

func (d *delivery) Run(port string) error {
	d.buildRoute()
	return http.ListenAndServe(port, d.mux)
}

func (d *delivery) buildRoute() {
	d.mux.HandleFunc("/register", d.routeRegister)
	d.mux.HandleFunc("/users", d.routeUsers)
	d.mux.HandleFunc("/posts", d.routePosts)
}
