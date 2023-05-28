# Встановлення базового образу Go
FROM golang:1.20.4

# Встановлення директорії робочої зони
WORKDIR /gses2.app

# Копіювання файлів go.mod та go.sum
COPY go.mod go.sum ./

# Виконання завантаження залежностей
RUN go mod download

# Копіювання решти файлів проекту до контейнера
COPY . .

# Встановлення необхідних залежностей
RUN go get github.com/go-mail/gomail

# Збірка Go-проекту
RUN go build -o main .

# Встановлення вхідного точки для контейнера
CMD ["./main"]