package gdstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type AccountService struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func NewStorageService() (*drive.Service, error) {
	// * Get the account service json content from .env
	serviceAccountJsonContent := os.Getenv("GOOGLE_ACCOUNT_SERVICE_JSON")
	if serviceAccountJsonContent == "" {
		return &drive.Service{}, errors.New("Missing account service json data")
	}

	// * Set project id env
	var accountService AccountService

	// * Parse json content into an object
	if err := json.Unmarshal([]byte(serviceAccountJsonContent), &accountService); err != nil {
		return &drive.Service{}, errors.New(fmt.Sprintf("Error parsing account service JSON content: %v", err))
	}

	// * Set project id
	if err := os.Setenv("GOOGLE_PROJECT_ID", accountService.ProjectID); err != nil {
		return &drive.Service{}, errors.New(fmt.Sprintf("Error setting project id: %v", err))
	}

	// * Getting working directory
	wd, err := os.Getwd()
	log.Printf("Working dir: %s", wd)
	if err != nil {
		return &drive.Service{}, errors.New(fmt.Sprintf("Error retrieving working directory: %v", err))
	}

	// * Generate a JSON file based on your service account data stored in env
	if _, err := os.Stat(filepath.Join(wd, "svracc.json")); err != nil && os.IsNotExist(err) {
		err := ioutil.WriteFile(filepath.Join(wd, "svracc.json"), []byte(serviceAccountJsonContent), os.ModePerm)

		if err != nil {
			return &drive.Service{}, errors.New(fmt.Sprintf("Error writing service account file: %v", err))
		}
	} else if err != nil && !os.IsNotExist(err) {
		return &drive.Service{}, errors.New(fmt.Sprintf("Error loading service account file: %v", err))
	} else if err == nil {
		// * If file actually exists
		jsonContent, err := os.ReadFile(filepath.Join(wd, "svracc.json"))
		if err != nil {
			return &drive.Service{}, errors.New(fmt.Sprintf("Error reading service account file: %v", err))
		}

		// * Check if the content is the same as env or not
		// * If false, override with the env content
		if string(jsonContent) != serviceAccountJsonContent {
			err := ioutil.WriteFile(filepath.Join(wd, "svracc.json"), []byte(serviceAccountJsonContent), os.ModePerm)

			if err != nil {
				return &drive.Service{}, errors.New(fmt.Sprintf("Error overriding service account file: %v", err))
			}
		}
	}

	// * Create a new drive service
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(filepath.Join(wd, "svracc.json")))

	if err != nil {
		return &drive.Service{}, errors.New(fmt.Sprintf("Error creating drive service api: %v", err))
	}

	return srv, nil
}
