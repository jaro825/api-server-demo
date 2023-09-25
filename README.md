# api-server-demo
This is a HTTP Server written in Go for demonstration purposes only.
The functionality currently supported includes:
- adding user
- getting user

Users are stored in relational database.

The application also uses cache to store users until TTL expires.

When querying for user application first searches its cache
and only if user is not found it queries a database.

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
 
### example output
```shell
./cmd/server/api-server --debug
12:21AM INF api server created
12:22AM WRN request completed bytes=19 ip=127.0.0.1:64453 latency=0.186542 method=POST path=/api/v1/users query= status=404
12:22AM DBG UserAPI.CreateUser saved user in cache: {ID:89757d09-bf3a-4718-b8bf-4ffc2b3d527e Name:jaro}
12:22AM DBG PostgresStore.CreateUser: &{ID:89757d09-bf3a-4718-b8bf-4ffc2b3d527e Name:jaro}
12:22AM INF request completed bytes=60 ip=127.0.0.1:64453 latency=152.815708 method=POST path=/api/v1/user query= status=201
12:22AM DBG UserAPI.GetUser user found in cache: {ID:89757d09-bf3a-4718-b8bf-4ffc2b3d527e Name:jaro}
12:22AM INF request completed bytes=60 ip=127.0.0.1:64453 latency=5.832125 method=GET path=/api/v1/user/89757d09-bf3a-4718-b8bf-4ffc2b3d527e query= status=200
12:22AM WRN request completed bytes=40 ip=127.0.0.1:64453 latency=0.075667 method=GET path=/api/v1/user/xyz query= status=400
```