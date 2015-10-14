package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GithubConnect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	_, err := getValidProjectIDFromToken(r)
	if err != nil {
		http.Error(w, "The token is invalid", http.StatusBadRequest)
		return
	}

	conf := newGithub()
	conf.RedirectURL = config.ApiBaseURL + "/api/v1/github/callback?token=" + token
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	projectID, err := getValidProjectIDFromToken(r)
	if err != nil {
		http.Error(w, "The token is invalid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	conf := newGithub()
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		logger.Error("could not exhange gituh code for token", "err", err, "code", code)
		http.Error(w, "could not exchange github code for token", http.StatusBadRequest)
		return
	}

	repository, err := models.FindRepositoryByProjectID(projectID)
	if err == models.ErrNotFound {
		repository = &models.Repository{
			ProjectID: projectID,
			Provider:  "github",
		}
	} else if err != nil {
		logger.Error("could not find repository", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	repository.GithubScope = "repo"
	repository.GithubToken = tok.AccessToken
	repository.GithubTokenType = tok.TokenType

	err = models.InsertRepository(repository)
	if err != nil {
		logger.Error("could not insert repository", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/projects/%d/config", config.AppBaseURL, projectID)
	http.Redirect(w, r, url, http.StatusFound)
}

func newGithub() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"repo"},
	}
}

func getValidProjectIDFromToken(r *http.Request) (int64, error) {
	authorization := r.URL.Query().Get("token")
	if authorization == "" {
		return 0, errors.New("no project token present")
	}

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.TokenSigningKey, nil
	})

	if err != nil {
		return 0, errors.New("bad project token")
	}

	id := token.Claims["project_id"].(float64)
	return int64(id), nil
}
