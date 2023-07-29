#!/bin/sh

# Adding autocompletion
# Automatically detect and add autocompletion for bash, zsh, fish and powershell

if [ -n "$BASH_VERSION" ]; then
  echo "Bash detected"
  echo "Adding autocompletion for bash"
  gloader completion bash > /etc/bash_completion.d/gloader
fi

if [ -n "$ZSH_VERSION" ]; then
  echo "Zsh detected"
  echo "Adding autocompletion for zsh"
  gloader completion zsh > "${fpath[1]}/_gloader"
fi

if [ -n "$FISH_VERSION" ]; then
  echo "Fish detected"
  echo "Adding autocompletion for fish"
  gloader completion fish > ~/.config/fish/completions/gloader.fish
fi