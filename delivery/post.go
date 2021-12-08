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
		d.handlePosts(w, r)
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

func (d *delivery) handlePosts(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Form.Get("method")
	if method == "store" {
		d.storePosts(w, r)
	} else if method == "vote" {
		d.votePost(w, r)
	} else {
		FailureResponse(w, errors.New("unsupported method"))
		return
	}
}

func (d *delivery) storePosts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// r.ParseForm()

	token := r.Form.Get("token")
	if token == "" {
		FailureResponse(w, errors.New("invalid token."))
		return
	}
	title := r.Form.Get("title")
	if title == "" {
		FailureResponse(w, errors.New("invalid title."))
		return
	}
	user, err := d.userusecase.FetchUserByToken(ctx, token)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	err = d.postusecase.StorePost(ctx, &domain.Post{
		Title: title,
	}, user)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Write([]byte("Success"))
}

func (d *delivery) votePost(w http.ResponseWriter, r *http.Request) {
	postId := r.Form.Get("postID")
	if postId == "" {
		FailureResponse(w, errors.New("invalid postID"))
		return
	}
	token := r.Form.Get("token")
	if token == "" {
		FailureResponse(w, errors.New("invalid token"))
		return
	}
	ctx := context.Background()
	user, err := d.userusecase.FetchUserByToken(ctx, token)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	post, err := d.postusecase.FetchPostByID(ctx, postId)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	err = d.postusecase.Vote(ctx, post, user)
	if err != nil {
		FailureResponse(w, err)
		return
	}
	w.Write([]byte("Success"))
}
