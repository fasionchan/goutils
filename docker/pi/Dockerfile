ARG PI_CODING_AGENT_VERSION=v0.82.1

FROM node:24-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g @earendil-works/pi-coding-agent:${PI_CODING_AGENT_VERSION}

RUN useradd -m -u 1001 -g 1001 -s /bin/bash -d /home/pi pi

VOLUME /home/pi/.pi
VOLUME /workspace

WORKDIR /workspace

ENV PI_SANDBOX=danger-full-access

ENTRYPOINT ["pi"]

CMD ["--help"]