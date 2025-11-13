package logging

import (
	"net/http"
)

type HttpConfig interface {
	Logger() func(http.Handler) http.Handler
}
