#!/bin/sh
set -e

go build -o vfs-shell ./

./vfs-shell --vfs=tests/stage3/vfs_multiple.xml --script=tests/stage5/script_rm_chown.txt &
sleep 20 && kill $! 2>/dev/null