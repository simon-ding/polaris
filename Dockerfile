FROM instrumentisto/flutter:3 AS flutter
WORKDIR /app
COPY ./ui/pubspec.yaml ./ui/pubspec.lock ./
RUN flutter pub get
COPY ./ui/ ./
RUN flutter build web --no-web-resources-cdn --web-renderer html

# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.23 as builder

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
ENV TZ="Asia/Shanghai" GIN_MODE=release PUID=0 PGID=0 UMASK=0 PERMS=true

WORKDIR /app
RUN apt-get update && apt-get -y install ca-certificates tzdata gosu tini && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \ 
    && echo "${TZ}" > /etc/timezone && apt-get clean

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/polaris .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x /app/entrypoint.sh

VOLUME /app/data
EXPOSE 8080

ENTRYPOINT ["tini","./entrypoint.sh"]