#!/bin/bash

go install
supervisorctl restart isaac-racing-server
