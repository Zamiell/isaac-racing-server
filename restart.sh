#!/bin/bash

# Get the directory of this script
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd "$DIR/src"
GOPATH="/root/go" "/usr/local/go/bin/go" install
if [ $? -eq 0 ]; then
        mv "/root/go/bin/src" "/root/go/bin/isaac-racing-server"
	supervisorctl restart isaac-racing-server
else
	echo "isaac-racing-server - Go compilation failed!"
fi
