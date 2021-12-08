package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/domain/mocks"
)

func TestRegisterUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fakeName, fakePwd := "fakeName", "fakePwd"
		val := url.Values{}
		val.Set("name", fakeName)
		val.Set("password", fakePwd)
		req := httptest.NewRequest("POST", "/register", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		userusecase := new(mocks.UserUsecase)
		postusecase := new(mocks.PostUsecase)

		userusecase.On("StoreUser", mock.Anything, mock.MatchedBy(
			func(user *domain.User) bool {
				return user.Password == fakePwd && user.Name == fakeName
			})).Return(nil)

		mux := http.NewServeMux()
		delivery := NewHttpDelivery(mux, userusecase, postusecase)

		delivery.buildRoute()

		mux.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &responseMap)

		_, ok := responseMap["token"]
		if !ok {
			t.Errorf("Returen not contain token, failed.")
		}
	})
	t.Run("failure", func(t *testing.T) {
		fakeName, fakePwd := "", ""

		val := url.Values{}
		val.Set("name", fakeName)
		val.Set("password", fakePwd)
		req := httptest.NewRequest("POST", "/register", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		userusecase := new(mocks.UserUsecase)
		postusecase := new(mocks.PostUsecase)
		mux := http.NewServeMux()
		delivery := NewHttpDelivery(mux, userusecase, postusecase)

		delivery.buildRoute()

		mux.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &responseMap)

		_, ok := responseMap["err"]
		if !ok {
			t.Error("Return not contain err, failed.")
		}
	})

}
