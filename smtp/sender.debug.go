package smtp

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

type DebugSender struct {
	writer io.Writer
}

func NewDebugSender(writer io.Writer) *DebugSender {
	if writer == nil {
		writer = os.Stdout
	}

	return &DebugSender{writer: writer}
}

func (sender *DebugSender) SendMail(_ MailUsers, t *template.Template, tName string, data any) error {
	err := t.ExecuteTemplate(sender.writer, tName, data)
	if err != nil {
		return fmt.Errorf("execute template err: %w", err)
	}

	return nil
}

func (sender *DebugSender) Ping() error {
	return nil
}
