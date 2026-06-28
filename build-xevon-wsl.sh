#!/bin/sh
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/os/.bun/bin
export GOTMPDIR=/tmp
export GOCACHE=/home/os/.cache/go-build
export GOPATH=/home/os/go
export CGO_ENABLED=1

cd /mnt/c/Users/PC/Downloads/xevon-main

echo "[*] Disk space check..."
df -h /

echo "[*] Running go build with Linux tmp/cache dirs..."
VER=$(grep -E 'Version[[:space:]]+=' pkg/cli/version.go | cut -d '"' -f 2)
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo local)
BUILDTIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
PKG=github.com/xevonlive-dev/xevon/pkg/cli

echo "[*] Version=$VER Commit=$COMMIT BuildTime=$BUILDTIME"

go build \
  -ldflags "-s -w -X ${PKG}.Version=${VER} -X ${PKG}.Commit=${COMMIT} -X ${PKG}.BuildTime=${BUILDTIME}" \
  -o bin/xevon \
  ./cmd/xevon

STATUS=$?
if [ $STATUS -eq 0 ]; then
  echo ""
  echo "[✓] Build SUCCESS! Binary at bin/xevon"
  ls -lh bin/xevon
  echo ""
  echo "[*] Running: bin/xevon version"
  ./bin/xevon version
else
  echo "[!] Build FAILED with exit code $STATUS"
fi
