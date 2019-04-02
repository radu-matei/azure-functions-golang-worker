#! /bin/bash

# Directory name for start.sh
DIR="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"
echo "executing $DIR/worker $@"
$DIR/worker $@
