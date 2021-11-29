package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	val := url.Values{}
	val.Set("name", "testname")
	val.Set("password", "1234")
	req := httptest.NewRequest("POST", "/register", strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	usecase := NewMockUsecase()
	mux := http.NewServeMux()
	delivery := NewHttpDelivery(mux, usecase)

	delivery.buildRoute()

	mux.ServeHTTP(w, req)

	responseMap := make(map[string]interface{}, 0)
	json.Unmarshal(w.Body.Bytes(), &responseMap)

	if _, ok := responseMap["token"]; !ok {
		t.Errorf("Returen not contain token, failed.")
	}
}
