package models

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nerdyworm/sourcemap"
)

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
			return nil, err
		}
		defer response2.Close()

		b, err := ioutil.ReadAll(response2)
		if err != nil {
			return nil, err
		}
		sm, err := sourcemap.Parse(sourceMapURL, b)

		//sm, err := sourcemap.Read(response2)
		if err != nil {
			log.Printf("error sourcemap.Read(%s): %v\n", sourceMapURL, err)
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

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		h.cache[url] = contents
	}

	reader := bytes.NewReader(h.cache[url])
	return ioutil.NopCloser(reader), nil
}
