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
    "mimeType": "image/png",
    "createdAt": "timestamp"
  },
  {
    "id": "yyyyyyyyyyyyyyyyy",
    "name": "storage_semi-cool-app_Semi Cool App",
    "url": "https://url-to-storage.dir/yyyyyyyyyyyyyyyyy",
    "mimeType": "image/png",
    "createdAt": "timestamp"
  },
  {
    "id": "zzzzzzzzzzzzzzzzz",
    "name": "storage_generic-app_Generic App",
    "url": "https://url-to-storage.dir/zzzzzzzzzzzzzzzzz",
    "mimeType": "image/png",
    "createdAt": "timestamp"
  }
]
```

### CreateAppStorage
Create a storage root directory for current GCP project. It will give you the directory id which you need to assign to `DRIVE_APP_DIR_ID` environment variable. Before you run it make sure to set `APP_NAME`
```go
storage, err := gds.CreateAppStorage()
if err != nil {
  fmt.Println(err.Error())

  return
}

fmt.Println(storage)
```