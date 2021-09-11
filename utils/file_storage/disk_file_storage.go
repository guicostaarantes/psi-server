package file_storage

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guicostaarantes/psi-server/utils/logging"
)

type DiskFileStorageUtil struct {
	BaseFolder  string
	LoggingUtil logging.ILoggingUtil
}

func (u DiskFileStorageUtil) WriteFile(data []byte) (string, error) {
	fileName := fmt.Sprintf("%x", sha256.Sum256(data))
	fileLocation := filepath.Join(u.BaseFolder, fileName)

	writeErr := os.WriteFile(fileLocation, data, 0644)
	if writeErr != nil {
		u.LoggingUtil.Error("9a429953", writeErr)
		return "", writeErr
	}

	return fileName, nil
}

func (u DiskFileStorageUtil) ReadFile(name string) ([]byte, error) {
	fileLocation := filepath.Join(u.BaseFolder, name)

	data, readErr := os.ReadFile(fileLocation)
	if readErr != nil {
		if readErr.Error() == fmt.Sprintf(`open /data/files/%s: no such file or directory`, name) {
			return []byte{}, nil
		} else {
			u.LoggingUtil.Error("486f275d", readErr)
			return []byte{}, readErr
		}
	}

	return data, nil
}
