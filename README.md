# Logger Service
Logger Service is a lightweight and extensible logging service built with Go. It enables structured logging and publishes logs asynchronously to Redis for further processing or storage.

## Features
- Supports multiple log levels: INFO, WARN, and ERROR.
- Asynchronous log publishing to Redis for scalability.
- JSON serialization for structured log storage.
- Dependency injection for modularity and testability.

## Usage
### Prerequisites
- Docker and Docker Compose installed on your system.
- Redis is used as the message broker.

### Running the Service
1. Clone the repository:
```bash
git clone https://github.com/yourusername/logger-service.git
cd logger-service
```
2. Build and run the service using Docker Compose:
```bash
docker-compose up
```
3. The logger service will be accessible on http://localhost:7979.

## Example Logging API
- Log an Info Message:
```go
service.LogInfo("MyService", "This is an info message.")
```
- Log a Warning Message:
```go
service.LogWarning("MyService", "This is a warning message.")
```
- Log an Error Message:
```go
service.LogError("MyService", "This is an error message.")
```
- Log a Custom Message:
```go
service.LogMessage("MyService", "Custom log message.", "DEBUG")
```

## Environment Variables
- REDIS_HOST: Hostname of the Redis server.
- REDIS_PORT: Port of the Redis server.

These are configured in docker-compose.yml.

## Docker Compose Configuration
- Redis Service:
    * Image: redis:latest
    * Ports: 6379:6379
    * Persistent storage via redis_data.
- Logger Service:
    * Built from the Dockerfile.
    * Ports: 7979:7979.
    * Depends on the Redis service.

Acknowledgments

Special thanks to [Albert Karapetyan](https://github.com/AlbertKarapetyan) for leading the development of this service.