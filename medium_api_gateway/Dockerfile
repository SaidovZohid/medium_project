FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /medium

COPY . .


RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /medium
RUN mkdir media

COPY --from=builder /medium/main .

EXPOSE 8080

CMD [ "/medium/main" ]