ARG NODE_VERSION=24-bookworm-slim
ARG PI_CODING_AGENT_VERSION=0.82.1
ARG OPENCODE_VERSION=1.18.15
ARG CODEX_VERSION=0.147.0

FROM node:${NODE_VERSION}

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    fd-find \
    git \
    ripgrep \
    openssh-client \
    && rm -rf /var/lib/apt/lists/*

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
