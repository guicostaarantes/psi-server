package file_storage

// IFileStorageUtil is an abstraction for a utility that stores files
type IFileStorageUtil interface {
	WriteFile(data []byte) (string, error)
	ReadFile(name string) ([]byte, error)
}
