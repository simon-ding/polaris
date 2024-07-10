FROM ubuntu:20.04 as flutter_build

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update 
RUN apt-get install -y curl git wget unzip libgconf-2-4 gdb libstdc++6 libglu1-mesa fonts-droid-fallback python3
RUN apt-get clean


# download Flutter SDK from Flutter Github repo
RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter

# Set flutter environment path
ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:${PATH}"

# Run flutter doctor
RUN flutter doctor

# Enable flutter web
RUN flutter channel stable
RUN flutter upgrade
RUN flutter config --enable-web

# Copy files to container and build
RUN mkdir /app/
COPY . /app/
WORKDIR /app/
RUN flutter build web


# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.20 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY --from=flutter_build go.mod .
COPY --from=flutter_build go.sum .
RUN go mod download

COPY --from=flutter_build . .

# 指定OS等，并go build
RUN GOOS=linux GOARCH=amd64 go build -o polaris ./cmd/ 

FROM debian:12

WORKDIR /app
RUN apt-get update && apt-get -y install ca-certificates

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/polaris .

EXPOSE 8080

ENTRYPOINT ["./polaris"]