# dsh 运行环境镜像
#
# 以 node 为基础镜像，提供 golang + python 开发环境，
# 并预装 DeepSeek Harness CLI（dsh）。
#
# 用法:
#   docker build -f docker/dsh/Dockerfile -t dsh:latest .
#
#   # 启动 Web UI (默认端口 3080)
#   docker run -it --rm -p 3080:3080 dsh:latest
#
#   # 运行 headless 任务 (需注入 DEEPSEEK_API_KEY)
#   docker run -it --rm -e DEEPSEEK_API_KEY=xxx dsh:latest \
#       --profile headless "your task"

ARG NODE_VERSION=24-bookworm-slim
ARG GOLANG_VERSION=1.26
ARG PYTHON_VERSION=3.12
ARG DSH_VERSION=latest

FROM golang:${GOLANG_VERSION} AS GoVendor

FROM node:${NODE_VERSION}

ARG PYTHON_VERSION
ARG DSH_VERSION

# 基础工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    fd-find \
    git \
    make \
    openssh-client \
    ripgrep && \
    rm -rf /var/lib/apt/lists/*

# Golang 开发环境
COPY --from=GoVendor /usr/local/go/ /usr/local/go/

ENV GOROOT=/usr/local/go \
    GOPATH=/go \
    PATH="/usr/local/go/bin:/go/bin:$PATH"

# Python 开发环境 (uv 管理)
ENV UV_INSTALL_DIR=/usr/local/bin
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

ENV UV_PYTHON_INSTALL_DIR=/opt/uv/python \
    UV_PYTHON_BIN_DIR=/usr/local/bin \
    PATH="/usr/local/bin:$PATH"

RUN mkdir -p /opt/uv/python && chmod 755 /opt/uv/python && \
    uv python install ${PYTHON_VERSION} && uv python pin --global ${PYTHON_VERSION} && \
    ln -sf "$(uv python find ${PYTHON_VERSION})" /usr/local/bin/python3 && \
    ln -sf python3 /usr/local/bin/python

# dsh CLI
RUN npm install -g --ignore-scripts @deepseek-ai/dsh@${DSH_VERSION}

# pnpm: dsh plugin 命令会转发给 pnpm 管理 profile 插件
RUN npm install -g pnpm

# dsh 数据目录 (profiles / 会话) 与工作区
VOLUME /root/.dsh
VOLUME /workspace

WORKDIR /workspace

# dsh web 默认端口
EXPOSE 3080

# 默认启动 Web UI; 可覆盖为 headless 等其他入口
CMD ["dsh", "web"]
