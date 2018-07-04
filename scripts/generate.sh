#!/usr/bin/env bash

ROOT_PACKAGE="github.com/bitshifta/crd-controller"
CUSTOM_RESOURCE_NAME="podcounter"
CUSTOM_RESOURCE_VERSION="v1"

go get -u k8s.io/code-generator/...

cd $GOPATH/src/k8s.io/code-generator

./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
