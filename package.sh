#!/bin/sh
TOPDIR=$(dirname "$0")
TOPDIR=$(realpath "$TOPDIR")

if [ -n "$1" ]
then
    export GOARCH="$1"
fi

clean_build_tmp() {
    if [ -n "$BUILD_TEMP_DIR" ] && [ -d "$BUILD_TEMP_DIR" ]; then
      log_info "clean build tmp: $BUILD_TEMP_DIR"
      rm -r "$BUILD_TEMP_DIR"
    fi
}

log_fault() {
    clean_build_tmp
    echo "[FAULT] $1"
    exit 1
}

log_info() {
    echo "[INFO] $1"
}

log_warn() {
    echo "[WARN] $1"
}

log_info "building web/frontend ..."

if ! cd "$TOPDIR/web/frontend"; then
    log_fault "Failed to cd to web/frontend directory"
fi

if ! pnpm install
then
    log_fault "Failed to install pnpm dependencies"
fi

if [ -d ../dist ]
then
  if ! rm -rf ../dist; then
    log_fault "Failed to remove old web dist directory"
  fi
fi

if ! pnpm run build; then
    log_fault "Failed to build web/frontend"
fi

if ! cd "$TOPDIR"; then
    log_fault "Failed to cd to $TOPDIR"
fi

BUILD_TEMP_DIR="/tmp/picp-build"
if [ -d "$BUILD_TEMP_DIR" ];
then
  if ! rm -rf "$BUILD_TEMP_DIR"; then
    log_fault "Failed to remove old build temp directory"
  fi
fi

log_info "copy to build temp directory $BUILD_TEMP_DIR"
if ! cp -r install "$BUILD_TEMP_DIR"; then
    log_fault "Failed to copy install directory to $BUILD_TEMP_DIR"
fi
log_info "building picp ..."
BUILD_BIN_DIR="$BUILD_TEMP_DIR/opt/picp"
if ! go mod tidy; then
    log_fault "Failed to go mod tidy"
fi

if ! go build -o "$BUILD_BIN_DIR/picp" -ldflags "-X 'main.GoVersion=$(go version)' -X 'main.BuildTime=$(date "+%F %T")'" .; then
    log_fault "Failed to build picp"
fi

if ! cp picp.ini $BUILD_TEMP_DIR/opt/picp; then
    log_fault "Failed to copy picp.ini to $BUILD_TEMP_DIR/opt/picp"
fi

if ! CURRENT_ARCH=$(go env GOARCH); then
  log_fault "Failed to get current architecture"
fi

if [ "$CURRENT_ARCH" = "arm" ]; then
  CURRENT_ARCH="armhf"
fi

if ! VERSION=$(cat VERSION); then
  log_fault "Failed to get version"
fi

if ! SIZE_RET=$(du -sb "$BUILD_TEMP_DIR/opt" "$BUILD_TEMP_DIR/etc" | cut -f1); then
  log_fault "Failed to get installed size"
fi

TOTAL_SIZE=0
for size in $(echo "$SIZE_RET"); do
  TOTAL_SIZE=$((TOTAL_SIZE + size))
done
TOTAL_SIZE=$((TOTAL_SIZE / 1024))
if ! sed -i "s/Installed-Size:/Installed-Size: $TOTAL_SIZE/g" "$BUILD_TEMP_DIR/DEBIAN/control"; then
  log_fault "Failed to update control file Installed-Size"
fi

if ! sed -i "s/Version:/Version: $VERSION/g" "$BUILD_TEMP_DIR/DEBIAN/control"; then
  log_fault "Failed to update control file Version"
fi

if ! sed -i "s/Architecture:/Architecture: $CURRENT_ARCH/g" "$BUILD_TEMP_DIR/DEBIAN/control"; then
  log_fault "Failed to update control file Architecture"
fi

if dpkg-deb -b "$BUILD_TEMP_DIR" "picp-$VERSION-$CURRENT_ARCH.deb"
then
    clean_build_tmp
    log_info "build success"
else
    log_fault "Failed to build deb package"
fi