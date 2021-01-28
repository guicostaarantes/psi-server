package identifier

// IIdentifierUtil is an abstraction for a utility that creates virtually unique identifiers
type IIdentifierUtil interface {
	GenerateIdentifier() ([]byte, string, error)
}
