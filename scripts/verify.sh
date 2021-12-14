#!/bin/bash

# Get the directory of this script
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source "$DIR/../.env"

if [ -z "$1" ]; then
  echo "Supply the username of the person to verify as an argument."
  exit 1
fi

mysql -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "UPDATE users SET verified = 1 WHERE username = '$1';"
