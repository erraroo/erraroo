package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/app"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

var (
	_app     *app.App
	_account *models.Account
	_user    *models.User
	_token   string
)

func TestMain(m *testing.M) {
	config.Env = "test"
	config.Postgres = "dbname=erraroo_test sslmode=disable"

	models.Setup()
	models.SetupForTesting()
	api.Limiter = api.NoLimiter()

	_app = app.New()
	_app.SetupForTesting()
	defer _app.Shutdown()

	var err error

	_account, err = models.CreateAccount()
	if err != nil {
		panic(err)
	}

	_user, err = models.Users.Create(uniqEmail(), "password", _account)
	if err != nil {
		panic(err)
	}

	req, res := rr("POST", "/api/v1/sessions", api.SigninRequest{api.Signin{_user.Email, "password"}})
	_app.ServeHTTP(res, req)
	response := map[string]string{}
	json.NewDecoder(res.Body).Decode(&response)
	_token = response["token"]

	ret := m.Run()
	os.Exit(ret)
}

func uniqEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().Nanosecond())
}

func rr(method string, path string, payload interface{}) (*http.Request, *httptest.ResponseRecorder) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	return req, httptest.NewRecorder()
}

func signin(t *testing.T, email string, password string) string {
	req, res := rr("POST", "/api/v1/sessions", api.SigninRequest{api.Signin{email, password}})
	_app.ServeHTTP(res, req)

	if !assert.Equal(t, http.StatusCreated, res.Code) {
		errors, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("Could not login %d: %s\n", res.Code, string(errors))
	}

	if !assert.NotEmpty(t, res.Header().Get("Set-Cookie")) {
		errors, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("Did not set the cookie! %d: %s\n", res.Code, string(errors))
	}

	response := map[string]string{}
	json.NewDecoder(res.Body).Decode(&response)
	return response["token"]
}

func getCookie(res *httptest.ResponseRecorder) http.Cookie {
	setCookie := res.Header().Get("Set-Cookie")
	parts := strings.Split(strings.TrimSpace(setCookie), ";")
	parts = strings.Split(parts[0], "=")
	return http.Cookie{Name: parts[0], Value: parts[1]}
}
