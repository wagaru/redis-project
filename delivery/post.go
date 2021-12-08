package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wagaru/redis-project/domain"
)

func (d *delivery) routePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.getPosts(w, r)
	case http.MethodPost:
		d.storePosts(w, r)
	default:
		FailureResponse(w, errors.New("invalid request"))
	}
}

func (d *delivery) getPosts(w http.ResponseWriter, r *http.Request) {
	var queryParams domain.PostQueryParams
	err := decoder.Decode(&queryParams, r.Form)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	posts, err := d.postusecase.FetchPosts(context.Background(), &queryParams)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	json, err := json.Marshal(posts)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (d *delivery) storePosts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	r.ParseForm()
	user, err := d.userusecase.FetchUserByToken(ctx, r.Form.Get("token"))
	if err != nil {
		FailureResponse(w, err)
		return
	}
	err = d.postusecase.StorePost(ctx, &domain.Post{
		Title: r.Form.Get("title"),
	}, user)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Write([]byte("Success"))
}
