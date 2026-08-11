ARG NODE_VERSION=24-bookworm-slim
ARG PI_CODING_AGENT_VERSION=0.82.1
ARG OPENCODE_VERSION=1.18.15
ARG CODEX_VERSION=0.147.0
ARG GOLANG_VERSION=1.26
ARG PYTHON_VERSION=3.14.7

From golang:${GOLANG_VERSION} as GoVendor

FROM node:${NODE_VERSION}

ARG PYTHON_VERSION

COPY --from=GoVendor /usr/local/go/ /usr/local/go/

ENV GOROOT=/usr/local/go \
    PYENV_ROOT=/opt/pyenv \
    PATH=/opt/pyenv/shims:/opt/pyenv/bin:/usr/local/go/bin:$PATH

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    fd-find \
    git \
    ripgrep \
    openssh-client \
    # pyenv / CPython build dependencies
    build-essential \
    libssl-dev \
    zlib1g-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    libncursesw5-dev \
    xz-utils \
    tk-dev \
    libxml2-dev \
    libxmlsec1-dev \
    libffi-dev \
    liblzma-dev \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://pyenv.run | bash \
    && pyenv install "${PYTHON_VERSION}" \
    && pyenv global "${PYTHON_VERSION}" \
    && pyenv rehash \
    && printf '%s\n' \
        'export PYENV_ROOT=/opt/pyenv' \
        'export PATH="$PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH"' \
        'eval "$(pyenv init -)"' \
        > /etc/profile.d/pyenv.sh \
    && python --version \
    && pip --version

RUN npm install -g --ignore-scripts @earendil-works/pi-coding-agent@${PI_CODING_AGENT_VERSION}
RUN npm install -g opencode-ai@${OPENCODE_VERSION}
RUN npm install -g @openai/codex@${CODEX_VERSION}
RUN curl -fsSL https://raw.githubusercontent.com/multica-ai/multica/main/scripts/install.sh | bash

COPY docker/devteam/bin/multica.sh /usr/local/bin/multica.sh

VOLUME /root/.pi
VOLUME /home/node/.pi
VOLUME /workspace

WORKDIR /workspace

ENV PI_SANDBOX=danger-full-access

ENTRYPOINT ["multica.sh"]
