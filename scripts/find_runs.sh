#!/bin/bash

USER="Zamiel"
BUILD=8

if [ -z "$1" ]
  then
    echo "No argument supplied"
fi

source "$DIR/../.env"

mysql -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" <<EOF
SELECT races.id, starting_build, race_participants.run_time / 1000
FROM races
JOIN race_participants ON races.id = race_participants.race_id
JOIN users ON race_participants.user_id = users.id
WHERE users.username = "$USER"
/*AND races.starting_build = $BUILD*/
AND races.ranked = 1
AND races.solo = 1
EOF
