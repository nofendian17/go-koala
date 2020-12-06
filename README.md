### Dev requirements

 - MySQL => v.5.7 
 - Redis => 6.9.0
 - Go => 1.15.5
 
###  Project Setup

- Import schema.sql from folder Database/schema.sql to local enviroment
- Edit .env configuration at rest/.env 
- Import Postman Collection at postman_collection folder

### Test Run

```sh
$ cd rest
$ go mod download #download dependency
$ go run main.go
```

### Rest API test

- Open Url from postman or browser

```sh
http://localhost:8080
```

### Feedback 
```sh
$ beni_anwar@yahoo.com
```