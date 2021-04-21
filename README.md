# Go Google Drive Storage API
Simple Go Google Drive API for free webhosters (e.g. Heroku) which the sole purpose is to store generic files such as images and pdf/docs input by users.

## 1. Installation
Install package.
```
$ go get github.com/muktiwbw/gdstorage
```

## 2. Configurations
There are couple of preparations to do before you can use it.

### 2.a Google Cloud Project
1. Create a new project in [GCP dashboard](https://console.cloud.google.com/)
2. Go to Navigation menu > APIs & Services > Credentials
3. Create credentials and choose Service account, fill in the Service account details (you can skip the optionals)
4. Click the service account you just created 
5. Go to Keys tab, click Add key > Create new key, you will get a JSON file containing your service account data

### 2.b Environment Variables
1. **APP_NAME**, your app name, it will be used for your storage root directory name
2. **GOOGLE_ACCOUNT_SERVICE_JSON**, paste the content of the JSON file you just downloaded earlier as a string (make sure to remove the spaces just in case)
3. **DRIVE_APP_DIR_ID**, the ID of your storage root directory (you will get it after creating one, will be explained later)
4. **DRIVE_ORGANIZER_EMAIL**, your personal email so that you can organize files and folders from your drive's **Shared with me**
5. **GOOGLE_PROJECT_ID**, no need to set it yourself, it will be set automatically once you create a service instance
