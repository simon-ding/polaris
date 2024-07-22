FROM golang:1.22-alpine3.20 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o polaris ./cmd/

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache bash ca-certificates

COPY --from=builder /app/polaris .

EXPOSE 8080

ENTRYPOINT ["./polaris"]