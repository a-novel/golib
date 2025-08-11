package smtp

import (
	"errors"
	"sync"
	"text/template"
)

var ErrPingTestSender = errors.New("pinging test sender: make sure this is not a misconfiguration")

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

func (sender *TestSender) SendMail(to MailUsers, _ *template.Template, _ string, data any) error {
	sender.mu.Lock()
	defer sender.mu.Unlock()

	sender.mails = append(sender.mails, &TestMail{
		To:   to.Emails(),
		Data: data,
	})

	return nil
}

func (sender *TestSender) Ping() error {
	return ErrPingTestSender
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
