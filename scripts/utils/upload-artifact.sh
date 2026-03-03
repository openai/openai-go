#!/usr/bin/env bash
set -exuo pipefail

DIST_DIR="dist"
FILENAME="source.zip"

mapfile -d '' files < <(
  find . -type f \
    \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) \
    ! -path "./${DIST_DIR}/*" \
    -print0
)

if [[ ${#files[@]} -eq 0 ]]; then
  echo -e "\033[31mNo Go source files found for packaging.\033[0m"
  exit 1
fi

mkdir -p "$DIST_DIR"
rm -f "${DIST_DIR}/${FILENAME}"

relative_files=()
for file in "${files[@]}"; do
  relative_files+=("${file#./}")
done

zip "${DIST_DIR}/${FILENAME}" "${relative_files[@]}"

RESPONSE=$(curl -X POST "$URL?filename=$FILENAME" \
  -H "Authorization: Bearer $AUTH" \
  -H "Content-Type: application/json")

SIGNED_URL=$(echo "$RESPONSE" | jq -r '.url')

if [[ "$SIGNED_URL" == "null" ]]; then
  echo -e "\033[31mFailed to get signed URL.\033[0m"
  exit 1
fi

UPLOAD_RESPONSE=$(curl -v -X PUT \
  -H "Content-Type: application/zip" \
  --data-binary "@${DIST_DIR}/${FILENAME}" "$SIGNED_URL" 2>&1)

if echo "$UPLOAD_RESPONSE" | grep -q "HTTP/[0-9.]* 200"; then
  echo -e "\033[32mUploaded build to Stainless storage.\033[0m"
  echo -e "\033[32mInstallation: Download and unzip: 'https://pkg.stainless.com/s/openai-go/$SHA'. Run 'go mod edit -replace github.com/openai/openai-go/v3=/path/to/unzipped_directory'.\033[0m"
else
  echo -e "\033[31mFailed to upload artifact.\033[0m"
  exit 1
fi
