#!/bin/bash

cfg=$(mktemp)
trap 'rm -f "$cfg"' EXIT

cat > "$cfg" <<EOF
log
proxy -a -p10080
EOF

podman run -d --replace --name 3proxy-lo \
	-v "$cfg:/etc/3proxy/3proxy.cfg:ro" \
	-p 10080:10080 \
	3proxy/3proxy:latest
