package resolvers

import (
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	users_services "github.com/guicostaarantes/psi-server/modules/users/services"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

var databaseUtil = database.MongoDatabaseUtil
var hashUtil = hash.BcryptHashUtil
var identifierUtil = identifier.UUIDIdentifierUtil
var mailUtil = mail.SMTPMailUtil
var matchUtil = match.RegexpMatchUtil
var serializingUtil = serializing.JSONSerializerUtil
var tokenUtil = token.RngTokenUtil
var secondsToExpire = int64(1800)

var ActivateUserService = &users_services.ActivateUserService{
	DatabaseUtil: databaseUtil,
}

var AuthenticateUserService = &users_services.AuthenticateUserService{
	DatabaseUtil:    databaseUtil,
	HashUtil:        hashUtil,
	IdentifierUtil:  identifierUtil,
	SerializingUtil: serializingUtil,
	TokenUtil:       tokenUtil,
	SecondsToExpire: secondsToExpire,
}

var CreateUserService = &users_services.CreateUserService{
	DatabaseUtil:    databaseUtil,
	IdentifierUtil:  identifierUtil,
	MatchUtil:       matchUtil,
	SerializingUtil: serializingUtil,
}

var CreateUserWithPasswordService = &users_services.CreateUserWithPasswordService{
	DatabaseUtil:    databaseUtil,
	HashUtil:        hashUtil,
	IdentifierUtil:  identifierUtil,
	MatchUtil:       matchUtil,
	SerializingUtil: serializingUtil,
}

var GetUserByIdService = &users_services.GetUserByIdService{
	DatabaseUtil: databaseUtil,
}

var ProcessPendingMailsService = &mails_services.ProcessPendingMailsService{
	DatabaseUtil: databaseUtil,
	MailUtil:     mailUtil,
}

var ValidateUserTokenService = &users_services.ValidateUserTokenService{
	DatabaseUtil:    databaseUtil,
	SerializingUtil: serializingUtil,
	SecondsToExpire: secondsToExpire,
}
