#!/usr/bin/env bash

set -euo pipefail

current_goreleaser_version=$(go mod graph | grep 'github.com/goreleaser/goreleaser' |tr ' ' '@' | awk -F@ '{ print $2 }' | head -n 2 | tail -n 1)
echo "current goreleaser version from go.mod: ${current_goreleaser_version}"

echo "removing old fork..."
rm -rf ./goreleaser

goreleaser_version="${current_goreleaser_version:1}"
goreleaser_zip="v${goreleaser_version}.zip"

echo "downloading goreleaser@${goreleaser_version}"
wget "https://github.com/goreleaser/goreleaser/archive/refs/tags/${goreleaser_zip}"
unzip "${goreleaser_zip}"

mv "goreleaser-${goreleaser_version}" goreleaser

echo "renaming internal to int..."
mv goreleaser/internal goreleaser/int
find ./goreleaser -type f -name "*.go" -exec sed -i '' "s/github\.com\/goreleaser\/goreleaser\/internal/github\.com\/goreleaser\/goreleaser\/int/g" {} \;
echo "!/int/pipe/dist" >> goreleaser/.gitignore

echo "cleanup..."
rm "${goreleaser_zip}"
