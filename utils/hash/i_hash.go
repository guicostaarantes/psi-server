package hash

// IHashUtil is an abstraction for a utility that hashes strings for data security (like passwords)
type IHashUtil interface {
	Hash(plain string) (string, error)
	Compare(plain string, hashed string) error
	GetWrongPasswordError() string
}
