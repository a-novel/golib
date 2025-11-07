package httpf

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SendJSON[Data any](_ context.Context, w http.ResponseWriter, span trace.Span, data Data) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	span.RecordError(err)
	span.SetStatus(lo.Ternary(err == nil, codes.Ok, codes.Error), "")
}
