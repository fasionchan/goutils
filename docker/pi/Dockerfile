ARG PI_CODING_AGENT_VERSION=0.82.1

FROM node:24-bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    fd-find \
    git \
    ripgrep \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g --ignore-scripts @earendil-works/pi-coding-agent@${PI_CODING_AGENT_VERSION}

VOLUME /root/.pi
VOLUME /home/node/.pi
VOLUME /workspace

WORKDIR /workspace

ENV PI_SANDBOX=danger-full-access

ENTRYPOINT ["pi"]
