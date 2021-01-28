package serializing

// ISerializingUtil is an abstraction for a utility that transform bytes into Go variables and vice-versa
type ISerializingUtil interface {
	BytesToVariable(bytes []byte, receiver interface{}) error
	VariableToBytes(provider interface{}) ([]byte, error)
}
