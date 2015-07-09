package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventJSExtractsTags(t *testing.T) {
	event := &Event{
		Kind:    "js.error",
		Payload: `{"useragent":"iterm","url":"http://example.com"}`,
	}

	tags := event.Tags()
	assert.NotEmpty(t, tags)
	assert.Equal(t, "js.useragent", tags[0].Key)
	assert.Equal(t, "iterm", tags[0].Value)
	assert.Equal(t, "js.url", tags[1].Key)
	assert.Equal(t, "http://example.com", tags[1].Value)
}
