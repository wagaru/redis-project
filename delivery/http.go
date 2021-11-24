package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/usecase"
)

var decoder = schema.NewDecoder()

type delivery struct {
	usecase usecase.Usecase
}

type httpDelivery interface {
	Run(string) error
}

var LoginUser = make(map[string]*domain.User)

func NewHttpDelivery(usecase usecase.Usecase) httpDelivery {
	return &delivery{
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
	mux := http.NewServeMux()
	d.buildRoute(mux)
	return http.ListenAndServe(port, mux)
}

func (d *delivery) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/register", d.routeRegister)
	mux.HandleFunc("/users", d.routeUsers)
	mux.HandleFunc("/posts", d.routePosts)
}
