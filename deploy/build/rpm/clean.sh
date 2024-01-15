#!/bin/bash 

set -Eeuo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1

help(){
    echo "clean.sh - rpm clean up action"
    echo "  --deep removes old container"
}

[ -d "build" ] && sudo rm -rf build/

[ $# -ge 1 ] && [ "$1" = "--deep" ] && docker compose rm -fsv
[ $# -ge 1 ] && [ "$1" = "help" ] && help

exit 0
