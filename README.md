# go-mongo-modular-sample

#### To run the application
```shell
go run main.go
```

## Architecture
### Architecture
![Architecture](https://camo.githubusercontent.com/c17d4dfaab39cf7223f7775c9e973bb936e4169e8bd0011659e83cec755c8f26/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a42526b437272622d5f417637395167737142556b48672e706e67)

### Architecture with standard features: config, health check, logging, middleware log tracing
![Architecture with standard features: config, health check, logging, middleware log tracing](https://camo.githubusercontent.com/fa1158e7f94bf96e09aef42fcead23366839baf71190133d5df10f3006b2e041/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a6d494e3344556569365676316c755a376747727655412e706e67)

#### [core-go/search](https://github.com/core-go/search)
- Build the search model at http handler
- Build dynamic SQL for search
  - Build SQL for paging by page index (page) and page size (limit)
  - Build SQL to count total of records
### Search users: Support both GET and POST 
#### POST /users/search
##### *Request:* POST /users/search
In the below sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending
```json
{
    "page": 1,
    "limit": 20,
    "sort": "phone,-id",
    "email": "tony",
    "dateOfBirth": {
        "min": "1953-11-16T00:00:00+07:00",
        "max": "1976-11-16T00:00:00+07:00"
    }
}
```
##### GET /users/search?page=1&limit=2&email=tony&dateOfBirth.min=1953-11-16T00:00:00+07:00&dateOfBirth.max=1976-11-16T00:00:00+07:00&sort=phone,-id
In this sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending

#### *Response:*
- total: total of users, which is used to calculate numbers of pages at client 
- list: list of users
```json
{
    "list": [
        {
            "id": "ironman",
            "username": "tony.stark",
            "email": "tony.stark@gmail.com",
            "phone": "0987654321",
            "dateOfBirth": "1963-03-24T17:00:00Z"
        }
    ],
    "total": 1
}
```

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource, refer to [core](https://github.com/core-go/core) and [mongo](https://github.com/core-go/mongo)  
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "mongo": {
            "status": "UP"
        }
    }
}
```

## API design for users
#### *Resource:* users

### Get all users
#### *Request:* GET /users
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "username": "peter.parker",
        "email": "peter.parker@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one user by id
#### *Request:* GET /users/:id
```shell
GET /users/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```

### Create a new user
#### *Request:* POST /users 
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 4: error
```json
{
    "status": 1,
    "value": {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T00:00:00+07:00"
    }
}
```
#### *Fail case sample:* 
- Request:
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett",
    "phone": "0987654321a",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
- Response: in this below sample, email and phone are not valid
```json
{
    "status": 4,
    "errors": [
        {
            "field": "email",
            "code": "email"
        },
        {
            "field": "phone",
            "code": "phone"
        }
    ]
}
```

### Update one user by id
#### *Request:* PUT /users/:id
```shell
PUT /users/wolverine
```
```json
{
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 2: version error, 4: error
```json
{
    "status": 1,
    "value": {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T00:00:00+07:00"
    }
}
```

### Patch one user by id
Perform a partial update of user. For example, if you want to update 2 fields: email and phone, you can send the request body of below.
#### *Request:* PATCH /users/:id
```shell
PATCH /users/wolverine
```
```json
{
    "email": "james.howlett@gmail.com",
    "phone": "0987654321"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 2: version error, 4: error
```json
{
    "status": 1,
    "value": {
        "email": "james.howlett@gmail.com",
        "phone": "0987654321"
    }
}
```

#### Problems for patch
If we pass a struct as a parameter, we cannot control what fields we need to update. So, we must pass a map as a parameter.
```go
type UserService interface {
    Update(ctx context.Context, user *User) (int64, error)
    Patch(ctx context.Context, user map[string]interface{}) (int64, error)
}
```
We must solve 2 problems:
1. At http handler layer, we must convert the user struct to map, with json format, and make sure the nested data types are passed correctly.
2. At service layer or repository layer, from json format, we must convert the json format to database format (in this case, we must convert to bson of Mongo)

#### Solutions for patch  
1. At http handler layer, we use [core-go/core](https://github.com/core-go/core), to convert the user struct to map, to make sure we just update the fields we need to update
```go
import "github.com/core-go/core"

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
    var user User
    userType := reflect.TypeOf(user)
    _, jsonMap := core.BuildMapField(userType)
    body, _ := core.BuildMapAndStruct(r, &user)
    json, er1 := core.BodyToJson(r, user, body, ids, jsonMap, nil)

    res, er2 := h.service.Patch(r.Context(), json)
    if er2 != nil {
        http.Error(w, er2.Error(), http.StatusInternalServerError)
        return
    }
    respond(w, res)
}
```

2. At service layer or repository layer, we use [core-go/mongo](https://github.com/core-go/mongo), to convert from json to bson 
```go
import mgo "github.com/core-go/mongo"

func (p *MongoUserService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
    userType := reflect.TypeOf(User{})
    maps := mgo.MakeBsonMap(userType)
    filter := mgo.BuildQueryByIdFromMap(user, "id")
    bson := mgo.MapToBson(user, maps)
    return mgo.PatchOne(ctx, p.Collection, bson, filter)
}
```

### Delete a new user by id
#### *Request:* DELETE /users/:id
```shell
DELETE /users/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include HealthHandler, HealthChecker, MongoHealthChecker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): log and log middleware

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "mongo": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(root.Mongo.Uri))
    if err != nil {
        return nil, err
    }
    db := client.Database(root.Mongo.Database)

    mongoChecker := mongo.NewHealthChecker(db)
    healthHandler := health.NewHealthHandler(mongoChecker)
```

To handler routing
```go
    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler.Check).Methods("GET")
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package main

import "github.com/core-go/config"

type Root struct {
    DB DatabaseConfig `mapstructure:"db"`
}

type DatabaseConfig struct {
    Driver                 string `mapstructure:"driver"`
    DataSourceName string `mapstructure:"data_source_name"`
}

func main() {
    var conf Root
    err := config.Load(&conf, "configs/config")
    if err != nil {
        panic(err)
    }
}
```

### core-go/log *&* core-go/log/middleware
```go
import (
	"github.com/core-go/config"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/gorilla/mux"
)

func main() {
	var conf app.Root
	config.Load(&conf, "configs/config")

	r := mux.NewRouter()

	log.Initialize(conf.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewLogger()
	r.Use(mid.Logger(conf.MiddleWare, log.InfoFields, logger))
	r.Use(mid.Recover(log.ErrorMsg))
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```