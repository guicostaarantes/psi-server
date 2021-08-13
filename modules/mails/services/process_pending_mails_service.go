package mails_services

import (
	"context"
	"sync"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/mail"
)

// ProcessPendingMailsService is a service that checks the mails table and processes pending messages
type ProcessPendingMailsService struct {
	DatabaseUtil database.IDatabaseUtil
	MailUtil     mail.IMailUtil
}

// Execute is the method that runs the business logic of the service
func (p ProcessPendingMailsService) Execute() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	cursor, findErr := p.DatabaseUtil.FindMany("mails", map[string]interface{}{"processed": false})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(ctx)

	var wg sync.WaitGroup

	for cursor.Next(ctx) {

		dbmsg := &mails_models.TransientMailMessage{}

		decodeErr := cursor.Decode(&dbmsg)
		if decodeErr != nil {
			return decodeErr
		}

		wg.Add(1)

		go func(msg *mails_models.TransientMailMessage) {

			sendErr := p.MailUtil.Send(mail.Message{
				FromAddress: dbmsg.FromAddress,
				FromName:    dbmsg.FromName,
				To:          dbmsg.To,
				Cc:          dbmsg.Cc,
				Cco:         dbmsg.Cco,
				Subject:     dbmsg.Subject,
				HTML:        dbmsg.Html,
			})

			if sendErr == nil {
				dbmsg.Processed = true
				p.DatabaseUtil.UpdateOne("mails", map[string]interface{}{"id": dbmsg.ID}, dbmsg)
			}

			wg.Done()
		}(dbmsg)

	}

	wg.Wait()

	return nil
}
