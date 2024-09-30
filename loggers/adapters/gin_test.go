package adapters_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/adapters"
	"github.com/a-novel/golib/loggers/formatters"
	formattersmocks "github.com/a-novel/golib/loggers/formatters/mocks"
)

func TestGin(t *testing.T) {
	// Prevents useless logs.
	gin.SetMode(gin.ReleaseMode)

	testCases := []struct {
		name string

		status int
		errors []error

		method string
		path   string
		query  string

		userAgent   string
		clientIP    string
		contentType string
		proto       string
		projectID   string
		trace       string

		expectLogLevel loggers.LogLevel
		expectConsole  *regexp.Regexp
		expectJSON     interface{}
	}{
		{
			name:   "SimpleRequest",
			status: http.StatusOK,

			method:      http.MethodGet,
			path:        "/foo",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",

			expectLogLevel: loggers.LogLevelInfo,
			expectConsole:  regexp.MustCompile("^✓ 200 \\[GET /foo] \\([^)]+\\)\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"severity":    "INFO",
				"ip":          "127.0.0.1",
				"contentType": "application/json",
				"errors":      []string(nil),
				"query":       url.Values{},
			},
		},
		{
			name:   "Warning",
			status: http.StatusNotFound,

			method:      http.MethodGet,
			path:        "/foo",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",

			expectLogLevel: loggers.LogLevelWarning,
			expectConsole:  regexp.MustCompile("^⟁ 404 \\[GET /foo] \\([^)]+\\)\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        404,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"severity":    "WARNING",
				"ip":          "127.0.0.1",
				"contentType": "application/json",
				"errors":      []string(nil),
				"query":       url.Values{},
			},
		},
		{
			name:   "Error",
			status: http.StatusInternalServerError,

			method:      http.MethodGet,
			path:        "/foo",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",

			expectLogLevel: loggers.LogLevelError,
			expectConsole:  regexp.MustCompile("^✗ 500 \\[GET /foo] \\([^)]+\\)\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        500,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"severity":    "ERROR",
				"ip":          "127.0.0.1",
				"contentType": "application/json",
				"errors":      []string(nil),
				"query":       url.Values{},
			},
		},

		{
			name:   "WithTrace",
			status: http.StatusOK,

			method:      http.MethodGet,
			path:        "/foo",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",
			trace:       "xyz",

			expectLogLevel: loggers.LogLevelInfo,
			expectConsole:  regexp.MustCompile("^✓ 200 \\[GET /foo] \\([^)]+\\)\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"logging.googleapis.com/trace": "projects/hello-world/traces/xyz",
				"severity":                     "INFO",
				"ip":                           "127.0.0.1",
				"contentType":                  "application/json",
				"errors":                       []string(nil),
				"query":                        url.Values{},
			},
		},

		{
			name:   "WithQuery",
			status: http.StatusOK,

			method:      http.MethodGet,
			path:        "/foo",
			query:       "?a=1&a=3&b=2",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",

			expectLogLevel: loggers.LogLevelInfo,
			expectConsole:  regexp.MustCompile("^✓ 200 \\[GET /foo] \\([^)]+\\)┌─┬─┐\n│a│1│\n│ │3│\n│b│2│\n└─┴─┘\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        200,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"severity":    "INFO",
				"ip":          "127.0.0.1",
				"contentType": "application/json",
				"errors":      []string(nil),
				"query": url.Values{
					"a": []string{"1", "3"},
					"b": []string{"2"},
				},
			},
		},
		{
			name:   "WithErrors",
			status: http.StatusInternalServerError,
			errors: []error{
				errors.New("foo"),
				errors.New("bar"),
			},

			method:      http.MethodGet,
			path:        "/foo",
			userAgent:   "Netscape",
			clientIP:    "127.0.0.1",
			contentType: "application/json",
			proto:       "HTTP/1.1",
			projectID:   "hello-world",

			expectLogLevel: loggers.LogLevelError,
			expectConsole:  regexp.MustCompile("^✗ 500 \\[GET /foo] \\([^)]+\\)\n\n  - foo\n  - bar\n\n$"),
			expectJSON: map[string]interface{}{
				"httpRequest": map[string]interface{}{
					"requestMethod": "GET",
					"requestUrl":    "/foo",
					"status":        500,
					"userAgent":     "Netscape",
					"remoteIp":      "127.0.0.1",
					"protocol":      "HTTP/1.1",
				},
				"severity":    "ERROR",
				"ip":          "127.0.0.1",
				"contentType": "application/json",
				"errors":      []string{"foo", "bar"},
				"query":       url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			c, r := gin.CreateTestContext(rw)
			r.TrustedPlatform = "abcd"

			req := httptest.NewRequest(tc.method, tc.path+tc.query, nil)
			req.Header.Set("User-Agent", tc.userAgent)
			req.Header.Set("abcd", tc.clientIP)
			req.Header.Set("Content-Type", tc.contentType)
			req.Header.Set("X-Cloud-Trace-Context", tc.trace)
			req.Proto = tc.proto

			c.Request = req

			formatter := formattersmocks.NewMockFormatter(t)
			adapter := adapters.NewGin(formatter, tc.projectID)

			var consoleResult string
			var jsonResult interface{}

			formatter.
				On("Log", mock.Anything, tc.expectLogLevel).
				Run(func(args mock.Arguments) {
					content := args.Get(0).(formatters.LogSplit)
					consoleResult = content.RenderConsole()
					jsonResult = content.RenderJSON()
				})

			r.Handle(tc.method, tc.path, adapter.Middleware(), func(c *gin.Context) {
				for _, err := range tc.errors {
					_ = c.Error(err)
				}

				c.String(tc.status, "")
			})

			r.ServeHTTP(rw, req)

			require.Regexp(t, tc.expectConsole, consoleResult)
			require.Empty(t, cmp.Diff(
				tc.expectJSON, jsonResult,
				cmpopts.IgnoreMapEntries(func(k string, _ interface{}) bool {
					return k == "start" || k == "latency"
				}),
			))

			formatter.AssertExpectations(t)
		})
	}
}
