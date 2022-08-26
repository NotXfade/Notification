# On Boarding Portal Document Manage Service #

Onboarding Portal Document Manage service is used to upload and access documents

## Requirements ##
1. go (1.13 or above)
2. Postgres Database
3. JWT token
4. NATS Server 2.1 or above

## Documentation ##
1. [Api Documentation on Swagger](openapi.yaml)
2. [Sample Configuration using toml](example.toml)

## Environment Variables  ##
```
//Postgres database
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="Postgres"
export DB_PASS="Postgres"
export DB_IDEAL_CONNECTIONS = "50"

//Service related settings
export ENVIRONMENT="development"
export BUILD_IMAGE="image"
export SERVICE_PORT=""

//NATS related settings
export NATS_URL="localhost:4222"
export NATS_TOKEN="auth token" // either token or user name/password can be used for authentication
export NATS_USERNAME="user"
export NATS_PASSWORD="password"

//jwt Configuration
export PRIVATE_KEY="same among all the backend services"

//AWS Configuration
export AWS_ACCESS_ID=""
export AWS_ACCESS_KEY=""
export AWS_BUCKET_NAME=""
```

## How to run the app ##

### 1. Configuration using environment variables ###

```
i.    Export above all environment variables
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/document-manage --conf=environment`
```

### 2. Configuration using TOML file ###

```
i.    Create a configuration toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/document-manage --conf=toml --file=<path of toml file>`
```
### 3. Run the Test Case ###

```
i.    Run the Test case -> go test ./... -p 1 -v -coverprofile=coverage.out
ii.   To see the Test Coverage Output in html page -> go tool cover -html=coverage.out 
```


`Note :- for any help regarding flags, run this command '$GOPATH/bin/accounts --help'`
