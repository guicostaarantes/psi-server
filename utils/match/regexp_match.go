package match

import (
	"errors"
	"regexp"

	"github.com/guicostaarantes/psi-server/utils/logging"
)

type regexpMatchUtil struct {
	loggingUtil logging.ILoggingUtil
}

func (r regexpMatchUtil) IsPasswordStrong(password string) error {
	checks := []string{"[A-Z]", "[a-z]", "[0-9]", "[~!@#$%^&*()_+]", ".{8,}"}

	for _, check := range checks {
		match, matchErr := regexp.MatchString(check, password)
		if matchErr != nil {
			r.loggingUtil.Error("78dbbb35", matchErr)
			return errors.New("internal server error")
		}

		if !match {
			return errors.New("weak password")
		}
	}

	return nil
}

func (r regexpMatchUtil) IsEmailValid(email string) error {
	match, matchErr := regexp.MatchString("^[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$", email)
	if matchErr != nil {
		r.loggingUtil.Error("7c54d929", matchErr)
		return errors.New("internal server error")
	}

	if !match {
		return errors.New("invalid email")
	}

	return nil
}

// RegexpMatchUtil is an implementation of IMatchUtil that uses regexp.
var RegexpMatchUtil = regexpMatchUtil{
	loggingUtil: logging.PrintLogUtil,
}
