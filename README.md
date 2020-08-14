# Health building block

Go project to provide rest service for rokwire building block health results.

The service is based on clear hexagonal architecture.

## Set Up

### Prerequisites

MongoDB v4.2.2+

Go v1.13+

### Environment variables
The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Value|Required|Description
---|---|---|---
ROKWIRE_API_KEYS | <value1,value2,value3> | yes | Comma separated list of rokwire api keys
HEALTH_MONGO_AUTH | <mongodb://USER:PASSWORD@HOST:PORT/DATABASE NAME> | yes | MongoDB authentication string. The user must have read/write privileges.
HEALTH_MONGO_DATABASE | < value > | yes | MongoDB database name
HEALTH_MONGO_TIMEOUT | < value > | no | MongoDB timeout in milliseconds. Set default value(500 milliseconds) if omitted
HEALTH_NEWS_RSS_URL | < value > | yes | News RSS url
HEALTH_RESOURCES_URL | < value > | yes | Resources url
HEALTH_SMTP_HOST | < value > | yes | SMTP host
HEALTH_SMTP_PORT | < value > | yes | SMTP port
HEALTH_SMTP_USER | < value > | yes | SMTP user
HEALTH_SMTP_PASSWORD | < value > | yes | SMTP password
HEALTH_EMAIL_FROM | < value > | yes | Email from
HEALTH_EMAIL_TO | < value > | yes | Email to
HEALTH_OIDC_PROVIDER | < value > | yes | OIDC provider
HEALTH_OIDC_APP_CLIENT_ID | < value > | yes | OIDC app client id
HEALTH_OIDC_ADMIN_CLIENT_ID | < value > | yes | OIDC admin client id
HEALTH_PHONE_SECRET | < value > | yes | Phone secret
HEALTH_PROVIDERS_KEY | <value1,value2,value3> | yes | Comma separated list of providers api keys
HEALTH_HOST | < value > | yes | Host
HEALTH_FIREBASE_PROJECT_ID | < value > | yes | Firebase project ID
HEALTH_FIREBASE_AUTH | < value > | yes | Firebase authentication file content
HEALTH_PROFILE_HOST | < value > | yes | Profile building block host
HEALTH_PROFILE_API_KEY | < value > | yes | Profile building block api key

### Run Application

#### Run locally without Docker

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Make the project  
```
$ make
...
▶ building executable(s)… 1.9.0 2020-08-13T10:00:00+0300
```

4. Run the executable
```
$ ./bin/health
```

#### Run locally as Docker container

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Create Docker image  
```
docker build -t health .
```
4. Run as Docker container
```
docker run -e ROKWIRE_API_KEYS -e HEALTH_MONGO_AUTH -e HEALTH_MONGO_DATABASE -e HEALTH_MONGO_TIMEOUT -e HEALTH_NEWS_RSS_URL -e HEALTH_RESOURCES_URL -e HEALTH_SMTP_HOST -e HEALTH_SMTP_PORT -e HEALTH_SMTP_USER -e HEALTH_SMTP_PASSWORD -e HEALTH_EMAIL_FROM -e HEALTH_EMAIL_TO -e HEALTH_OIDC_PROVIDER -e HEALTH_OIDC_APP_CLIENT_ID -e HEALTH_OIDC_ADMIN_CLIENT_ID -e HEALTH_PHONE_SECRET -e HEALTH_PROVIDERS_KEY -e HEALTH_HOST -e HEALTH_FIREBASE_PROJECT_ID -e HEALTH_FIREBASE_AUTH -e HEALTH_PROFILE_HOST -e HEALTH_PROFILE_API_KEY -p 80:80 health
```

#### Tools

##### Run tests
```
$ make tests
```

##### Run code coverage tests
```
$ make cover
```

##### Run golint
```
$ make lint
```

##### Run gofmt to check formatting on all source files
```
$ make checkfmt
```

##### Run gofmt to fix formatting on all source files
```
$ make fixfmt
```

##### Cleanup everything
```
$ make clean
```

##### Run help
```
$ make help
```

##### Generate Swagger docs
```
$ make swagger
```

### Test Application APIs

Verify the service is running as calling the get version API.

#### Call get version API

curl -X GET -i http://localhost/health/version

Response
```
1.9.0
```

## Documentation

The documentation is placed here - https://api-dev.rokwire.illinois.edu/docs/?urls.primaryName=Health%20Building%20Block

Alternativelly the documentation is served by the service on the following url - https://api-dev.rokwire.illinois.edu/health/doc/ui/
