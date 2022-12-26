#!/bin/bash
FILES=$(git diff --cached --name-only --diff-filter=ACMR)

golangci-lint run --new --fix

git add $FILES
