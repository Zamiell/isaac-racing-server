#!/bin/bash

go install
if [ $? -eq 0 ]; then
	supervisorctl restart isaac-racing-server
fi
