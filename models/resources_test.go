package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockResourceGetter struct{}

func (m *mockResourceGetter) Get(url string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("../test/fixtures/%x", md5.Sum([]byte(url)))
	_, err := os.Open(filename)
	if err != nil {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		io.Copy(f, resp.Body)
	}

	return os.Open(filename)
}

func TestResourceStore(t *testing.T) {
	sourceURL := "https://d16vxe267myqks.cloudfront.net/assets/erraroo-a97979fee2e8d42dc9a281a233f895b2.js"
	resources := &resourcesStore{&mockResourceGetter{}}

	resource, err := resources.FindByURL(sourceURL)
	assert.Nil(t, err)
	assert.NotEmpty(t, resource.SourceMap)

	context := resource.Context(1, 13983)
	assert.Equal(t, "        throw new Error('i threw an error' + Math.random());", context.ContextLine)
	assert.Equal(t, 17, context.OrigLine)
	assert.Equal(t, 0, context.OrigColumn)
	assert.Equal(t, "erraroo/routes/application.js", context.OrigFile)
}

func TestResourcesWithoutSourceMap(t *testing.T) {
	sourceURL := "http://example.com/assets/lame-application.js"
	resources := &resourcesStore{&mockResourceGetter{}}

	resource, err := resources.FindByURL(sourceURL)
	assert.Nil(t, err)
	assert.Nil(t, resource.SourceMap)

	context := resource.Context(4, 1)
	assert.Equal(t, "  cnsole.log('buggy');", context.ContextLine)
	assert.Equal(t, 4, context.OrigLine)
	assert.Equal(t, 1, context.OrigColumn)
	assert.Equal(t, "/assets/lame-application.js", context.OrigFile)
}
