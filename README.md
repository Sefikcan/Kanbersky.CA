### Golang Clean Architecture Sample

### Project Full List what has been used:
* [echo](https://github.com/labstack/echo) - Web framework
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [swag](https://github.com/swaggo/swag) - Swagger
* [Docker](https://www.docker.com/) - Docker

### Recommendation for local development most comfortable usage
    make local
    make run

### Docker-compose files:
    docker-compose.yml - run postgresql, redis, prometheus, grafana container

### Local development usage:
    make local
    make run

### SWAGGER UI:
http://localhost:5000/swagger/index.html

### Jaeger UI:
http://localhost:16686

### Prometheus UI:
http://localhost:9090

### Grafana UI:
http://localhost:3000