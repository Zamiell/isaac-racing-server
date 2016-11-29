/*
    sqlite3 database.sqlite < install/database_schema.sql
*/

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id                         INTEGER    PRIMARY KEY  AUTOINCREMENT,
    auth0_id                   TEXT       NOT NULL,
    username                   TEXT       NOT NULL,
    datetime_created           INTEGER    NOT NULL,
    last_login                 INTEGER    NOT NULL,
    last_ip                    TEXT       NOT NULL,
    admin                      INTEGER    DEFAULT 0,
    unseeded_adjusted_average  INTEGER    DEFAULT 0,
    unseeded_real_average      INTEGER    DEFAULT 0,
    num_unseeded_races         INTEGER    DEFAULT 0,
    num_forfeits               INTEGER    DEFAULT 0,
    forfeit_penalty            INTEGER    DEFAULT 0,
    lowest_unseeded_time       INTEGER    DEFAULT 0,
    last_unseeded_race         INTEGER    DEFAULT 0,
    elo                        INTEGER    DEFAULT 0,
    last_elo_change            INTEGER    DEFAULT 0,
    num_seeded_races           INTEGER    DEFAULT 0,
    last_seeded_race           INTEGER    DEFAULT 0
);
CREATE UNIQUE INDEX users_index_auth0_id ON users (auth0_id);
CREATE UNIQUE INDEX users_index_username ON users (username COLLATE NOCASE);

DROP TABLE IF EXISTS races;
CREATE TABLE races (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    name                  TEXT                  DEFAULT "-",
    status                TEXT                  DEFAULT "open", /* starting, in progress, finished */
    format                TEXT                  DEFAULT "unseeded", /* seeded, diversity */
    character             TEXT                  DEFAULT "Isaac", /* Isaac, Magdalene, Cain, Judas, Blue Baby, Eve, Samson, Azazel, Lazarus, Eden, The Lost, Lilith, Keeper */
    goal                  TEXT                  DEFAULT "Blue Baby", /* The Lamb, Mega Satan */
    starting_build        INTEGER               DEFAULT -1, /* -1 for unseeded/diversity races, setting it to 0 means "keep it as it is" */
    seed                  TEXT                  DEFAULT "-",
    captain               INTEGER               NOT NULL,
    datetime_created      INTEGER               NOT NULL,
    datetime_started      INTEGER               DEFAULT 0,
    datetime_finished     INTEGER               DEFAULT 0,
    FOREIGN KEY(captain)  REFERENCES users(id)
);
CREATE INDEX races_index_status ON races (status);
CREATE INDEX races_index_datetime_finished ON races (datetime_finished);

DROP TABLE IF EXISTS race_participants;
CREATE TABLE race_participants (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id               INTEGER               NOT NULL,
    race_id               INTEGER               NOT NULL,
    status                TEXT                  DEFAULT "not ready", /* ready, racing, finished, quit, disqualified */
    datetime_joined       INTEGER               NOT NULL,
    datetime_finished     INTEGER               DEFAULT 0,
    place                 INTEGER               DEFAULT 0,
    comment               TEXT                  DEFAULT "-",
    seed                  TEXT                  DEFAULT "-",
    starting_item         INTEGER               DEFAULT 0,
    floor                 TEXT                  DEFAULT "1-0",
    FOREIGN KEY(user_id)  REFERENCES users(id),
    FOREIGN KEY(race_id)  REFERENCES races(id),
    UNIQUE(user_id, race_id)
);
CREATE INDEX race_participants_index_user_id ON race_participants (user_id);
CREATE INDEX race_participants_index_race_id ON race_participants (race_id);
CREATE INDEX race_participants_index_datetime_joined ON race_participants (datetime_joined);

DROP TABLE IF EXISTS race_participant_items;
CREATE TABLE race_participant_items (
    id                                INTEGER                           PRIMARY KEY  AUTOINCREMENT,
    race_participant_id               INTEGER                           NOT NULL,
    item_id                           INTEGER                           NOT NULL,
    floor                             TEXT                              NOT NULL,
    FOREIGN KEY(race_participant_id)  REFERENCES race_participants(id)
);
CREATE INDEX race_participant_items_index_race_participant_id ON race_participant_items (race_participant_id);

DROP TABLE IF EXISTS banned_users;
CREATE TABLE banned_users (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id                         INTEGER               NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime_banned                 INTEGER               NOT NULL,
    FOREIGN KEY(user_id)            REFERENCES users(id),
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);
CREATE UNIQUE INDEX banned_users_index_user_id ON banned_users (user_id);

DROP TABLE IF EXISTS banned_ips;
CREATE TABLE banned_ips (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    ip                              TEXT                  NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime_banned                 INTEGER               NOT NULL,
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);
CREATE UNIQUE INDEX banned_ips_index_ip ON banned_ips (ip);

DROP TABLE IF EXISTS squelched_users;
CREATE TABLE squelched_users (
    id                              INTEGER               PRIMARY KEY  AUTOINCREMENT,
    user_id                         INTEGER               NOT NULL,
    admin_responsible               INTEGER               NOT NULL,
    datetime_squelched              INTEGER               NOT NULL,
    FOREIGN KEY(user_id)            REFERENCES users(id),
    FOREIGN KEY(admin_responsible)  REFERENCES users(id)
);
CREATE UNIQUE INDEX squelched_users_index_user_id ON squelched_users (user_id);

DROP TABLE IF EXISTS chat_log;
CREATE TABLE chat_log (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    room                  TEXT                  NOT NULL,
    user_id               INTEGER               NOT NULL,
    message               TEXT                  NOT NULL,
    datetime_sent         INTEGER               NOT NULL,
    FOREIGN KEY(user_id)  REFERENCES users(id)
);
CREATE INDEX chat_log_index_room ON chat_log (room);
CREATE INDEX chat_log_index_user_id ON chat_log (user_id);
CREATE INDEX chat_log_index_datetime ON chat_log (datetime_sent);

DROP TABLE IF EXISTS chat_log_pm;
CREATE TABLE chat_log_pm (
    id                         INTEGER               PRIMARY KEY  AUTOINCREMENT,
    recipient_id               INTEGER               NOT NULL,
    user_id                    INTEGER               NOT NULL,
    message                    TEXT                  NOT NULL,
    datetime_sent              INTEGER               NOT NULL,
    FOREIGN KEY(user_id)       REFERENCES users(id),
    FOREIGN KEY(recipient_id)  REFERENCES users(id)
);
CREATE INDEX chat_log_pm_index_recipient_id ON chat_log_pm (recipient_id);
CREATE INDEX chat_log_pm_index_user_id ON chat_log_pm (user_id);
CREATE INDEX chat_log_pm_index_datetime ON chat_log_pm (datetime_sent);

DROP TABLE IF EXISTS achievements;
CREATE TABLE achievements (
    id                    INTEGER               PRIMARY KEY  AUTOINCREMENT,
    name                  TEXT                  NOT NULL,
    description           TEXT                  NOT NULL
);
CREATE UNIQUE INDEX achievements_index_name ON achievements (name COLLATE NOCASE);

DROP TABLE IF EXISTS user_achievements;
CREATE TABLE user_achievements (
    id                           INTEGER                     PRIMARY KEY  AUTOINCREMENT,
    user_id                      INTEGER                     NOT NULL,
    achievement_id               INTEGER                     NOT NULL,
    datetime_achieved            INTEGER                     NOT NULL,
    FOREIGN KEY(user_id)         REFERENCES users(id),
    FOREIGN KEY(achievement_id)  REFERENCES achievement(id),
    UNIQUE(user_id, achievement_id)
);
CREATE INDEX user_achievements_index_user_id ON user_achievements (user_id);
CREATE INDEX user_achievements_index_achievement_id ON user_achievements (achievement_id);

DROP TABLE IF EXISTS seeds;
CREATE TABLE seeds (
    id    INTEGER  PRIMARY KEY  AUTOINCREMENT,
    seed  TEXT     NOT NULL
);
CREATE UNIQUE INDEX seeds_index_seed ON seeds (seed COLLATE NOCASE);
