#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Variables
BASEDIR=$(dirname $(realpath "$0"))
echo "Base Dir: $BASEDIR"
TARGET="bon-voyage-cli"
SRC_DIR="$BASEDIR/src"
BUILD_DIR="$BASEDIR/build"

# Determine the operating system
OS="$(uname -s)"
case "$OS" in
    Linux*)     GOOS="linux";;
    Darwin*)    GOOS="darwin";;
    *)          echo "Unsupported OS: $OS"; exit 1;;
esac

# Clean previous builds
echo "Cleaning previous builds..."
rm -f $BUILD_DIR/$TARGET

# Build the CLI
echo "Building CLI ..."
cd $SRC_DIR

GIT_HASH=$(git rev-parse HEAD)
go build -o $BUILD_DIR/$TARGET -ldflags="-X main.Commit=$GIT_HASH"
cd $BASEDIR
echo "Build finished. See README.md for usage."
