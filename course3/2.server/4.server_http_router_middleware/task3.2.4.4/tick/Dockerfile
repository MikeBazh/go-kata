# Используем официальный образ Go как базовый
FROM golang:1.16-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Копируем папку graph внутрь контейнера
COPY graph /app/graph

# Скачиваем все зависимости
RUN go mod tidy

# Собираем приложение
RUN go build -o main

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/graph /app/graph
COPY --from=builder /app/main /main

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
