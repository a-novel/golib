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

	require.NoError(t, sender.SendMail([]string{"user"}, nil, "", map[string]string{"test": "foo"}))
	require.NoError(t, sender.SendMail([]string{"user"}, nil, "", map[string]string{"test": "bar"}))

	res1, ok1 := sender.FindTestMail(func(mail *smtp.TestMail) bool {
		return mail.Data.(map[string]string)["test"] == "foo"
	}, time.Second)
	res2, ok2 := sender.FindTestMail(func(mail *smtp.TestMail) bool {
		return mail.Data.(map[string]string)["test"] == "bar"
	}, time.Second)
	res3, ok3 := sender.FindTestMail(func(mail *smtp.TestMail) bool {
		return mail.Data.(map[string]string)["test"] == "baz"
	}, time.Second)

	require.True(t, ok1)
	require.True(t, ok2)
	require.False(t, ok3)

	require.Equal(t, &smtp.TestMail{
		To:   []string{"user"},
		Data: map[string]string{"test": "foo"},
	}, res1)
	require.Equal(t, &smtp.TestMail{
		To:   []string{"user"},
		Data: map[string]string{"test": "bar"},
	}, res2)
	require.Nil(t, res3)
}
