FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /medium

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /medium

COPY --from=builder /medium/main .
COPY templates ./templates

CMD [ "/medium/main" ]