package gdstorage

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
)

type DriveFile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	MimeType  string    `json:"mimeType"`
	CreatedAt time.Time `json:"createdAt"`
}

type StoreFileInput struct {
	Name       string // * New name from system, not real file name
	FileHeader *multipart.FileHeader
	FileSource *multipart.File
}

type GoogleDriveStorage interface {
	GetAppStorages() ([]DriveFile, error)
	CreateAppStorage() (DriveFile, error)
	GetDirectory(dirID string) (DriveFile, error)
	StoreFile(file *StoreFileInput, parentID string) (string, error)
	StoreFiles(files []*StoreFileInput, parentID string) ([]string, error)
	DeleteFile(fileID string) error
	DeleteFiles(fileIDs []string) error
}

type googleDriveStorage struct {
	service *drive.Service
}

func New(srv *drive.Service) GoogleDriveStorage {
	return &googleDriveStorage{srv}
}

func formatDriveFile(f *drive.File) (DriveFile, error) {
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

func GetURL(id string) string {
	return "https://drive.google.com/uc?id=" + id
}

// * =========================================================================================
// * API Functions ===========================================================================
// * =========================================================================================

// * Get all app storages in root directory
func (s *googleDriveStorage) GetAppStorages() ([]DriveFile, error) {
	fileList := []DriveFile{}

	driveFileList, err := s.service.Files.List().Q("'root' in parents and mimeType='application/vnd.google-apps.folder' and name contains 'storage_'").Fields("files(id, name, webViewLink, mimeType, createdTime)").Do()

	if err != nil {
		return fileList, err
	}

	for _, file := range driveFileList.Files {
		formattedDriveFile, err := formatDriveFile(file)

		if err != nil {
			return fileList, err
		}

		fileList = append(fileList, formattedDriveFile)
	}

	return fileList, nil
}

// * Create a new app storage, put the id to .env file by the name DRIVE_APP_DIR_ID
func (s *googleDriveStorage) CreateAppStorage() (DriveFile, error) {
	appName := fmt.Sprintf("storage_%s_%s", os.Getenv("GOOGLE_PROJECT_ID"), os.Getenv("APP_NAME"))

	appDir, err := s.service.Files.Create(&drive.File{Name: appName, MimeType: "application/vnd.google-apps.folder", Parents: []string{"root"}}).Do()

	if err != nil {
		return DriveFile{}, err
	}

	// * Set to read only permission for anyone
	_, err = s.service.Permissions.Create(appDir.Id, &drive.Permission{Type: "anyone", Role: "reader"}).Do()

	if err != nil {
		return DriveFile{}, err
	}

	// ? Adding your real email so that you can easily organize sub-folders in the website
	if em := os.Getenv("DRIVE_ORGANIZER_EMAIL"); em != "" {
		_, err = s.service.Permissions.Create(appDir.Id, &drive.Permission{Type: "user", Role: "writer", EmailAddress: em}).Do()

		if err != nil {
			return DriveFile{}, err
		}
	} else {
		return DriveFile{}, errors.New("Missing DRIVE_ORGANIZER_EMAIL in .env")
	}

	return DriveFile{ID: appDir.Id, Name: appDir.Name}, nil
}

// * Get directory by id
func (s *googleDriveStorage) GetDirectory(dirID string) (DriveFile, error) {
	// * If you provides appdir id, check by the id
	appDir, err := s.service.Files.Get(dirID).Fields("id, name, webViewLink, mimeType, createdTime").Do()

	// * Check if err is caused by something other than not found
	if err != nil {
		e := strings.Split(err.Error(), ", ")
		if e[len(e)-1] != "notFound" {
			return DriveFile{}, nil
		}
		return DriveFile{}, err
	}

	// * Format output data
	formattedDriveFile, err := formatDriveFile(appDir)
	if err != nil {
		return formattedDriveFile, err
	}

	return formattedDriveFile, nil
}

// * Store a file
// TODO - Currently only accepts image files
func (s *googleDriveStorage) StoreFile(file *StoreFileInput, parentID string) (string, error) {
	// * Is parent directory available?
	parentDir, err := s.GetDirectory(parentID)
	if err != nil {
		return "", err
	} else if err == nil && parentDir.ID == "" {
		return "", errors.New(fmt.Sprintf("Unable to find parent directory with id of: %s", parentID))
	}

	// * Extract the file from multipart file header
	src, err := file.FileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// * Building the file data
	newFileName := fmt.Sprintf("%s_%d", file.Name, time.Now().UnixNano())
	fileExt := (strings.Split(filepath.Ext(file.FileHeader.Filename), "."))[1]
	fullFileName := fmt.Sprintf("%s.%s", newFileName, fileExt)
	mimeType := "image/"

	if strings.ToLower(fileExt) == "jpeg" || strings.ToLower(fileExt) == "jpg" {
		mimeType += "jpeg"
	} else if strings.ToLower(fileExt) == "png" {
		mimeType += "png"
	} else {
		return "", errors.New(fmt.Sprintf("Unable to store files with type %s", fileExt))
	}

	// * Store the file
	driveFile, err := s.service.Files.Create(&drive.File{Name: fullFileName, MimeType: mimeType, Parents: []string{parentDir.ID}}).Media(src).Do()
	if err != nil {
		return "", err
	}

	return driveFile.Id, nil
}

// * Store multiple files
func (s *googleDriveStorage) StoreFiles(files []*StoreFileInput, parentID string) ([]string, error) {
	// * Prepare the slice
	driveFiles := []string{}

	// * Is parent directory available?
	parentDir, err := s.GetDirectory(parentID)
	if err != nil {
		return driveFiles, err
	} else if err == nil && parentDir.ID == "" {
		return driveFiles, errors.New(fmt.Sprintf("Unable to find parent directory with id of: %s", parentID))
	}

	// * Extract the file from multipart file header
	// * So that if there's error at some point here, the function returns
	for _, file := range files {
		src, err := file.FileHeader.Open()
		if err != nil {
			return driveFiles, err
		}

		// * Needs to be assigned via pointer so that it saves. It won't save using normal assignment
		file.FileSource = &src

		defer src.Close()
	}

	for _, file := range files {
		// * Building the file data
		newFileName := fmt.Sprintf("%s_%d", file.Name, time.Now().UnixNano())
		fileExt := (strings.Split(filepath.Ext(file.FileHeader.Filename), "."))[1]
		fullFileName := fmt.Sprintf("%s.%s", newFileName, fileExt)
		mimeType := "image/"

		if strings.ToLower(fileExt) == "jpeg" || strings.ToLower(fileExt) == "jpg" {
			mimeType += "jpeg"
		} else if strings.ToLower(fileExt) == "png" {
			mimeType += "png"
		} else {
			return driveFiles, errors.New(fmt.Sprintf("Unable to store files with type %s", fileExt))
		}

		// * Store each file
		driveFile, err := s.service.Files.Create(&drive.File{Name: fullFileName, MimeType: mimeType, Parents: []string{parentDir.ID}}).Media(*file.FileSource).Do()
		if err != nil {
			return driveFiles, err
		}

		driveFiles = append(driveFiles, driveFile.Id)
	}

	return driveFiles, nil
}

// * Delete a file
func (s *googleDriveStorage) DeleteFile(fileID string) error {
	if err := s.service.Files.Delete(fileID).Do(); err != nil {
		e := strings.Split(err.Error(), ", ")
		if e[len(e)-1] == "notFound" {
			return errors.New(fmt.Sprintf("Unable to find file with ID %s", fileID))
		}

		return err
	}

	return nil
}

// * Delete multiple files
func (s *googleDriveStorage) DeleteFiles(fileIDs []string) error {
	for _, fileID := range fileIDs {
		if err := s.service.Files.Delete(fileID).Do(); err != nil {
			e := strings.Split(err.Error(), ", ")
			if e[len(e)-1] == "notFound" {
				return errors.New(fmt.Sprintf("Unable to find file with ID %s", fileID))
			}

			return err
		}
	}

	return nil
}
