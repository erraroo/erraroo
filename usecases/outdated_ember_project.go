package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	"golang.org/x/oauth2"

	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/google/go-github/github"
)

var ErrNoRepo = errors.New("no repository")

type DependencyChecker interface {
	Outdated(*models.Repository) (*models.Revision, error)
}

func CheckEmberDependencies(projectID int64, checker DependencyChecker) error {
	repository, err := models.FindRepositoryByProjectID(projectID)
	if err == models.ErrNotFound {
		return ErrNoRepo
	}

	if err != nil {
		return err
	}

	if checker == nil {
		checker = &githubNodeDepencyChecker{}
	}

	revision, err := checker.Outdated(repository)
	if err != nil {
		return err
	}

	err = models.SaveRevision(revision)
	if err != nil {
		return err
	}

	return nil
}

type githubNodeDepencyChecker struct{}

func (g githubNodeDepencyChecker) Outdated(r *models.Repository) (*models.Revision, error) {
	root := "/tmp/" + uuid() + "/"
	os.MkdirAll(root, 0766)
	defer os.RemoveAll(root)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.GithubToken},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	bowerjson, _, _, err := client.Repositories.GetContents(r.GithubOrg, r.GithubRepo, "bower.json", nil)
	if err != nil {
		return nil, err
	}

	bowerjsonbytes, err := bowerjson.Decode()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(root+"bower.json", bowerjsonbytes, 0644)
	if err != nil {
		logger.Error("could not write bower.json", "err", err)
		return nil, err
	}

	packagejson, _, _, err := client.Repositories.GetContents(r.GithubOrg, r.GithubRepo, "package.json", nil)
	if err != nil {
		logger.Error("could not get package.json", "err", err)
		return nil, err
	}

	packagejsonbytes, err := packagejson.Decode()
	if err != nil {
		logger.Error("could not decode package.json", "err", err)
		return nil, err
	}

	err = ioutil.WriteFile(root+"package.json", packagejsonbytes, 0644)
	if err != nil {
		logger.Error("could not write package.json", "err", err)
		return nil, err
	}

	branch, _, err := client.Repositories.GetBranch("erraroo", "erraroo-app", "master")
	if err != nil {
		logger.Error("could not get branch", "err", err, "repository", r)
		return nil, err
	}

	outdated, err := errarooNodeOutdated(root)
	if err != nil {
		return nil, err
	}

	outdated.SHA = *branch.Commit.SHA
	outdated.ProjectID = r.ProjectID
	return outdated, nil
}

func errarooNodeOutdated(path string) (*models.Revision, error) {
	cmd := "erraroo-node-outdated"

	c := exec.Command(cmd)
	c.Dir = path

	var out bytes.Buffer
	var errs bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errs
	err := c.Run()
	if err != nil {
		logger.Error("erraroo-node-outdated", "stdout", string(out.Bytes()), "stderr", string(errs.Bytes()), "path", path, "err", err)
		return nil, err
	}

	outdated := &models.Revision{}
	err = json.Unmarshal(out.Bytes(), outdated)
	if err != nil {
		logger.Error("could not unmarhsal erraroo-node-oudated", "stdout", string(out.Bytes()), "stderr", string(errs.Bytes()), "path", path, "err", err)
		return nil, err
	}

	return outdated, nil
}
