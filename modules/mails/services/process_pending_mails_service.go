package mails_services

import (
	"strings"
	"sync"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// ProcessPendingMailsService is a service that checks the mails table and processes pending messages
type ProcessPendingMailsService struct {
	MailUtil mail.IMailUtil
	OrmUtil  orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s ProcessPendingMailsService) Execute() error {

	pendingMails := []*mails_models.TransientMailMessage{}

	s.OrmUtil.Db().Where("processed = false").Find(&pendingMails)

	var wg sync.WaitGroup

	for _, dbmsg := range pendingMails {

		wg.Add(1)

		go func(msg *mails_models.TransientMailMessage) {

			sendErr := s.MailUtil.Send(mail.Message{
				FromAddress: msg.FromAddress,
				FromName:    msg.FromName,
				To:          strings.Split(msg.To, ","),
				Cc:          strings.Split(msg.Cc, ","),
				Cco:         strings.Split(msg.Cco, ","),
				Subject:     msg.Subject,
				HTML:        msg.Html,
			})

			if sendErr == nil {
				msg.Processed = true
				s.OrmUtil.Db().Save(msg)
			}

			wg.Done()
		}(dbmsg)

	}

	wg.Wait()

	return nil
}
