package tasks

import (
	"github.com/mayron-dev/gowork/internal/email"
	"github.com/mayron-dev/gowork/internal/file"
)

var (
	emailService email.IEmailService
	fileService  file.IFileService
)

func Init() error {
	emailService = email.NewEmailService(&email.EmailConfig{})
	fs, err := file.NewFileService("", "", "", "")
	if err != nil {
		return err
	}
	fileService = fs
	return nil
}
