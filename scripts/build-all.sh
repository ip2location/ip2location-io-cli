#!/bin/bash

# Build binary for all architectures

set -e

DIR=`dirname $0`
ROOT=$DIR/..

NAME=$1
VERSION=$2

for t in                                                                      \
    darwin_amd64                                                              \
    darwin_arm64                                                              \
    dragonfly_amd64                                                           \
    freebsd_386                                                               \
    freebsd_amd64                                                             \
    freebsd_arm                                                               \
    freebsd_arm64                                                             \
    linux_386                                                                 \
    linux_amd64                                                               \
    linux_arm                                                                 \
    linux_arm64                                                               \
    netbsd_386                                                                \
    netbsd_amd64                                                              \
    netbsd_arm                                                                \
    netbsd_arm64                                                              \
    openbsd_386                                                               \
    openbsd_amd64                                                             \
    openbsd_arm                                                               \
    openbsd_arm64                                                             \
    solaris_amd64                                                             \
    windows_386                                                               \
    windows_amd64                                                             \
    windows_arm ;
do
    os="${t%_*}"
    arch="${t#*_}"
    output="${NAME}_${VERSION}_${os}_${arch}"

    if [ "$os" == "windows" ] ; then
        output+=".exe"
    fi

    echo "Building ${output}"
    GOOS=$os GOARCH=$arch go build                                            \
        -o $ROOT/build/${output}                                              \
        $ROOT/${NAME}
done

wait
