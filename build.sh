#! /usr/bin/env bash
set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

EXE_PATH="GCS.app/Contents/MacOS/gcs"
COPYRIGHT_YEARS=2019-$(date "+%Y")
GOLANGCI_LINT_VERSION=1.45.0

# Process args
for arg in "$@"; do
  case "$arg" in
  --all | -a)
    LINT=1
    TEST=1
    RACE=-race
    ;;
  --lint | -l) LINT=1 ;;
  --race | -r)
    TEST=1
    RACE=-race
    ;;
  --test | -t) TEST=1 ;;
  --help | -h)
    echo "$0 [options]"
    echo "  -a, --all  Equivalent to --lint --race"
    echo "  -l, --lint Run the linters"
    echo "  -r, --race Run the tests with race-checking enabled"
    echo "  -t, --test Run the tests"
    echo "  -h, --help This help text"
    exit 0
    ;;
  *)
    echo "Invalid argument: $arg"
    exit 1
    ;;
  esac
done

# Setup the tools we'll need
TOOLS_DIR=$PWD/tools
mkdir -p "$TOOLS_DIR"
if [ -z $SKIP_LINTERS ]; then
  if [ ! -e "$TOOLS_DIR/golangci-lint" ] || [ "$("$TOOLS_DIR/golangci-lint" version 2>&1 | awk '{ print $4 }' || true)x" != "${GOLANGCI_LINT_VERSION}x" ]; then
    echo -e "\033[33mInstalling version $GOLANGCI_LINT_VERSION of golangci-lint into $TOOLS_DIR...\033[0m"
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$TOOLS_DIR" v$GOLANGCI_LINT_VERSION
  fi
fi
export PATH=$TOOLS_DIR:$PATH

# Setup version info
if command -v git 2>&1 >/dev/null; then
  if [ -z "$(git status --porcelain)" ]; then
    STATE=clean
  else
    STATE=dirty
  fi
  GIT_VERSION=$(git rev-parse HEAD)-$STATE
  GIT_TAG=$(git tag --points-at HEAD)
  if [ -z "$GIT_TAG" ]; then
    GIT_TAG=$(git tag --list --sort -version:refname | head -1)
    if [ -n "$GIT_TAG" ]; then
      GIT_TAG=$GIT_TAG~
    fi
  fi
  if [ -n "$GIT_TAG" ]; then
    VERSION=$(echo "$GIT_TAG" | sed -E "s/^v//")
  else
    VERSION=""
  fi
else
  GIT_VERSION=Unknown
  VERSION=""
fi
BUILD_NUMBER=$(date -u "+%Y%m%d%H%M%S")

echo -e "\033[33mBuilding...\033[0m"

# Generate the source
go generate ./gen/enumgen.go

# Prepare the bundle
/bin/rm -rf "GCS.app"
mkdir -p "GCS.app/Contents/MacOS"
mkdir -p "GCS.app/Contents/Resources"
cp bundle/*.icns "GCS.app/Contents/Resources/"

# Build our code
LINK_FLAGS="-X 'github.com/richardwilkes/toolbox/cmdline.AppVersion=$VERSION'"
LINK_FLAGS="$LINK_FLAGS -X 'github.com/richardwilkes/toolbox/cmdline.BuildNumber=$BUILD_NUMBER'"
LINK_FLAGS="$LINK_FLAGS -X 'github.com/richardwilkes/toolbox/cmdline.GitVersion=$GIT_VERSION'"
go build -o "$EXE_PATH" -v -ldflags=all="$LINK_FLAGS" .
sed -e "s/SHORT_APP_VERSION/$($EXE_PATH -v)/" -e "s/LONG_APP_VERSION/$($EXE_PATH -V)/" -e "s/COPYRIGHT_YEARS/$COPYRIGHT_YEARS/" bundle/Info.plist >"GCS.app/Contents/Info.plist"
touch "GCS.app" # Allows the Finder to notice changes

# Run the tests
if [ "$TEST"x == "1x" ]; then
  if [ -n "$RACE" ]; then
    echo -e "\033[32mTesting with -race enabled...\033[0m"
  else
    echo -e "\033[32mTesting...\033[0m"
  fi
  go test $RACE ./...
fi

# Run the linters
if [ "$LINT"x == "1x" ]; then
  echo -e "\033[32mLinting...\033[0m"
  golangci-lint run
fi
