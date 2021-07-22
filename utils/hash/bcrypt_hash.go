package hash

import (
	"errors"

	"github.com/guicostaarantes/psi-server/utils/logging"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHashUtil struct {
	Cost        int
	LoggingUtil logging.ILoggingUtil
}

func (h BcryptHashUtil) Hash(plain string) (string, error) {
	hashedBytes, hashErr := bcrypt.GenerateFromPassword([]byte(plain), h.Cost)
	if hashErr != nil {
		h.LoggingUtil.Error("e3092a73", hashErr)
		return "", errors.New("internal server error")
	}
	return string(hashedBytes), hashErr
}

func (h BcryptHashUtil) Compare(plain string, hashed string) error {
	compareErr := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if compareErr != nil {
		if compareErr.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return errors.New("hashedPassword is not the hash of the given password")
		}
		if compareErr.Error() == "crypto/bcrypt: hashedSecret too short to be a bcrypted password" {
			return errors.New("hashedPassword is not the hash of the given password")
		}
		h.LoggingUtil.Error("1f3b1849", compareErr)
		return errors.New("internal server error")
	}
	return nil
}

func (h BcryptHashUtil) GetWrongPasswordError() string {
	return "hashedPassword is not the hash of the given password"
}
