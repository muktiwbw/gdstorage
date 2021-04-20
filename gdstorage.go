package gdstorage

import (
	"time"

	"google.golang.org/api/drive/v3"
)

type DriveFile struct {
	ID        string
	Name      string
	URL       string
	MimeType  string
	CreatedAt time.Time
}

func FormatDriveFile(f *drive.File) (DriveFile, error) {
	createdAt, err := time.Parse(time.RFC3339, f.CreatedTime)

	if err != nil {
		return DriveFile{}, err
	}

	return DriveFile{
		ID:        f.Id,
		Name:      f.Name,
		URL:       f.WebViewLink,
		MimeType:  f.MimeType,
		CreatedAt: createdAt,
	}, nil
}

func GetAppStorages(srv *drive.Service) ([]DriveFile, error) {
	fileList := []DriveFile{}

	driveFileList, err := srv.Files.List().Q("'root' in parents and mimeType='application/vnd.google-apps.folder' and name contains 'storage_'").Fields("files(id, name, webViewLink, mimeType, createdTime)").Do()

	if err != nil {
		return fileList, err
	}

	for _, file := range driveFileList.Files {
		formattedDriveFile, err := FormatDriveFile(file)

		if err != nil {
			return fileList, err
		}

		fileList = append(fileList, formattedDriveFile)
	}

	return fileList, nil
}
