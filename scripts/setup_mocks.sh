#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
mockgen -source=${SCRIPT_DIR}/../internal/input/input.go -destination=${SCRIPT_DIR}/../.mock/mock_internal/input/input.go
mockgen -source=${SCRIPT_DIR}/../internal/storage/storage.go -destination=${SCRIPT_DIR}/../.mock/mock_internal/storage/storage.go
mockgen -source=${SCRIPT_DIR}/../internal/worker.go -destination=${SCRIPT_DIR}/../.mock/mock_internal/worker.go