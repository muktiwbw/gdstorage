# Go Google Drive Storage API
Simple Go Google Drive API for free webhosters (e.g. Heroku) which the sole purpose is to store generic files such as images and pdf/docs input by users. But as of now it can only store image files (jpg, jpeg, png).

## 1. Installation
Install package.
```
$ go get github.com/muktiwbw/gdstorage
```

## 2. Configurations
There are couple of preparations to do before you can use it.

### Google Cloud Project
1. Create a new project in [GCP dashboard](https://console.cloud.google.com/)
2. Go to Navigation menu > APIs & Services > Credentials
3. Create credentials and choose Service account, fill in the Service account details (you can skip the optionals)
4. Click the service account you just created 
5. Go to Keys tab, click Add key > Create new key, you will get a JSON file containing your service account data

### Environment Variables
1. `APP_NAME`, your app name, it will be used for your storage root directory name
2. `GOOGLE_ACCOUNT_SERVICE_JSON`, paste the content of the JSON file you just downloaded earlier as a string (make sure to remove the spaces just in case)
3. `DRIVE_APP_DIR_ID`, the ID of your storage root directory (you will get it after creating one, will be explained later)
4. `DRIVE_ORGANIZER_EMAIL`, your personal email so that you can organize files and folders from your drive's **Shared with me**
5. `GOOGLE_PROJECT_ID`, no need to set it yourself, it will be set automatically once you create a service instance

```
APP_NAME="Super Cool App"
GOOGLE_ACCOUNT_SERVICE_JSON=
DRIVE_APP_DIR_ID=
DRIVE_ORGANIZER_EMAIL=example@gmail.com
```

## 3. Service Initialization
Create a service to access functions.
```go
func main() {
  // Load .env
  if err := godotenv.Load(); err != nil {
    log.Fatalf("Error loading .env: %v\n", err)
  }

  // Create a Google API service
  srv, err := gdstorage.NewStorageService()
  if err != nil {
    log.Fatalf("Error creating storage service: %v\n", err)
  }

  gds := gdstorage.New(srv)

  // List all of your app storages in current GCP project
  storages, err := gds.GetAppStorages()
  if err != nil {
    fmt.Println(err.Error())

    return
  }
  
  fmt.Println(storages)
  
  // ...
}
```

## 4. Storage Functions
### GetAppStorages
List all of your app storages in current GCP project
```go
storages, err := gds.GetAppStorages()
if err != nil {
  fmt.Println(err.Error())

  return
}

fmt.Println(storages)
```
Will output something like this
```json
[
  {
    "id": "xxxxxxxxxxxxxxxx",
    "name": "storage_super-cool-app_Super Cool App",
    "url": "https://url-to-storage.dir/xxxxxxxxxxxxxxxx",
    "mimeType": "application/vnd.google-apps.folder",
    "createdAt": "timestamp"
  },
  {
    "id": "yyyyyyyyyyyyyyyyy",
    "name": "storage_semi-cool-app_Semi Cool App",
    "url": "https://url-to-storage.dir/yyyyyyyyyyyyyyyyy",
    "mimeType": "application/vnd.google-apps.folder",
    "createdAt": "timestamp"
  },
  {
    "id": "zzzzzzzzzzzzzzzzz",
    "name": "storage_generic-app_Generic App",
    "url": "https://url-to-storage.dir/zzzzzzzzzzzzzzzzz",
    "mimeType": "application/vnd.google-apps.folder",
    "createdAt": "timestamp"
  }
]
```

### CreateAppStorage
Create a storage root directory for current GCP project.
```go
storage, err := gds.CreateAppStorage()
if err != nil {
  fmt.Println(err.Error())

  return
}

fmt.Println(storage)
```
It will give you the directory `id` which you need to assign to `DRIVE_APP_DIR_ID` environment variable. Before you run it make sure to set `APP_NAME` because your storage name will need app name in it.
```json
{
  "id": "xxxxxxxxxxxxxxxx",
  "name": "storage_super-cool-app_Super Cool App",
  "url": "",
  "mimeType": "",
  "createdAt": "timestamp"
}
```

### GetDirectory
Returns a directory by its `id`. It can be used to check if certain parent directory exists or not.
```go
parentID := "xxxxxxxxxxxxxxxxxx"

storage, err := gds.GetDirectory(parentID)
if err != nil {
  fmt.Println(err.Error())

  return
}

// You can check if it exists or not
if storage.ID != "" {
  fmt.Println(fmt.Sprintf("Directory with id %s exists.", storage.ID))
} else {
  fmt.Println(fmt.Sprintf("Directory with id %s does not exist.", storage.ID))
}

fmt.Println(storage)
```
```json
{
  "id": "xxxxxxxxxxxxxxxx",
  "name": "storage_super-cool-app_Super Cool App",
  "url": "https://url-to-storage.dir/xxxxxxxxxxxxxxxx",
  "mimeType": "application/vnd.google-apps.folder",
  "createdAt": "timestamp"
}
```

### StoreFile
Store a single file to parent directory.
```go
fileName := "test-file"
parentID := "xxxxxxxxxxxxxxxxxx"

// file is *multipart.FileHeader
fileID, err := gds.StoreFile(&gdstorage.StoreFileInput{Name: fileName, FileHeader: file}, parentID)
if err != nil {
  fmt.Println(err.Error())

  return
}

// You can get the direct URL using gdstorage.GetURL(id), this is what you want to save in DB
fmt.Printf("Successfully created file. The URL is: %s", gdstorage.GetURL(fileID))
```

### StoreFiles
Store a single file to parent directory.
```go
fileInputs := []*gdstorage.StoreFileInput{}

// files is []*multipart.FileHeader
for i, file := range files {
  fileInputs = append(fileInputs, &gdstorage.StoreFileInput{Name: fmt.Sprintf("test-multiple-file-%d", i), FileHeader: file})
}

parentID := "xxxxxxxxxxxxxxxxxx"

fileIDs, err := gds.StoreFiles(fileInputs, parentID)
if err != nil {
  fmt.Println(err.Error())

  return
}

fileURLs := []string{}
for _, fileID := range fileIDs {
  fileURLs = append(fileURLs, gdstorage.GetURL(fileID))
}

fmt.Printf("Successfully created multiple files. The URLs are: %v", fileURLs)
```

### DeleteFile
Delete a file by id.
```go
fileID := "xxxxxxxxxxxxxxxxxx"

if err := gds.DeleteFile(fileID); err != nil {
  fmt.Println(err.Error())

  return
}

fmt.Printf("Successfully deleted file with id: %s", fileID)
```

### DeleteFiles
Delete multiple files by ids.
```go
fileIDs := []string{"xxxxxxxxxxxxxxxx", "yyyyyyyyyyyyyyyyyyyyyy", "zzzzzzzzzzzzzzzzzzz"}

if err := gds.DeleteFiles(fileIDs); err != nil {
  fmt.Println(err.Error())

  return
}

fmt.Printf("Successfully deleted files with ids: %v", fileIDs)
```

### GetURL
Get direct file URL to be displayed in your web.
```go
url := gdstorage.GetURL(fileID)
fmt.Println(url)
```
