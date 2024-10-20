package adapters

import (
	"fmt"
	"github.com/a-novel/golib/logger"
	"github.com/a-novel/golib/logger/formatters"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"time"
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
	level       logger.LogLevel
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

func (g *ginImpl) getReport(c *gin.Context) *ginReport {
	start := time.Now()
	c.Next()
	end := time.Now()

	output := ginReport{
		start:         start,
		end:           end,
		latency:       end.Sub(start),
		level:         logger.LogLevelInfo,
		color:         lipgloss.Color("#00A7FF"),
		consolePrefix: "✓ ",
		code:          c.Writer.Status(),
		errors:        c.Errors.Errors(),
		query:         c.Request.URL.Query(),
		path:          c.FullPath(),
		verb:          c.Request.Method,
		userAgent:     c.Request.UserAgent(),
		remoteIP:      c.ClientIP(),
		protocol:      c.Request.Proto,
		contentType:   c.ContentType(),
	}

	// Allow logs to be grouped in log explorer.
	// https://cloud.google.com/run/docs/logging#run_manual_logging-go
	if g.projectID != "" {
		traceHeader := c.GetHeader("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			output.trace = fmt.Sprintf("projects/%s/traces/%s", g.projectID, traceParts[0])
		}
	}

	if output.code > 499 {
		output.level = logger.LogLevelError
		output.color = "#FF0000"
		output.consolePrefix = "✗ "
	} else if output.code > 399 || len(output.errors) > 0 {
		output.level = logger.LogLevelWarning
		output.color = "#FF8000"
		output.consolePrefix = "⟁ "
	}

	return &output
}

func (g *ginImpl) getConsoleMessage(report *ginReport) string {
	queryMessage := ""
	if len(report.query) > 0 {
		var rows [][]string
		for k, v := range report.query {
			rows = append(rows, []string{k, strings.Join(v, "\n")})
		}

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

	errorMessage := "\n\n"
	if len(report.errors) > 0 {
		errorMessage = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#FF3232")).
			Render("\n\n" + strings.Join(report.errors, "\n") + "\n\n")
	}

	latencyMessage := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" (%s)", report.latency))

	return lipgloss.NewStyle().Foreground(report.color).Bold(true).Render(fmt.Sprintf("%s%v", report.consolePrefix, report.code)) +
		lipgloss.NewStyle().Foreground(report.color).Render(fmt.Sprintf(" [%s %s]", report.verb, report.path)) +
		latencyMessage +
		queryMessage +
		errorMessage
}

func (g *ginImpl) getJSONMessageGCP(report *ginReport) map[string]interface{} {
	severity := lo.Switch[logger.LogLevel, string](report.level).
		Case(logger.LogLevelInfo, "INFO").
		Case(logger.LogLevelWarning, "WARNING").
		Case(logger.LogLevelError, "ERROR").
		Case(logger.LogLevelFatal, "ERROR").
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

func (g *ginImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		report := g.getReport(c)
		message := formatters.NewSplit().
			SetConsoleMessage(g.getConsoleMessage(report)).
			SetJSONMessage(g.getJSONMessageGCP(report))

		g.formatter.Log(message, report.level)
	}
}

func NewGin(formatter formatters.Formatter, projectID string) Gin {
	return &ginImpl{
		formatter: formatter,
		projectID: projectID,
	}
}
