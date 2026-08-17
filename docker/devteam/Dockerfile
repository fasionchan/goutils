ARG NODE_VERSION=24-bookworm-slim
ARG GOLANG_VERSION=1.26
ARG PYTHON_VERSION=3.12

# ARG PI_CODING_AGENT_VERSION=0.82.1
# ARG OPENCODE_VERSION=1.18.15
# ARG CODEX_VERSION=0.147.0

ARG DSH_VERSION=latest
ARG PI_CODING_AGENT_VERSION=latest
ARG OPENCODE_VERSION=latest
ARG CODEX_VERSION=latest

FROM golang:${GOLANG_VERSION} AS GoVendor

FROM node:${NODE_VERSION}

ARG PYTHON_VERSION
ARG DSH_VERSION
ARG PI_CODING_AGENT_VERSION
ARG OPENCODE_VERSION
ARG CODEX_VERSION

# 基础包
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

# Golang
COPY --from=GoVendor /usr/local/go/ /usr/local/go/

ENV GOROOT=/usr/local/go \
    PATH=/usr/local/go/bin:$PATH

# Python
ENV UV_INSTALL_DIR=/usr/local/bin
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

ENV UV_PYTHON_INSTALL_DIR=/opt/uv/python \
    UV_PYTHON_BIN_DIR=/usr/local/bin \
    PATH="/usr/local/bin:$PATH"

RUN mkdir -p /opt/uv/python && chmod 755 /opt/uv/python && \
    uv python install ${PYTHON_VERSION} && uv python pin --global ${PYTHON_VERSION}

# Agents
RUN npm install -g @deepseek-ai/dsh@${DSH_VERSION}
RUN npm install -g --ignore-scripts @earendil-works/pi-coding-agent@${PI_CODING_AGENT_VERSION}
RUN npm install -g opencode-ai@${OPENCODE_VERSION}
RUN npm install -g @openai/codex@${CODEX_VERSION}

ENV PATH="/root/.local/bin:$PATH"
RUN curl https://cursor.com/install -fsS | bash && cursor-agent sandbox disable

# Multica
RUN curl -fsSL https://raw.githubusercontent.com/multica-ai/multica/main/scripts/install.sh | bash

COPY docker/devteam/bin/multica.sh /usr/local/bin/multica.sh

# VOLUME /root/.pi
# VOLUME /home/node/.pi
# VOLUME /workspace

WORKDIR /workspace

ENV PI_SANDBOX=danger-full-access

ENTRYPOINT ["multica.sh"]
