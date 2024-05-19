FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o app .
ENTRYPOINT ["./app"]
