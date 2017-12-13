/*
    Setting up the database is covered in the README.md file.

    The length of NVARCHAR and VARCHAR columns was deliberately chosen to be liberal;
    application-level constraints limit these values to be much smaller than they are expressed here.
*/

USE isaac;

/*
    We have to disable foreign key checks so that we can drop the tables;
    this will only disable it for the current session
 */
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    /* Main values */
    id                   INT           NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    steam_id             BIGINT        NOT NULL  UNIQUE, /* All authentication is through Steam, so we don't need to store a password for the user */
    /* steam_id must be a BIGINT, as they are 17 digits long */
    username             NVARCHAR(50)  NOT NULL  UNIQUE, /* MariaDB is case insensitive by default, which is what we want */
    datetime_created     TIMESTAMP     NOT NULL  DEFAULT NOW(),
    datetime_last_login  TIMESTAMP     NOT NULL  DEFAULT NOW(),
    last_ip              VARCHAR(40)   NOT NULL,
    admin                INT           NOT NULL  DEFAULT 0, /* 0 is not an admin, 1 is a staff, 2 is a full administrator */
    verified             TINYINT(1)    NOT NULL  DEFAULT 0, /* Used to show who is a legitimate player on the leaderboard */

    /* Seeded leaderboard values */
    seeded_trueskill        FLOAT      NOT NULL  DEFAULT 0,
    seeded_trueskill_mu     FLOAT      NOT NULL  DEFAULT 25,
    seeded_trueskill_sigma  FLOAT      NOT NULL  DEFAULT 8.333,
    seeded_num_races        INT        NOT NULL  DEFAULT 0,
    seeded_last_race        TIMESTAMP  NULL      DEFAULT NULL,

    /* Seeded solo leaderboard values */
    seeded_solo_trueskill        FLOAT      NOT NULL  DEFAULT 0,
    seeded_solo_trueskill_mu     FLOAT      NOT NULL  DEFAULT 25,
    seeded_solo_trueskill_sigma  FLOAT      NOT NULL  DEFAULT 8.333,
    seeded_solo_num_races        INT        NOT NULL  DEFAULT 0,
    seeded_solo_last_race        TIMESTAMP  NULL      DEFAULT NULL,

    /* Unseeded leaderboard values */
    unseeded_trueskill        FLOAT      NOT NULL  DEFAULT 0,
    unseeded_trueskill_mu     FLOAT      NOT NULL  DEFAULT 25,
    unseeded_trueskill_sigma  FLOAT      NOT NULL  DEFAULT 8.333,
    unseeded_num_races        INT        NOT NULL  DEFAULT 0,
    unseeded_last_race        TIMESTAMP  NULL      DEFAULT NULL,

    /* Unseeded solo leaderboard values */
    unseeded_solo_adjusted_average  INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_solo_real_average      INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_solo_num_races         INT        NOT NULL  DEFAULT 0,
    unseeded_solo_num_forfeits      INT        NOT NULL  DEFAULT 0,
    unseeded_solo_forfeit_penalty   INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_solo_lowest_time       INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_solo_last_race         TIMESTAMP  NULL      DEFAULT NULL,

    /* Diversity leaderboard values */
    diversity_trueskill         FLOAT      NOT NULL  DEFAULT 25,
    diversity_trueskill_sigma   FLOAT      NOT NULL  DEFAULT 8.333,
    diversity_trueskill_change  FLOAT      NOT NULL  DEFAULT 0, /* The amount changed in the last race (can be positive or negative) */
    diversity_num_races         INT        NOT NULL  DEFAULT 0,
    diversity_last_race         TIMESTAMP  NULL      DEFAULT NULL,

    /* Stream values */
    stream_url                 NVARCHAR(50)  NOT NULL  DEFAULT "-", /* Their stream URL */
    twitch_bot_enabled         TINYINT(1)    NOT NULL  DEFAULT 0, /* Either 0 or 1 */
    twitch_bot_delay           INT           NOT NULL  DEFAULT 15 /* Between 0 and 60 */
);
CREATE UNIQUE INDEX users_index_steam_id ON users (steam_id);
CREATE UNIQUE INDEX users_index_username ON users (username);
INSERT INTO users (steam_id, username, last_ip) VALUES (0, "[SERVER]", "-");

DROP TABLE IF EXISTS races;
CREATE TABLE races (
    /* Main fields */
    id       INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    finished TINYINT(1)     NOT NULL  DEFAULT 0,
    name     NVARCHAR(100)  NULL,

    /* Format */
    ranked          TINYINT(1)   NULL, /* 0 for unranked, 1 for ranked */
    solo            TINYINT(1)   NULL, /* 0 for solo, 1 for multiplayer */
    format          VARCHAR(50)  NULL, /* unseeded, seeded, diversity, unseeded-lite, custom */
    player_type     VARCHAR(50)  NULL, /* You can't name columns "character" in MariaDB, so we will use the Lua name for this instead */
    /* Isaac, Magdalene, Cain, Judas, Blue Baby, Eve, Samson, Azazel, Lazarus, Eden, The Lost, Lilith, Keeper, Samael */
    goal            VARCHAR(50)  NULL, /* Blue Baby, The Lamb, Mega Satan, Everything, custom */
    starting_build  INT          NULL  DEFAULT -1, /* -1 for unseeded & diversity races, otherwise matches the build number */

    /* Other fields */
    seed               VARCHAR(50)  NULL      DEFAULT "-",
    captain            INT          NULL,
    datetime_created   TIMESTAMP    NOT NULL  DEFAULT NOW(),
    datetime_started   TIMESTAMP    NULL      DEFAULT NULL,
    datetime_finished  TIMESTAMP    NULL      DEFAULT NULL,

    FOREIGN KEY(captain) REFERENCES users(id)
);
CREATE INDEX races_index_datetime_finished ON races (datetime_finished);

DROP TABLE IF EXISTS race_participants;
CREATE TABLE race_participants (
    id                 INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    user_id            INT            NOT NULL,
    race_id            INT            NOT NULL,
    datetime_joined    TIMESTAMP      NOT NULL  DEFAULT 0,
    seed               VARCHAR(50)    NOT NULL,
    starting_item      INT            NOT NULL,
    place              INT            NOT NULL, /* -1 is quit, -2 is disqualified */
    datetime_finished  TIMESTAMP      NOT NULL  DEFAULT 0,
    run_time           INT            NOT NULL, /* in milliseconds */
    comment            NVARCHAR(150)  NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(race_id) REFERENCES races(id) ON DELETE CASCADE,
    /* If the race is deleted, automatically delete all of the race participant rows */
    UNIQUE(user_id, race_id)
);
CREATE INDEX race_participants_index_user_id ON race_participants (user_id);
CREATE INDEX race_participants_index_race_id ON race_participants (race_id);
CREATE INDEX race_participants_index_datetime_joined ON race_participants (datetime_joined);

DROP TABLE IF EXISTS race_participant_items;
CREATE TABLE race_participant_items (
    id                   INT        NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    race_participant_id  INT        NOT NULL,
    item_id              INT        NOT NULL,
    floor_num            INT        NOT NULL,
    stage_type           INT        NOT NULL,
    datetime_acquired    TIMESTAMP  NOT NULL,

    FOREIGN KEY(race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
    /* If the race participant entry is deleted, automatically delete all of their items */
);
CREATE INDEX race_participant_items_index_race_participant_id ON race_participant_items (race_participant_id);

DROP TABLE IF EXISTS race_participant_rooms;
CREATE TABLE race_participant_rooms (
    id                   INT          NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    race_participant_id  INT          NOT NULL,
    room_id              VARCHAR(50)  NOT NULL,
    floor_num            INT          NOT NULL,
    stage_type           INT          NOT NULL,
    datetime_arrived     TIMESTAMP    NOT NULL,

    FOREIGN KEY(race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
    /* If the race participant entry is deleted, automatically delete all of their rooms */
);
CREATE INDEX race_participant_rooms_index_race_participant_id ON race_participant_rooms (race_participant_id);

DROP TABLE IF EXISTS banned_users;
CREATE TABLE banned_users (
    id                 INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    user_id            INT            NOT NULL,
    admin_responsible  INT            NOT NULL,
    reason             NVARCHAR(150)  NOT NULL  DEFAULT "-",
    datetime_banned    TIMESTAMP      NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE, /* If the user is deleted, automatically delete the banned_users entry */
    FOREIGN KEY(admin_responsible) REFERENCES users(id)
);
CREATE UNIQUE INDEX banned_users_index_user_id ON banned_users (user_id);

DROP TABLE IF EXISTS banned_ips;
CREATE TABLE banned_ips (
    id                 INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    ip                 VARCHAR(40)    NOT NULL,
    user_id            INT            NULL      DEFAULT NULL, /* If specified, this IP address is associated with the respective user */
    admin_responsible  INT            NOT NULL,
    reason             NVARCHAR(150)  NOT NULL  DEFAULT "-",
    datetime_banned    TIMESTAMP      NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE, /* If the user is deleted, automatically delete the banned_ips entry */
    FOREIGN KEY(admin_responsible) REFERENCES users(id)
);
CREATE UNIQUE INDEX banned_ips_index_ip ON banned_ips (ip);
CREATE UNIQUE INDEX banned_ips_index_user_id ON banned_ips (user_id);

DROP TABLE IF EXISTS muted_users;
CREATE TABLE muted_users (
    id                 INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    user_id            INT            NOT NULL,
    admin_responsible  INT            NOT NULL,
    reason             NVARCHAR(150)  NOT NULL  DEFAULT "-",
    datetime_muted     TIMESTAMP      NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE, /* If the user is deleted, automatically delete the muted_users entry */
    FOREIGN KEY(admin_responsible) REFERENCES users(id)
);
CREATE UNIQUE INDEX muted_users_index_user_id ON muted_users (user_id);

DROP TABLE IF EXISTS chat_log;
CREATE TABLE chat_log (
    id             INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    room           VARCHAR(50)    NOT NULL,
    user_id        INT            NOT NULL,
    message        NVARCHAR(200)  NOT NULL,
    datetime_sent  TIMESTAMP      NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX chat_log_index_room ON chat_log (room);
CREATE INDEX chat_log_index_user_id ON chat_log (user_id);
CREATE INDEX chat_log_index_datetime ON chat_log (datetime_sent);

DROP TABLE IF EXISTS chat_log_pm;
CREATE TABLE chat_log_pm (
    id             INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    recipient_id   INT            NOT NULL,
    user_id        INT            NOT NULL,
    message        NVARCHAR(500)  NOT NULL,
    datetime_sent  TIMESTAMP      NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(recipient_id) REFERENCES users(id)
);
CREATE INDEX chat_log_pm_index_recipient_id ON chat_log_pm (recipient_id);
CREATE INDEX chat_log_pm_index_user_id ON chat_log_pm (user_id);
CREATE INDEX chat_log_pm_index_datetime ON chat_log_pm (datetime_sent);

DROP TABLE IF EXISTS achievements;
CREATE TABLE achievements (
    id           INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    name         NVARCHAR(100)  NOT NULL,
    description  NVARCHAR(300)  NOT NULL
);
CREATE UNIQUE INDEX achievements_index_name ON achievements (name);

DROP TABLE IF EXISTS user_achievements;
CREATE TABLE user_achievements (
    id                 INT        NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    user_id            INT        NOT NULL,
    achievement_id     INT        NOT NULL,
    datetime_achieved  TIMESTAMP  NOT NULL  DEFAULT NOW(),

    FOREIGN KEY(user_id) REFERENCES users(id),
    UNIQUE(user_id, achievement_id)
);
CREATE INDEX user_achievements_index_user_id ON user_achievements (user_id);
CREATE INDEX user_achievements_index_achievement_id ON user_achievements (achievement_id);

DROP TABLE IF EXISTS user_season_stats;
CREATE TABLE user_season_stats (
    /* Main values */
    id       INT  NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    user_id  INT  NOT NULL,
    season   INT  NOT NULL,

    /* Seeded leaderboard values */
    seeded_trueskill         FLOAT      NOT NULL  DEFAULT 0,
    seeded_num_races         INT        NOT NULL  DEFAULT 0,
    seeded_last_race         TIMESTAMP  NULL      DEFAULT NULL,

    /* Unseeded leaderboard values */
    unseeded_adjusted_average  INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_real_average      INT        NOT NULL  DEFAULT 0, /* Rounded to the second */
    unseeded_num_races         INT        NOT NULL  DEFAULT 0,
    unseeded_num_forfeits      INT        NOT NULL  DEFAULT 0,
    unseeded_forfeit_penalty   INT        NOT NULL  DEFAULT 0, /* Rounded to the second */

    FOREIGN KEY(user_id) REFERENCES users(id),
    UNIQUE(user_id, season)
);
CREATE INDEX user_season_stats_index_user_id ON user_season_stats (user_id);
