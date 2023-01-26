FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /medium

COPY . .

RUN apk add curl
RUN go build -o main cmd/main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.16

WORKDIR /medium

COPY --from=builder /medium/main .
COPY --from=builder /medium/migrate ./migrate
COPY migrations ./migrations

CMD [ "/medium/main" ]