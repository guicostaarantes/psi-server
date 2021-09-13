package files_services

import (
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/file_storage"
)

// UploadAvatarFileService is a service that uploads a new avatar file to the storage
type UploadAvatarFileService struct {
	DatabaseUtil    database.IDatabaseUtil
	FileStorageUtil file_storage.IFileStorageUtil
}

// Execute is the method that runs the business logic of the service
func (s UploadAvatarFileService) Execute(userID string, data []byte) (string, error) {
	fileName, writeErr := s.FileStorageUtil.WriteFile(data)

	return fileName, writeErr
}
