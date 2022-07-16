#!/bin/bash

VERSION="1.0"
TIMESTAMP="$(date)"

export GOOS GOARCH

for GOOS in darwin linux windows
do
	for GOARCH in amd64 arm64
	do
		echo "Building for ${GOOS}/${GOARCH}..."
		go build \
			-ldflags "-X main.version=${VERSION} -X 'main.timestamp=${TIMESTAMP}'" \
			-o gorwd-${GOOS}-${GOARCH} \
			main.go
	done
done
	
