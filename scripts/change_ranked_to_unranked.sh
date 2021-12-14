#!/bin/bash

RACEID=94416

if [ -z "$1" ]
  then
    echo "No argument supplied"
fi

source "$DIR/../.env"

mysql -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "UPDATE races SET ranked=0 WHERE id=$RACEID"
