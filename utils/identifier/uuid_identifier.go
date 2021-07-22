package identifier

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/guicostaarantes/psi-server/utils/logging"
)

type UuidIdentifierUtil struct {
	LoggingUtil logging.ILoggingUtil
}

func (u UuidIdentifierUtil) GenerateIdentifier() ([]byte, string, error) {
	uu, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		u.LoggingUtil.Error("0ca94e3b", uuidErr)
		return nil, "", errors.New("internal server error")
	}
	return uu.Bytes(), uu.String(), uuidErr
}
