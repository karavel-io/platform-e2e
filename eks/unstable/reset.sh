#!/usr/bin/env sh

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

rm -rf "$SCRIPT_DIR/applications"
rm -rf "$SCRIPT_DIR/projects"
rm -rf "$SCRIPT_DIR/kustomization.yml"
rm -rf "$SCRIPT_DIR/vendor"