#!/bin/bash

set -e

BINARY_NAME="opstl"
INSTALL_DIR="/usr/local/bin"

echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" .

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"

echo "Done! You can now use '$BINARY_NAME' from anywhere."
