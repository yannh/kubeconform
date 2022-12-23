#!/usr/bin/env sh

# This script was borrowed from https://github.com/quintush/helm-unittest

if [ -z "$HELM_PLUGIN_DIR" ]; then
    echo "No HELM_PLUGIN_DIR defined"

    exit 1
fi

PROJECT_NAME='kubeconform'
PROJECT_GH="yannh/$PROJECT_NAME"
PROJECT_CHECKSUM_FILE='CHECKSUMS'
HELM_PLUGIN_PATH="$HELM_PLUGIN_DIR"

# Convert the HELM_PLUGIN_PATH to unix if cygpath is
# available. This is the case when using MSYS2 or Cygwin
# on Windows where helm returns a Windows path but we
# need a Unix path
if type cygpath >/dev/null 2>&1; then
  echo 'Use Sygpath'

  HELM_PLUGIN_PATH=$(cygpath -u "$HELM_PLUGIN_PATH")
fi

if [ "$SKIP_BIN_INSTALL" = '1' ]; then
  echo 'Skipping binary install'

  exit
fi

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?

  if [ "$result" != '0' ]; then
    echo "Failed to install $PROJECT_NAME"
    echo 'For support, go to https://github.com/kubernetes/helm'
  fi

  exit $result
}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)

  case "$ARCH" in
    armv5*) ARCH='armv5';;
    armv6*) ARCH='armv6';;
    armv7*) ARCH='armv7';;
    aarch64) ARCH='arm64';;
    x86) ARCH='386';;
    x86_64) ARCH='amd64';;
    i686) ARCH='386';;
    i386) ARCH='386';;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname | tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Msys support
    msys*) OS='windows';;
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
    # MacOS
    darwin) OS='darwin';;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  supported='darwin-arm64,darwin-amd64,linux-amd64,linux-armv6,linux-arm64,linux-386,windows-arm64,windows-armv6,windows-amd64,windows-386'

  if ! echo "$supported" | grep -q "$OS-$ARCH"; then
    echo "No prebuild binary for $OS-$ARCH"

    exit 1
  fi

  if type curl >/dev/null 2>&1; then
    DOWNLOADER='curl'
  elif type wget >/dev/null 2>&1; then
    DOWNLOADER='wget'
  else
    echo 'Either curl or wget is required'

    exit 1
  fi

  echo "Support $OS-$ARCH"
}

# getDownloadURL checks the latest available version.
getDownloadURLs() {
  # Use the GitHub API to find the latest version for this project.
  latest_url="https://api.github.com/repos/$PROJECT_GH/releases/latest"

  if [ -z "$HELM_PLUGIN_UPDATE" ]; then
    version=$(git describe --tags --exact-match 2>/dev/null || true)

    if [ -n "$version" ]; then
      latest_url="https://api.github.com/repos/$PROJECT_GH/releases/tags/$version"
    fi
  fi

  echo "Retrieving $latest_url"

  if [ $DOWNLOADER = 'curl' ]; then
    DOWNLOAD_URL=$(curl -sL "$latest_url" | grep "$OS-$ARCH" | awk '/"browser_download_url":/{gsub(/[,"]/,"", $2); print $2}' 2>/dev/null)
    PROJECT_CHECKSUM=$(curl -sL "$latest_url" | grep "$PROJECT_CHECKSUM_FILE" | awk '/"browser_download_url":/{gsub(/[,"]/,"", $2); print $2}' 2>/dev/null)
  elif [ $DOWNLOADER = 'wget' ]; then
    DOWNLOAD_URL=$(wget -q -O - "$latest_url" | grep "$OS-$ARCH" | awk '/"browser_download_url":/{gsub(/[,"]/,"", $2); print $2}' 2>/dev/null)
    PROJECT_CHECKSUM=$(wget -q -O - "$latest_url" | grep "$PROJECT_CHECKSUM_FILE" | awk '/"browser_download_url":/{gsub(/[,"]/,"", $2); print $2}' 2>/dev/null)
  fi

  if [ -z "$DOWNLOAD_URL" ]; then
    echo 'Failed to get DOWNLOAD_URL'

    exit 1
  elif [ -z "$PROJECT_CHECKSUM" ]; then
    echo 'Failed to get PROJECT_CHECKSUM'

    exit 1
  fi
}

# downloadFiles downloads the latest binary package and also the checksum
# for that binary.
downloadFiles() {
  PLUGIN_TMP_FOLDER=$(mktemp -d)
  CHECKSUM_FILE_PATH="$PLUGIN_TMP_FOLDER/$PROJECT_CHECKSUM_FILE"

  echo "Downloading '$DOWNLOAD_URL' and '$PROJECT_CHECKSUM' to location $PLUGIN_TMP_FOLDER"

  if [ $DOWNLOADER = 'curl' ]; then
    (cd "$PLUGIN_TMP_FOLDER" && curl -sLO "$DOWNLOAD_URL")
    curl -s -L -o "$CHECKSUM_FILE_PATH" "$PROJECT_CHECKSUM"
  elif [ $DOWNLOADER = 'wget' ]; then
    wget -P "$PLUGIN_TMP_FOLDER" "$DOWNLOAD_URL"
    wget -q -O "$CHECKSUM_FILE_PATH" "$PROJECT_CHECKSUM"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  echo 'Verifying SHA for the file'

  DOWNLOAD_FILE=$(find "$PLUGIN_TMP_FOLDER" -name "*.tar.gz")

  if [ -z "$DOWNLOAD_FILE" ]; then
    DOWNLOAD_FILE=$(find "$PLUGIN_TMP_FOLDER" -name "*.zip")
  fi

  DOWNLOAD_FILE_NAME=$(basename "$DOWNLOAD_FILE")

  (
    cd "$PLUGIN_TMP_FOLDER"

    if type shasum >/dev/null 2>&1; then
      grep "$DOWNLOAD_FILE_NAME" "$CHECKSUM_FILE_PATH" | shasum -a 256 -c -s
    elif type sha256sum >/dev/null 2>&1; then
      if grep -q 'ID=alpine' /etc/os-release; then
        grep "$DOWNLOAD_FILE_NAME" "$CHECKSUM_FILE_PATH" | sha256sum -c -s
      else
        grep "$DOWNLOAD_FILE_NAME" "$CHECKSUM_FILE_PATH" | sha256sum -c --status
      fi
    else
      echo 'No Checksum as there is no shasum or sha256sum found'
    fi
  )

  HELM_TMP="$PLUGIN_TMP_FOLDER/$PROJECT_NAME"
  mkdir "$HELM_TMP"

  tar -C "$HELM_TMP" -xf "$DOWNLOAD_FILE"

  echo "Preparing to install into $HELM_PLUGIN_PATH"

  HELM_TMP_BIN="$HELM_TMP/$PROJECT_NAME"

  # Use * to also copy the file with the exe suffix on Windows
  cp "$HELM_TMP_BIN"* "$HELM_PLUGIN_PATH/bin/"

  rm -r "$HELM_TMP"
  rm -r "$PLUGIN_TMP_FOLDER"

  echo "$PROJECT_NAME installed into $HELM_PLUGIN_PATH/bin"
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  # To avoid to keep track of the Windows suffix,
  # call the plugin assuming it is in the PATH
  PATH="$HELM_PLUGIN_PATH/bin:$PATH"

  kubeconform -v
}

# Stop execution on any error
trap "fail_trap" EXIT
set -e

# Execution
initArch
initOS
verifySupported
getDownloadURLs
downloadFiles
installFile
testVersion
