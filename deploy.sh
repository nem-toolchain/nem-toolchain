#!/bin/bash

echo "Deploying..."
(export GITHUB_TOKEN=5311dab717c11c936a21d759e9364a1b63849df4 && curl -sL https://git.io/goreleaser | bash)