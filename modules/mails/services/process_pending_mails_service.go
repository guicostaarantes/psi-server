package mails_services

import (
	"context"
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

	cursor, findErr := p.DatabaseUtil.FindMany("psi_db", "mails", map[string]interface{}{"processed": false})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		dbmsg := &mails_models.TransientMailMessage{}

		decodeErr := cursor.Decode(&dbmsg)
		if decodeErr != nil {
			return decodeErr
		}

		sendErr := p.MailUtil.Send(mail.MailMessage{
			FromAddress: dbmsg.FromAddress,
			FromName:    dbmsg.FromName,
			To:          dbmsg.To,
			Cc:          dbmsg.Cc,
			Cco:         dbmsg.Cco,
			Subject:     dbmsg.Subject,
			Html:        dbmsg.Html,
		})
		if sendErr != nil {
			return sendErr
		}

		dbmsg.Processed = true

		updateErr := p.DatabaseUtil.UpdateOne("psi_db", "mails", map[string]interface{}{"id": dbmsg.ID}, dbmsg)
		if updateErr != nil {
			return updateErr
		}

	}

	return nil

}
