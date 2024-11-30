
FROM golang:1.21-alpine as builder


ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o myapp .


FROM alpine:3.18


RUN apk --no-cache add ca-certificates

WORKDIR /root/


RUN mkdir -p /root/audits


RUN chmod -R 755 /root/audits


COPY --from=builder /app/myapp .


COPY config.json ./config.json


COPY --from=builder /app/UI /root/UI


EXPOSE 8080


CMD ["./myapp"]
