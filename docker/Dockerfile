FROM golang:1.24-alpine as builder

WORKDIR /app

COPY ./go.mod ./go.sum


COPY . .

RUN go build -o /app/bin/main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main .

EXPOSE 1235

CMD ["./main"]