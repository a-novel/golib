package smtp_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/smtp"
)

func TestTestSender(t *testing.T) {
	t.Parallel()

	sender := smtp.NewTestSender()

	listenCondition1 := func(mail *smtp.TestMail) bool {
		mapData, ok := mail.Data.(map[string]any)
		if !ok {
			return false
		}

		return mapData["test"] == "foo"
	}

	listenCondition2 := func(mail *smtp.TestMail) bool {
		mapData, ok := mail.Data.(map[string]any)
		if !ok {
			return false
		}

		return mapData["test"] == "bar"
	}

	listenCondition3 := func(mail *smtp.TestMail) bool {
		mapData, ok := mail.Data.(map[string]any)
		if !ok {
			return false
		}

		return mapData["test"] == "baz"
	}

	waiter1 := sender.GetWaiter(listenCondition1, time.Second)
	waiter2 := sender.GetWaiter(listenCondition2, time.Second)
	waiter3 := sender.GetWaiter(listenCondition3, time.Second)

	go func() {
		_ = sender.SendMail([]string{"user"}, nil, "", map[string]any{"test": "foo"})
	}()

	go func() {
		_ = sender.SendMail([]string{"user"}, nil, "", map[string]any{"test": "bar"})
	}()

	res1, ok1 := waiter1.Wait()
	res2, ok2 := waiter2.Wait()
	res3, ok3 := waiter3.Wait()

	require.True(t, ok1)
	require.True(t, ok2)
	require.False(t, ok3)

	require.Equal(t, &smtp.TestMail{
		To:   []string{"user"},
		Data: map[string]any{"test": "foo"},
	}, res1)
	require.Equal(t, &smtp.TestMail{
		To:   []string{"user"},
		Data: map[string]any{"test": "bar"},
	}, res2)
	require.Nil(t, res3)
}
