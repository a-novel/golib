package smtp

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

type MailUser struct {
	Name  string
	Email string
}

func (mailUser MailUser) String() string {
	return fmt.Sprintf("%s <%s>", mailUser.Name, mailUser.Email)
}

type MailUsers []MailUser

func (mailUsers MailUsers) String() string {
	return strings.Join(lo.Map(mailUsers, func(item MailUser, _ int) string {
		return item.String()
	}), ", ")
}

func (mailUsers MailUsers) Emails() []string {
	return lo.Map(mailUsers, func(item MailUser, _ int) string {
		return item.Email
	})
}

type Sender interface {
	SendMail(to MailUsers, t *template.Template, tName string, data any) error
	Ping() error
}
