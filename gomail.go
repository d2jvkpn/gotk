package gotk

import (
	"net/mail"

	"gopkg.in/gomail.v2"
)

func NewSender(fp, field string) (sender *Sender, err error) {
	sender = new(Sender)
	objects := map[string]any{"sender": sender}

	if err = UnmarshalYamlObjects(fp, objects); err != nil {
		return nil, err
	}

	return sender, nil
}

type Sender struct {
	SMTPAddr string `mapstructure:"smtp_addr"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

type Email struct {
	Recipients []string `mapstructure:"recipients"`
	Title      string   `mapstructure:"title"`
	Body       string   `mapstructure:"body"`
	Attachs    []string `mapstructure:"attachs"`
}

func (sender *Sender) Send(email *Email) (err error) {
	for i := range email.Recipients {
		if _, err = mail.ParseAddress(email.Recipients[i]); err != nil {
			return err
		}
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", sender.Address)
	msg.SetHeader("To", email.Recipients...)
	msg.SetHeader("Subject", email.Title)
	msg.SetBody("text/html", email.Body)

	for i := range email.Attachs {
		msg.Attach(email.Attachs[i])
	}

	return gomail.NewDialer(
		sender.SMTPAddr, sender.SMTPPort, sender.Address, sender.Password,
	).DialAndSend(msg)
}
