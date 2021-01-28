package serializing

import (
	"encoding/json"
	"errors"

	"github.com/guicostaarantes/psi-server/utils/logging"
)

type jsonSerializer struct {
	loggingUtil logging.ILoggingUtil
}

func (j jsonSerializer) BytesToVariable(bytes []byte, receiver interface{}) error {
	jsonErr := json.Unmarshal(bytes, receiver)

	if jsonErr != nil {
		j.loggingUtil.Error("c36dcec7", jsonErr)
		return errors.New("internal server error")
	}

	return nil
}

func (j jsonSerializer) VariableToBytes(provider interface{}) ([]byte, error) {
	response, jsonErr := json.Marshal(provider)

	if jsonErr != nil {
		j.loggingUtil.Error("3d5252ea", jsonErr)
		return nil, errors.New("internal server error")
	}

	return response, nil
}

// JSONSerializingUtil is an implementation of ISerializingUtil that uses enconding/json
var JSONSerializingUtil = jsonSerializer{
	loggingUtil: logging.PrintLogUtil,
}
