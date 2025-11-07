package httpf

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/a-novel/golib/otel"
)

func HandleError(_ context.Context, w http.ResponseWriter, span trace.Span, errMap map[error]int, err error) {
	err = otel.ReportError(span, err)
	status := http.StatusInternalServerError

	for ref, refStatus := range errMap {
		if errors.Is(err, ref) {
			status = refStatus

			break
		}
	}

	http.Error(w, http.StatusText(status), status)
}
