FROM golang:1.22 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o polaris ./cmd/

FROM debian:12

WORKDIR /app
RUN apt-get update && apt-get -y install ca-certificates

COPY --from=builder /app/polaris .

EXPOSE 8080

ENTRYPOINT ["./polaris"]