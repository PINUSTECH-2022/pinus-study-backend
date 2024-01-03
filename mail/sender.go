package mail

import (
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	// to := []string{"recipient@example.net"}
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "!\r\n" +
		"\r\n" +
		content + "\r\n")
	return smtp.SendMail(smtpServerAddress, auth, sender.fromEmailAddress, to, msg)

	// e := email.NewEmail()
	// e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	// e.Subject = subject
	// e.HTML = []byte(content)
	// e.To = to
	// e.Cc = cc
	// e.Bcc = bcc

	// for _, f := range attachFiles {
	// 	_, err := e.AttachFile(f)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to attach file %s: %w", f, err)
	// 	}
	// }

	// smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	// return e.Send(smtpServerAddress, smtpAuth)
}
