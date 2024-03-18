# Используем официальный образ golang как базовый образ
FROM golang:latest AS builder

# Устанавливаем git, необходимый для загрузки зависимостей
RUN apt-get update && apt-get install -y git

# Копируем исходный код приложения
COPY ./ ./

# Если используется Go modules, копируем и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Собираем приложение
RUN go build -o app .

# Используем минимальный образ alpine в качестве базового образа для результата
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Указываем порт, который будет слушать приложение
EXPOSE 8080

# Команда для запуска приложения при запуске контейнера
CMD ["./app"]
