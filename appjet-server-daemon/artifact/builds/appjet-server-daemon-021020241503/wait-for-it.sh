#!/bin/bash
# wait-for-it.sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "1 minute before launching the app"
sleep 60

exec $cmd
