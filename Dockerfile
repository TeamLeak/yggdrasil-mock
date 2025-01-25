# Используем базовый образ Golang
FROM golang:1.20-alpine

# Устанавливаем зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o main .

# Указываем команду запуска
CMD ["./main"]
