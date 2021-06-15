package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/auth"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"gohanectl/hanectl/test/rest_test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func headerAuthorization(token string) func(next http.Handler) http.Handler {
	handlerFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token))
			next.ServeHTTP(w, r)
		})
	}
	return handlerFunc
}

func TestLogin(t *testing.T) {
	password := "secret"
	jwtSecret := "2222222222222"
	jwtAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	pwdHash, _ := auth.HashPassword(password)
	user := &model.User{
		Username:    "username",
		Password:    password,
		Authorities: []string{"USER"},
		Enabled:     true,
	}
	authRequest := AuthRequest{
		Username: "username",
		Password: pwdHash,
	}
	payloadContent, _ := json.Marshal(authRequest)
	payload := bytes.NewReader(payloadContent)

	userServiceMock := new(mock_test.UserServiceMock)
	userServiceMock.On("FindByUsername", authRequest.Username).Return(user, nil)

	r := chi.NewRouter()
	authApi(r, jwtAuth, userServiceMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "POST", "/auth/login", payload); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	assert.Empty(t, user.Password)
	userServiceMock.AssertExpectations(t)
}

func TestAuthenticated(t *testing.T) {
	jwtSecret := "2222222222222"
	jwtAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	userServiceMock := new(mock_test.UserServiceMock)
	_, token, _ := jwtAuth.Encode(model.Dictionary{UsernameClaimKey: "username", AuthoritiesClaimKey: []string{"USER"}})

	r := chi.NewRouter()
	r.Use(headerAuthorization(token))
	authApi(r, jwtAuth, userServiceMock)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if resp, body := rest_test.DoTestRequest(t, ts, "GET", "/auth/authenticated", nil); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}

	userServiceMock.AssertExpectations(t)
}
