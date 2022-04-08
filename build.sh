#! /usr/bin/env bash
set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

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

echo -e "\033[33mBuilding...\033[0m"

# Generate the source
go generate ./gen/enumgen.go

# Build our code
EXE="$(go env GOPATH)/bin/gcs"
case $(uname -s) in
Darwin*)
  if [ $(uname -p) == "arm" ]; then
    DEPLOYMENT_TARGET=11
  else
    DEPLOYMENT_TARGET=10.14
  fi
  MACOSX_DEPLOYMENT_TARGET=$DEPLOYMENT_TARGET go install -v .
  /bin/rm -rf GCS.app
  CONTENTS="GCS.app/Contents"
  mkdir -p "$CONTENTS/MacOS"
  mkdir -p "$CONTENTS/Resources"
  cp bundle/*.icns "$CONTENTS/Resources/"
  sed -e "s/SHORT_APP_VERSION/$($EXE -v | tr -d "\n")/" \
    -e "s/LONG_APP_VERSION/$($EXE -V |tr -d "\n")/" \
    -e "s/COPYRIGHT_YEARS/$($EXE --copyright-date | tr -d "\n")/" \
    bundle/Info.plist >"$CONTENTS/Info.plist"
  mv "$EXE" "$CONTENTS/MacOS/"
  ;;
Linux*)
  go install -v .
  /bin/rm -f gcs
  mv "$EXE" ./gcs
  ;;
MINGW*)
  go install -v -ldflags all="-H windowsgui" .
  /bin/rm -f gcs.exe
  mv "$EXE.exe" ./gcs.exe
  ;;
*)
  echo "Unsupported OS"
  false
  ;;
esac

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
  GOLANGCI_LINT_VERSION=1.45.2
  TOOLS_DIR=$PWD/tools
  mkdir -p "$TOOLS_DIR"
  if [ ! -e "$TOOLS_DIR/golangci-lint" ] || [ "$("$TOOLS_DIR/golangci-lint" version 2>&1 | awk '{ print $4 }' || true)x" != "${GOLANGCI_LINT_VERSION}x" ]; then
    echo -e "\033[33mInstalling version $GOLANGCI_LINT_VERSION of golangci-lint into $TOOLS_DIR...\033[0m"
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$TOOLS_DIR" v$GOLANGCI_LINT_VERSION
  fi
  echo -e "\033[32mLinting...\033[0m"
  $TOOLS_DIR/golangci-lint run
fi
