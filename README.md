# api-server-demo
This is a HTTP Server written in Go for demonstration purposes only.
The functionality currently supported includes:
- adding user
- getting user

### API endpoints
- ```GET /api/v1/user```
- ```POST /api/v1/user```

### application dependencies
API depends on the following:
- Redis - simple caching
- Postgresql - relational database

### configuration
create .env file with the following content:
```shell
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
API_PORT=8080
CACHE_TTL=30m
REDIS_ADDRESS=:6379
```

### run
To run the application with all dependencies run this command:
- ```make docker-compose-up```

### documentation
- [Openapi schema](docs/openapi.yaml)
- [Postman collection](docs/UserAPI.postman_collection.json)
 
