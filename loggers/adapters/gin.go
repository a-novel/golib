package adapters

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
)

type Gin interface {
	Middleware() gin.HandlerFunc
}

type ginImpl struct {
	formatter formatters.Formatter
	projectID string
}

type ginReport struct {
	start       time.Time
	end         time.Time
	latency     time.Duration
	level       loggers.LogLevel
	color       lipgloss.Color
	code        int
	errors      []string
	query       url.Values
	trace       string
	path        string
	verb        string
	userAgent   string
	remoteIP    string
	protocol    string
	contentType string

	consolePrefix string
}

func (ginAdapter *ginImpl) getReport(ctx *gin.Context) *ginReport {
	start := time.Now()
	ctx.Next()
	end := time.Now()

	output := ginReport{
		start:         start,
		end:           end,
		latency:       end.Sub(start),
		level:         loggers.LogLevelInfo,
		color:         lipgloss.Color("#00A7FF"),
		consolePrefix: "✓ ",
		code:          ctx.Writer.Status(),
		errors:        ctx.Errors.Errors(),
		query:         ctx.Request.URL.Query(),
		path:          ctx.FullPath(),
		verb:          ctx.Request.Method,
		userAgent:     ctx.Request.UserAgent(),
		remoteIP:      ctx.ClientIP(),
		protocol:      ctx.Request.Proto,
		contentType:   ctx.ContentType(),
	}

	// Allow logs to be grouped in log explorer.
	// https://cloud.google.com/run/docs/logging#run_manual_logging-go
	if ginAdapter.projectID != "" {
		traceHeader := ctx.GetHeader("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			output.trace = fmt.Sprintf("projects/%s/traces/%s", ginAdapter.projectID, traceParts[0])
		}
	}

	if output.code > 499 {
		output.level = loggers.LogLevelError
		output.color = "#FF3232"
		output.consolePrefix = "✗ "
	} else if output.code > 399 || len(output.errors) > 0 {
		output.level = loggers.LogLevelWarning
		output.color = "#FF8000"
		output.consolePrefix = "⟁ "
	}

	return &output
}

func (ginAdapter *ginImpl) getConsoleMessage(report *ginReport) string {
	queryMessage := ""
	if len(report.query) > 0 {
		var rows [][]string
		for k, v := range report.query {
			rows = append(rows, []string{k, strings.Join(v, "\n")})
		}

		slices.SortFunc(rows, func(a, b []string) int {
			// SOrt query parameters by alphabetical order.
			return strings.Compare(a[0], b[0])
		})

		queryMessage = table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))).
			StyleFunc(func(_, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
				}

				return lipgloss.NewStyle()
			}).
			Rows(rows...).
			String()
	}

	errorMessage := ""
	if len(report.errors) > 0 {
		listDisplay := list.New().
			ItemStyle(lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color("#FF3232"))).
			Enumerator(list.Dash).
			EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3232")))

		// Why we don't use initialization.
		// https://www.reddit.com/r/golang/comments/vco4rg/cannot_use_type_string_as_the_type_any/
		for _, err := range report.errors {
			listDisplay.Item(err)
		}

		errorMessage = "\n\n" + lipgloss.NewStyle().MarginLeft(2).Render(listDisplay.String())
	}

	latencyMessage := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" (%s)", report.latency))

	return lipgloss.NewStyle().
		Foreground(report.color).
		Bold(true).
		Render(fmt.Sprintf("%s%v", report.consolePrefix, report.code)) +
		lipgloss.NewStyle().Foreground(report.color).Render(fmt.Sprintf(" [%s %s]", report.verb, report.path)) +
		latencyMessage +
		queryMessage +
		errorMessage +
		"\n\n"
}

func (ginAdapter *ginImpl) getJSONMessageGCP(report *ginReport) map[string]interface{} {
	severity := lo.Switch[loggers.LogLevel, string](report.level).
		Case(loggers.LogLevelInfo, "INFO").
		Case(loggers.LogLevelWarning, "WARNING").
		Case(loggers.LogLevelError, "ERROR").
		Case(loggers.LogLevelFatal, "ERROR").
		Default("INFO")

	httpRequest := map[string]interface{}{
		"requestMethod": report.verb,
		"requestUrl":    report.path,
		"status":        report.code,
		"userAgent":     report.userAgent,
		"remoteIp":      report.remoteIP,
		"protocol":      report.protocol,
		"latency":       report.latency.String(),
	}

	output := map[string]interface{}{
		"httpRequest": httpRequest,
		"severity":    severity,
		"start":       report.start,
		"ip":          report.remoteIP,
		"contentType": report.contentType,
		"errors":      report.errors,
		"query":       report.query,
	}

	if len(report.trace) > 0 {
		output["logging.googleapis.com/trace"] = report.trace
	}

	return output
}

func (ginAdapter *ginImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		report := ginAdapter.getReport(c)

		message := formatters.NewSplit().
			SetConsoleRenderer(func() string {
				return ginAdapter.getConsoleMessage(report)
			}).
			SetJSONRenderer(func() interface{} {
				return ginAdapter.getJSONMessageGCP(report)
			})

		ginAdapter.formatter.Log(message, report.level)
	}
}

func NewGin(formatter formatters.Formatter, projectID string) Gin {
	return &ginImpl{
		formatter: formatter,
		projectID: projectID,
	}
}
