FROM golang:1.20.6-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

EXPOSE 8000

CMD ["go", "run", "main.go"]