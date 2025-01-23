#!/bin/bash

# Check if the release version is provided as an argument
if [ -z "$1" ]; then
    echo "Usage: $0 <release_version>"
    exit 1
fi

# Set the release version from the first argument
RELEASE_VERSION="$1"

# Create the output directory if it doesn't exist
mkdir -p ./bin

# Function to build for a specific OS and architecture
build() {
    local os=$1
    local arch=$2
    local output_name="geosite_editor_${RELEASE_VERSION}_${os}_${arch}"

    if [ "$os" == "windows" ]; then
        output_name="${output_name}.exe"
    fi

    # Set the environment variables for cross-compilation
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o "./bin/${output_name}" ./main.go
}

# Build for Linux
build linux amd64
build linux arm64

# Build for macOS
build darwin amd64
build darwin arm64

# Build for Windows
build windows amd64
build windows arm64

echo "Builds completed. Binaries are located in ./bin."
