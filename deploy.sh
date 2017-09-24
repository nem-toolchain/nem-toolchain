#!/bin/bash

echo "Deploying..."
(export GITHUB_TOKEN=d0898dbb529e387638730ccb78ac0b7f5be34c61 && curl -sL https://git.io/goreleaser | bash)