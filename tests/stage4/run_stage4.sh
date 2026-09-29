#!/bin/sh
set -e

go build -o vfs-shell ./

./vfs-shell --vfs=tests/stage3/vfs_nested.xml --script=tests/stage4/script_full.txt &
sleep 20 && kill $! 2>/dev/null