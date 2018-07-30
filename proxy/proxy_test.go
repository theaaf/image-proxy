package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfiguration_AllowsHost(t *testing.T) {
	for name, tc := range map[string]struct {
		Allowed  []string
		Host     string
		Expected bool
	}{
		"Match": {
			Allowed:  []string{"foo.example.com", "bar.example.com"},
			Host:     "foo.example.com",
			Expected: true,
		},
		"Mismatch": {
			Allowed:  []string{"foo.example.com", "bar.example.com"},
			Host:     "baz.example.com",
			Expected: false,
		},
		"WildcardMatch": {
			Allowed:  []string{"*.foo.example.com", "bar.example.com"},
			Host:     "foo.foo.example.com",
			Expected: true,
		},
		"WildcardMismatch": {
			Allowed:  []string{"*.foo.example.com", "bar.example.com"},
			Host:     "foo.example.com",
			Expected: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			config := &Configuration{
				AllowedHosts: tc.Allowed,
			}
			assert.Equal(t, tc.Expected, config.AllowsHost(tc.Host))
		})
	}
}

func TestProxy(t *testing.T) {
	for name, tc := range map[string]struct {
		Configuration      *Configuration
		URL                string
		ExpectedStatusCode int
	}{
		"ForbiddenHost": {
			Configuration: &Configuration{
				AllowedHosts: []string{"foo.example.com", "bar.example.com"},
			},
			URL:                "baz.example.com/foo.png",
			ExpectedStatusCode: http.StatusForbidden,
		},
	} {
		t.Run(name, func(t *testing.T) {
			url, err := url.Parse(tc.URL)
			require.NoError(t, err)

			config := tc.Configuration
			if config == nil {
				config = &Configuration{}
			}
			request, err := NewRequestFromURL(url)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			Proxy(config, rec, request)
			result := rec.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.ExpectedStatusCode, result.StatusCode)
		})
	}
}
