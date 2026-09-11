#!/bin/sh
case "$1" in
  *Username*) printf '%s' "${GIT_USERNAME:-oauth2}" ;;
  *Password*) printf '%s' "$GIT_TOKEN" ;;
esac
