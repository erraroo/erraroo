package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockResourceGetter struct{}

func (m *mockResourceGetter) Get(url string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("../test/fixtures/%x", md5.Sum([]byte(url)))
	return os.Open(filename)
}

func TestResourceStore(t *testing.T) {
	sourceURL := "http://erraroo.com/assets/erraroo-50c293f8d38851c3fa2c04bb62fb8e3f.js"
	resources := &resourcesStore{&mockResourceGetter{}}

	resource, err := resources.FindByURL(sourceURL)
	assert.Nil(t, err)
	assert.NotEmpty(t, resource.SourceMap)

	//context := resource.Context(536, 15)
	//assert.Equal(t, "        throw new Error('i threw an error');", context.ContextLine)
	//assert.Equal(t, 17, context.OrigLine)
	//assert.Equal(t, 1, context.OrigColumn)
	//assert.Equal(t, "erraroo/routes/application.js", context.OrigFile)
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
