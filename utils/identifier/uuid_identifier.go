package identifier

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/guicostaarantes/psi-server/utils/logging"
)

type uuidIdentifer struct {
	loggingUtil logging.ILoggingUtil
}

func (u uuidIdentifer) GenerateIdentifier() ([]byte, string, error) {
	uu, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		u.loggingUtil.Error("0ca94e3b", uuidErr)
		return nil, "", errors.New("internal server error")
	}
	return uu.Bytes(), uu.String(), uuidErr
}

// UUIDIdentifierUtil is an implementation of IIdentifierUtil that uses github.com/gofrs/uuid
var UUIDIdentifierUtil = uuidIdentifer{
	loggingUtil: logging.PrintLogUtil,
}
