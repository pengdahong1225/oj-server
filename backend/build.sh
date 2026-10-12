#!/bin/sh

set -e

build() {
  name=$1
  path=$2
  os=$3
  arch=$4
  out=${path}/bin

  mkdir -p "$out"

  GOOS=$os GOARCH=$arch go build -trimpath -o "$out/${name}" "$path"

  echo "✔ build $name success for $os/$arch"
}

# gateway
build gateway.exe ./svr/gateway windows amd64
build gateway ./svr/gateway linux amd64

# user
build user.exe ./svr/user windows amd64
build user ./svr/user linux amd64

# problem
build problem.exe ./svr/problem windows amd64
build problem ./svr/problem linux amd64

# judge
build judge.exe ./svr/judge windows amd64
build judge ./svr/judge linux amd64