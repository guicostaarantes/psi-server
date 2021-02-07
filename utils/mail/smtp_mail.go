package mail

import (
	"errors"
	"os"
	"strconv"

	"github.com/guicostaarantes/psi-server/utils/logging"
	"gopkg.in/gomail.v2"
)

type smtpMailer struct {
	loggingUtil logging.ILoggingUtil
}

func (s smtpMailer) GetMockedMessages() (*[]map[string]interface{}, error) {
	return nil, errors.New("this is not a mock implementation")
}

func (s smtpMailer) Send(msg Message) error {
	envSMTPPort, _ := strconv.Atoi(os.Getenv("PSI_SMTP_PORT"))

	dialer := gomail.NewDialer(os.Getenv("PSI_SMTP_HOST"), envSMTPPort, os.Getenv("PSI_SMTP_USERNAME"), os.Getenv("PSI_SMTP_PASSWORD"))

	gomsg := gomail.NewMessage()
	gomsg.SetHeader("From", gomsg.FormatAddress(msg.FromAddress, msg.FromName))
	gomsg.SetHeader("To", msg.To...)
	gomsg.SetHeader("Cc", msg.Cc...)
	gomsg.SetHeader("Cco", msg.Cco...)
	gomsg.SetHeader("Subject", msg.Subject)
	gomsg.SetBody("text/html", msg.HTML)

	sendErr := dialer.DialAndSend(gomsg)
	if sendErr != nil {
		s.loggingUtil.Error("fb59180c", sendErr)
		return errors.New("internal server error")
	}

	return nil
}

// SMTPMailUtil is an implementation of IMailUtil that uses SMTP via gopkg.in/gomail.v2
var SMTPMailUtil = smtpMailer{
	loggingUtil: logging.PrintLogUtil,
}
