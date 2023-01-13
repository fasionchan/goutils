#!/bin/sh

# Author: fasion
# Created time: 2022-11-24 09:49:40
# Last Modified by: fasion
# Last Modified time: 2022-11-24 10:28:09

SELF_PATH=`realpath "$0"`
SCRIPT_DIR_PATH=`dirname "$SELF_PATH"`
REPO_PATH=`dirname "$SCRIPT_DIR_PATH"`

last_version=`git describe --tags --abbrev=0`
version=`echo "${last_version}" | awk -F '.' '{ print $1"."$2"."$3+1 }'`

(
	cd "${REPO_PATH}"

	git log --name-status HEAD^..HEAD | cat

	echo "Last version: ${last_version}"
	echo "New version: ${version}"

	read -p "Are you sure to publish this version (y or n)? " value && \
		[ "$value" == "y" ] && \
		git tag "$version" && \
		git push origin "$version" && \
		GOPROXY=proxy.golang.org go list -m github.com/fasionchan/goutils@$version
)


