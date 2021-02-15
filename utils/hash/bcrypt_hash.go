package hash

import (
	"errors"

	"github.com/guicostaarantes/psi-server/utils/logging"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost               int
	wrongPasswordError string
	loggingUtil        logging.ILoggingUtil
}

func (h bcryptHasher) Hash(plain string) (string, error) {
	hashedBytes, hashErr := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if hashErr != nil {
		h.loggingUtil.Error("e3092a73", hashErr)
		return "", errors.New("internal server error")
	}
	return string(hashedBytes), hashErr
}

func (h bcryptHasher) Compare(plain string, hashed string) error {
	compareErr := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if compareErr != nil {
		if compareErr.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return errors.New(h.wrongPasswordError)
		}
		if compareErr.Error() == "crypto/bcrypt: hashedSecret too short to be a bcrypted password" {
			return errors.New(h.wrongPasswordError)
		}
		h.loggingUtil.Error("1f3b1849", compareErr)
		return errors.New("internal server error")
	}
	return nil
}

func (h bcryptHasher) GetWrongPasswordError() string {
	return h.wrongPasswordError
}

// BcryptHashUtil is an implementation of IHashUtil that uses golang.org/x/crypto/bcrypt
var BcryptHashUtil = bcryptHasher{
	cost:               8,
	wrongPasswordError: "hashedPassword is not the hash of the given password",
	loggingUtil:        logging.PrintLogUtil,
}

// WeakBcryptHashUtil is an implementation of IHashUtil that uses golang.org/x/crypto/bcrypt and minimum cost for testing
var WeakBcryptHashUtil = bcryptHasher{
	cost:               4,
	wrongPasswordError: "hashedPassword is not the hash of the given password",
	loggingUtil:        logging.PrintLogUtil,
}
