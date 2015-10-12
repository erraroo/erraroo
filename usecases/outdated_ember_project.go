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
	Outdated(*models.Repository) (*models.OutdatedRevision, error)
}

type githubNodeDepencyChecker struct{}

func (g githubNodeDepencyChecker) Outdated(r *models.Repository) (*models.OutdatedRevision, error) {
	// possibly need to pass the project in here to figure out
	// which DependencyChecker to use, so far it's just node and github
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

	packagejson, _, _, err := client.Repositories.GetContents("erraroo", "erraroo-app", "package.json", nil)
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

func OutdatedEmberProject(project *models.Project, checker DependencyChecker) error {
	repository, err := models.FindRepositoryByProjectID(project.ID)
	if err == models.ErrNotFound {
		return ErrNoRepo
	}

	if err != nil {
		return err
	}

	outdated, err := checker.Outdated(repository)
	if err != nil {
		return err
	}

	if !outdated.Empty() {
		err = models.InsertOutdatedRevision(outdated)
		if err != nil {
			return err
		}
	}

	return nil
}

func errarooNodeOutdated(path string) (*models.OutdatedRevision, error) {
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

	outdated := &models.OutdatedRevision{}
	err = json.Unmarshal(out.Bytes(), outdated)
	if err != nil {
		logger.Error("could not unmarhsal erraroo-node-oudated", "stdout", string(out.Bytes()), "stderr", string(errs.Bytes()), "path", path, "err", err)
		return nil, err
	}

	return outdated, nil
}
