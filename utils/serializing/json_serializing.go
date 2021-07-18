package serializing

import (
	"encoding/json"
	"errors"

	"github.com/guicostaarantes/psi-server/utils/logging"
)

type JsonSerializingUtil struct {
	LoggingUtil logging.ILoggingUtil
}

func (j JsonSerializingUtil) BytesToVariable(bytes []byte, receiver interface{}) error {
	jsonErr := json.Unmarshal(bytes, receiver)

	if jsonErr != nil {
		j.LoggingUtil.Error("c36dcec7", jsonErr)
		return errors.New("internal server error")
	}

	return nil
}

func (j JsonSerializingUtil) VariableToBytes(provider interface{}) ([]byte, error) {
	response, jsonErr := json.Marshal(provider)

	if jsonErr != nil {
		j.LoggingUtil.Error("3d5252ea", jsonErr)
		return nil, errors.New("internal server error")
	}

	return response, nil
}
