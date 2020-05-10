FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /src/main .

CMD ["./main"]
