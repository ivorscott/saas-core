package sendgrid

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender func(email *mail.SGMailV3) (*rest.Response, error)

func SendMail(email *mail.SGMailV3, send Sender) (*rest.Response, error) {
	var resp *rest.Response

	resp, err := send(email)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
