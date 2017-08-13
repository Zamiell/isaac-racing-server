#!/bin/bash

cd src
go install
mv "$GOPATH/bin/src" "$GOPATH/bin/isaac-racing-server"
if [ $? -eq 0 ]; then
	supervisorctl restart isaac-racing-server
fi
