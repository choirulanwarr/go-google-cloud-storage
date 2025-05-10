## GCS Integration

This is google cloud storage integration using Client with ADC or with Credential file .json key form Google Cloud. Using golang gin Gonic and postgres setup documentation using GO.

Help me to improve this code more flexible and esy to use

API 
- `POST` api/v1/upload
- `GET` api/v1/download?path=
- `GET` api/v1/list
- `GET` api/v1/list/:folderName
- `POST` api/v1/bucket/create
- `DELETE` api/v1/delete

**NEXT update ...**

- `POST` api/v1/bucket/config

## Feature golang

- Gin Gonic Framework 
```
go get -u github.com/gin-gonic/gin
```
- Gin Contrib Cors 
```
go get -u github.com/gin-contrib/cors
```
- Gorm
```
go get -u gorm.io/gorm
```
- PostgresSQL
```
go get -u gorm.io/driver/postgres
```
- Google Cloud Storage
```
go get -u cloud.google.com/go/storage
```
- Viper
```
go get -u github.com/spf13/viper
```
- Validator Go v10
```
go get -u github.com/go-playground/validator/v10
```
- Logrus
```
go get -u github.com/sirupsen/logrus
```
