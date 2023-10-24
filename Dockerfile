FROM golang:1.19 as server_image

WORKDIR /build

COPY . .

# 执行指令
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && export PATH=$PATH:/go/bin \
    && go build -o xiaomi_camera_merge main.go



# 使用官方的Ubuntu 20.04作为基础镜像
FROM ubuntu:20.04

# 设置环境变量，防止交互式安装
ENV DEBIAN_FRONTEND=noninteractive

# 更新包列表并安装依赖
RUN apt-get update && apt-get install -y \
    software-properties-common \
    build-essential \
    wget \
    yasm \
    unzip \
    # 安装FFmpeg和依赖
    && apt-get install -y ffmpeg \
    # 清理不再需要的包
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/*

# 设置工作目录，可以在其中运行FFmpeg命令
WORKDIR /app

COPY --from=server_image /build/xiaomi_camera_merge /app/xiaomi_camera_merge


CMD ["./xiaomi_camera_merge","-path","./video"]
