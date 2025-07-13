package smtp

import (
	"text/template"
	"time"

	"github.com/a-novel/golib/chans"
)

type TestMail struct {
	To   []string
	Data any
}

type TestSender struct {
	shared *chans.Shared[*TestMail]
}

func NewTestSender() *TestSender {
	return &TestSender{
		shared: chans.NewShared[*TestMail](),
	}
}

func (sender *TestSender) SendMail(to []string, _ *template.Template, _ string, data any) error {
	sender.shared.Send(&TestMail{
		To:   to,
		Data: data,
	})

	return nil
}

func (sender *TestSender) GetWaiter(
	condition func(mail *TestMail) bool, timeout time.Duration,
) *chans.Waiter[*TestMail] {
	listener := sender.shared.Register()

	waiter := chans.NewWaiter[*TestMail](listener, condition, timeout, func() {
		sender.shared.Unregister(listener)
	})

	return waiter
}
