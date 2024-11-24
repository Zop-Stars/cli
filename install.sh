#!/bin/bash

# Variables
GITHUB_USER="Zop-Stars"           # Replace with your GitHub username
GITHUB_REPO="cli"               # Replace with your GitHub repository name
BRANCH="cli-binary"                  # Replace with the branch name where the binary is stored
BINARY_PATH="gofr-gennie"     # Replace with the path to the binary file in the repo
DESTINATION="/usr/local/bin/gofr-gennie"

# Check if running as root or with sudo
if [[ "$EUID" -ne 0 ]]; then
  echo "Please run this script as root or with sudo."
  exit 1
fi

# Download the binary from the GitHub repo
echo "Downloading $BINARY_PATH from GitHub..."
curl -L -o "$BINARY_PATH" "https://raw.githubusercontent.com/$GITHUB_USER/$GITHUB_REPO/$BRANCH/$BINARY_PATH"
if [[ $? -ne 0 ]]; then
  echo "Failed to download $BINARY_PATH. Please check the URL and try again."
  exit 1
fi

# Make the binary executable
chmod +x "$BINARY_PATH"

# Move the binary to /usr/local/bin
echo "Moving $BINARY_PATH to $DESTINATION..."
mv "$BINARY_PATH" "$DESTINATION"
if [[ $? -eq 0 ]]; then
  echo "$BINARY_PATH successfully installed to $DESTINATION."
else
  echo "Failed to move $BINARY_PATH to $DESTINATION. Please check permissions and try again."
  exit 1
fi