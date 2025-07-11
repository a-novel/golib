package smtp

import (
	"sync"
	"text/template"
)

type TestMail struct {
	To   []string
	Data any
}

type TestSender struct {
	received []TestMail
	mu       sync.RWMutex
}

func NewTestSender() *TestSender {
	return &TestSender{
		received: make([]TestMail, 0),
	}
}

func (sender *TestSender) SendMail(to []string, _ *template.Template, _ string, data any) error {
	sender.mu.Lock()
	defer sender.mu.Unlock()

	sender.received = append(sender.received, TestMail{
		To:   to,
		Data: data,
	})

	return nil
}

func (sender *TestSender) GetReceived() []TestMail {
	sender.mu.RLock()
	defer sender.mu.RUnlock()

	return append([]TestMail{}, sender.received...)
}
