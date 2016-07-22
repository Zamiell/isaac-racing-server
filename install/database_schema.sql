/* sqlite3 database.sqlite < install/database_schema.sql */

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id                INTEGER    PRIMARY KEY  AUTOINCREMENT,
    auth0_id          TEXT       NOT NULL,
    username          TEXT       NOT NULL,
    datetime_created  INTEGER    DEFAULT (strftime('%s', 'now')),
    last_login        INTEGER    DEFAULT (strftime('%s', 'now')),
    last_ip           TEXT       NOT NULL,
    admin             INTEGER    DEFAULT 0
);

DROP TABLE IF EXISTS races;
CREATE TABLE races (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    name                  TEXT                  DEFAULT "-",
    status                TEXT                  DEFAULT "open", /* starting, in progress, finished */
    ruleset               TEXT                  DEFAULT "unseeded", /* seeded, diversity */
    seed                  TEXT                  DEFAULT "-",
    datetime_created      INTEGER               DEFAULT (strftime('%s', 'now')),
    datetime_started      INTEGER               DEFAULT 0,
    datetime_finished     INTEGER               DEFAULT 0,
    captain               INTEGER               NOT NULL,
    FOREIGN KEY(captain)  REFERENCES users(id)
);
CREATE INDEX races_index_status ON races (status);

DROP TABLE IF EXISTS race_participants;
CREATE TABLE race_participants (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id               INTEGER               NOT NULL,
    race_id               INTEGER               NOT NULL,
    status                TEXT                  DEFAULT "not ready", /* ready, racing, finished, quit */
    datetime_joined       INTEGER               DEFAULT (strftime('%s', 'now')),
    datetime_finished     INTEGER               DEFAULT 0,
    place                 INTEGER               DEFAULT 0,
    comment               TEXT                  DEFAULT "-",
    starting_item         INTEGER               DEFAULT 0,
    floor                 INTEGER               DEFAULT 1,
    FOREIGN KEY(user_id)  REFERENCES users(id),
    FOREIGN KEY(race_id)  REFERENCES races(id)
);
CREATE INDEX race_participants_index_user_id ON race_participants (user_id);
CREATE INDEX race_participants_index_race_id ON race_participants (race_id);

DROP TABLE IF EXISTS race_participant_items;
CREATE TABLE race_participant_items (
    id                                INTEGER                           PRIMARY KEY  AUTOINCREMENT,
    race_participant_id               INTEGER                           NOT NULL,
    item_id                           INTEGER                           NOT NULL,
    floor                             INTEGER                           NOT NULL,
    FOREIGN KEY(race_participant_id)  REFERENCES race_participants(id)
);
CREATE INDEX race_participant_items_index_race_participant_id ON race_participant_items (race_participant_id);

DROP TABLE IF EXISTS banned_users;
CREATE TABLE banned_users (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id                         INTEGER               NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime                        INTEGER               DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(user_id)            REFERENCES users(id)
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);

DROP TABLE IF EXISTS banned_ips;
CREATE TABLE banned_ips (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    ip                              TEXT                  NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime                        INTEGER               DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);

DROP TABLE IF EXISTS squelched_users;
CREATE TABLE squelched_users (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id                         INTEGER               NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime                        INTEGER               DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(user_id)            REFERENCES users(id)
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);

DROP TABLE IF EXISTS chat_log;
CREATE TABLE chat_log (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    room                  TEXT                  NOT NULL,
    user_id               INTEGER               NOT NULL,
    message               TEXT                  NOT NULL,
    datetime              INTEGER               DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(user_id)  REFERENCES users(id)
);

DROP TABLE IF EXISTS chat_log_pm;
CREATE TABLE chat_log_pm (
    id                         INTEGER               PRIMARY KEY  AUTOINCREMENT,
    recipient_id               INTEGER               NOT NULL,
    user_id                    INTEGER               NOT NULL,
    message                    TEXT                  NOT NULL,
    datetime                   INTEGER               DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(user_id)       REFERENCES users(id)
    FOREIGN KEY(recipient_id)  REFERENCES users(id)
);

DROP TABLE IF EXISTS achievements;
CREATE TABLE achievements (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    name                  TEXT                  NOT NULL,
    description           TEXT                  NOT NULL
);

DROP TABLE IF EXISTS user_achievements;
CREATE TABLE user_achievements (
    id                           INTEGER                     PRIMARY KEY  AUTOINCREMENT,
    user_id                      INTEGER                     NOT NULL,
    achievement_id               INTEGER                     NOT NULL,
    datetime                     INTEGER                     DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY(user_id)         REFERENCES users(id)
    FOREIGN KEY(achievement_id)  REFERENCES achievement(id)
);
