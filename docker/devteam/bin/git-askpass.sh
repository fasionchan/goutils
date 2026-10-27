#!/bin/sh

case "$1" in
  *Username*) printf '%s' "${GIT_USERNAME:-oauth2}" ;;
  *Password*) printf '%s' "$GIT_TOKEN" ;;
esac

# example:
#
# export GIT_ASKPASS=git-askpass.sh  
# export GIT_USERNAME=your_username
# export GIT_TOKEN=your_token
# git clone https://devops.gf.com.cn/gf/ums/docs.git