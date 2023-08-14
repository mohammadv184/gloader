#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
  echo "Generating $sh completion..."
  mkdir "completions/$sh"
	go run ./cmd/gloader completion "$sh" >"completions/$sh/gloader"
done