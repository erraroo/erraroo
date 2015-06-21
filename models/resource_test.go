package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceSourceMapURL(t *testing.T) {
	examples := []struct {
		url      string
		source   string
		expected string
	}{
		{
			"http://localhost:3000/assets/vendor.js",
			"//# sourceMappingURL=vendor.map",
			"http://localhost:3000/assets/vendor.map",
		},
		{
			"http://localhost:3000/assets/subfolder/vendor.js",
			"//# sourceMappingURL=vendor.map",
			"http://localhost:3000/assets/subfolder/vendor.map",
		},
		{
			"https://xxx.cloudfront.com/assets/vendor.js",
			"//# sourceMappingURL=assets/vendor.map",
			"https://xxx.cloudfront.com/assets/vendor.map",
		},
		{
			"https://xxx.cloudfront.com/assets/vendor.js",
			"//# sourceMappingURL=/assets/vendor.map",
			"https://xxx.cloudfront.com/assets/vendor.map",
		},
	}

	for _, example := range examples {
		resource := Resource{
			URL:    example.url,
			Source: example.source,
		}

		assert.Equal(t, example.expected, resource.SourceMapURL())
	}
}
