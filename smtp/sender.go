package smtp

import (
	"text/template"
)

type Sender interface {
	SendMail(to []string, t *template.Template, tName string, data any) error
	Ping() error
}
