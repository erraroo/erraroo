package models

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/nerdyworm/sourcemap"
)

var ErrCouldNotGet = errors.New("could not get resource")

type resourcesStore struct {
	Getter resourceGetter
}

func NewResourceStore() *resourcesStore {
	return &resourcesStore{newHttpResourceGetter()}
}

type resourceGetter interface {
	Get(string) (io.ReadCloser, error)
}

type httpResourseGetter struct {
	client *http.Client
	cache  map[string][]byte
}

func (s *resourcesStore) FindByURL(url string) (*Resource, error) {
	response, err := s.Getter.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	source, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}

	resource := &Resource{URL: url, Source: string(source)}

	sourceMapURL := resource.SourceMapURL()
	if sourceMapURL != "" {
		response2, err := s.Getter.Get(sourceMapURL)
		if err != nil {
			logger.Error("could not fetch source map", "url", sourceMapURL, "err", err)
			return nil, err
		}
		defer response2.Close()

		b, err := ioutil.ReadAll(response2)
		if err != nil {
			return nil, err
		}

		sm, err := sourcemap.Parse(sourceMapURL, b)
		if err != nil {
			logger.Error("sourcemap.Parse", "sourceMapURL", sourceMapURL, "err", err)
			return nil, err
		}

		resource.SourceMap = sm
	}

	return resource, nil
}

func newHttpResourceGetter() *httpResourseGetter {
	timeout := time.Duration(10 * time.Second)

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}

	client := &http.Client{
		Transport: &transport,
	}

	return &httpResourseGetter{client, make(map[string][]byte)}
}

func (h *httpResourseGetter) Get(url string) (io.ReadCloser, error) {
	if _, ok := h.cache[url]; !ok {
		response, err := h.client.Get(url)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		logger.Info("httpResourseGetter.Get", "url", url, "status", response.StatusCode)

		if response.StatusCode == 403 {
			return nil, ErrCouldNotGet
		}

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		h.cache[url] = contents
	}

	reader := bytes.NewReader(h.cache[url])
	return ioutil.NopCloser(reader), nil
}
