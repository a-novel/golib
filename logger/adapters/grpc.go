package adapters

import (
	"context"
	"fmt"
	"github.com/a-novel/golib/logger"
	"github.com/a-novel/golib/logger/formatters"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type GRPC interface {
	// Report logs the result of an RPC call to a specified service.
	//
	// Since gRPC is not aware of middlewares, you can use the WrapGRPCCall function over your services, to
	// automatically log the result of each call.
	Report(service string, err error)
}

type grpcImpl struct {
	formatter formatters.Formatter
}

type grpcReport struct {
	code  codes.Code
	level logger.LogLevel
	color lipgloss.Color

	consolePrefix string
}

type grpcMetrics struct {
	latency time.Duration
}

// Capture reporting from a GRPC call and format them for the logger.
func (g *grpcImpl) getReport(err error) *grpcReport {
	output := grpcReport{
		code:          codes.OK,
		level:         logger.LogLevelInfo,
		color:         lipgloss.Color("#00A7FF"),
		consolePrefix: "✓ ",
	}

	if err != nil {
		output.code = status.Code(err)
		if output.code == 0 {
			output.code = codes.Unknown
		}

		if output.code == codes.Unavailable || output.code == codes.Canceled || output.code == codes.Unimplemented {
			// Reserve special (yellow) treatment for codes that likely indicate a service / implementation issues.
			// GRPC does not have clear distinction between server-side and client-side errors (such as HTTP).
			output.level = logger.LogLevelWarning
			output.color = "#FF8000"
			output.consolePrefix = "⟁ "
		} else {
			// Regular codes are treated as standard errors.
			output.level = logger.LogLevelError
			output.color = "#FF3232"
			output.consolePrefix = "✗ "
		}
	}

	return &output
}

// Returns a console message tailored for the console logger.
func (g *grpcImpl) getConsoleMessage(service string, err error, report *grpcReport, metrics *grpcMetrics) string {
	errorMessage := "\n\n"
	if err != nil {
		errorMessage = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#FF3232")).
			Render("\n\n" + err.Error() + "\n\n")
	}

	latencyMessage := ""
	if metrics != nil {
		latencyMessage = lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" (%s)", metrics.latency))
	}

	return lipgloss.NewStyle().Foreground(report.color).Bold(true).Render(report.consolePrefix+report.code.String()) +
		lipgloss.NewStyle().Foreground(report.color).Render(fmt.Sprintf(" [%s]", service)) +
		latencyMessage +
		errorMessage
}

// Returns a JSON message tailored for Google Cloud Logging.
func (g *grpcImpl) getJSONMessageGCP(service string, err error, report *grpcReport, metrics *grpcMetrics) map[string]interface{} {
	severity := lo.Switch[logger.LogLevel, string](report.level).
		Case(logger.LogLevelInfo, "INFO").
		Case(logger.LogLevelWarning, "WARNING").
		Case(logger.LogLevelError, "ERROR").
		Case(logger.LogLevelFatal, "ERROR").
		Default("INFO")

	// TODO: check if we can add trace to GRPC requests.
	// TODO: improve formatting of GRPC messages.
	grpcRequest := map[string]interface{}{
		"service": service,
		"code":    report.code,
	}

	output := map[string]interface{}{
		"severity":    severity,
		"grpcRequest": grpcRequest,
	}

	if metrics != nil {
		grpcRequest["latency"] = metrics.latency
	}

	if err != nil {
		output["error"] = err.Error()
	}

	return output
}

func (g *grpcImpl) reportWith(service string, err error, metrics *grpcMetrics) {
	meta := g.getReport(err)

	message := formatters.NewSplit().
		SetConsoleMessage(g.getConsoleMessage(service, err, meta, metrics)).
		SetJSONMessage(g.getJSONMessageGCP(service, err, meta, metrics))

	g.formatter.Log(message, meta.level)
}

func (g *grpcImpl) Report(service string, err error) {
	g.reportWith(service, err, nil)
}

func NewGRPC(formatter formatters.Formatter) GRPC {
	return &grpcImpl{formatter}
}

type GRPCCallback[In any, Out any] func(ctx context.Context, in In) (Out, error)

// WrapGRPCCall wraps a GRPCCallback callback with a GRPC logger. It also adds extra reporting that are not available
// to the base class, such as latency info.
func WrapGRPCCall[In any, Out any](service string, logger GRPC, GRPCCallback GRPCCallback[In, Out]) GRPCCallback[In, Out] {
	return func(ctx context.Context, in In) (Out, error) {
		start := time.Now()
		out, err := GRPCCallback(ctx, in)
		end := time.Now()

		gcpLogger, ok := logger.(*grpcImpl)
		if !ok {
			logger.Report(service, err)
		} else {
			gcpLogger.reportWith(service, err, &grpcMetrics{latency: end.Sub(start)})
		}

		return out, err
	}
}
