FROM instrumentisto/flutter:3 AS flutter
WORKDIR /app
COPY ./ui/pubspec.yaml ./ui/pubspec.lock ./
RUN flutter pub get
COPY ./ui/ ./
RUN flutter build web --no-web-resources-cdn

# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.22 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY  go.mod .
COPY  go.sum .
RUN go mod download

COPY . .

COPY --from=flutter /app/build/web ./ui/build/web/
# 指定OS等，并go build
RUN CGO_ENABLED=1 go build -o polaris -ldflags="-X polaris/db.Version=$(git describe --tags --long)"  ./cmd/ 

FROM debian:stable-slim
ENV TZ="Asia/Shanghai" GIN_MODE=release

WORKDIR /app
RUN apt-get update && apt-get -y install ca-certificates

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/polaris .

EXPOSE 8080

USER 1000:1000

ENTRYPOINT ["./polaris"]