# Используем официальный образ Go
FROM golang:1.23 AS runtime

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose the application port
EXPOSE 7979

# Default command
CMD go run ./main.go
