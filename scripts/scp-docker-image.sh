#!/bin/bash

host=$1
image=$2

cleanup() {
    if [ -n $rfile ]; then
        echo "[R] Removing remote file" && \
            ssh $host "rm $rfile"
    fi

    if [ -n $lfile ]; then
        echo "[L] Removing local file" && \
            rm $lfile
    fi
}

trap cleanup EXIT

rfile=$(ssh $host "mktemp") && \
    lfile=$(mktemp) && \
    echo "[R] Saving image to $rfile" && \
    ssh $host "docker save $image > $rfile" && \
    echo "[L] rsyncing image to $lfile" && \
    rsync -avzP $host:$rfile $lfile && \
    echo "[L] Loading image from $lfile" && \
    docker load < $lfile
