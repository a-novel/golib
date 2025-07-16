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
	mails []*TestMail
	mu    sync.RWMutex
}

func NewTestSender() *TestSender {
	return &TestSender{}
}

func (sender *TestSender) SendMail(to []string, _ *template.Template, _ string, data any) error {
	sender.mu.Lock()
	defer sender.mu.Unlock()

	sender.mails = append(sender.mails, &TestMail{
		To:   to,
		Data: data,
	})

	return nil
}

func (sender *TestSender) FindTestMail(cmp func(*TestMail) bool) (*TestMail, bool) {
	sender.mu.RLock()
	defer sender.mu.RUnlock()

	for _, mail := range sender.mails {
		if cmp(mail) {
			return mail, true
		}
	}

	return nil, false
}
