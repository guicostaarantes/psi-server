package files_services

import (
	"github.com/guicostaarantes/psi-server/utils/file_storage"
)

// ReadFileService is a service that serves a file from the storage
type ReadFileService struct {
	FileStorageUtil file_storage.IFileStorageUtil
}

// Execute is the method that runs the business logic of the service
func (s ReadFileService) Execute(name string) ([]byte, error) {
	return s.FileStorageUtil.ReadFile(name)
}
