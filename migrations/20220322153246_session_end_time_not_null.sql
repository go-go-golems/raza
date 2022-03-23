--- this migration is more for the fun of writing a migration for sqlite3 rather than really necessary.
--- I folded a similar migration into init.sql for commands.

-- +goose Up
-- +goose StatementBegin
CREATE TEMPORARY TABLE sessions_backup
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    host               TEXT NOT NULL,
    user               TEXT NOT NULL,
    start_time_epoch_s INTEGER,
    end_time_epoch_s   INTEGER
);
INSERT INTO sessions_backup
SELECT id, host, user, start_time_epoch_s, end_time_epoch_s
FROM sessions;
DROP TABLE sessions;

UPDATE sessions_backup
SET start_time_epoch_s = datetime ('now', 'localtime')
WHERE start_time_epoch_s IS NULL;

UPDATE sessions_backup
SET end_time_epoch_s = start_time_epoch_s
WHERE end_time_epoch_s IS NULL;

CREATE TABLE sessions
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    host               TEXT NOT NULL,
    user               TEXT NOT NULL,
    start_time_epoch_s INTEGER NOT NULL,
    end_time_epoch_s   INTEGER NOT NULL
);

CREATE INDEX sessions_host_idx ON sessions (host);
CREATE INDEX sessions_user_idx ON sessions (user);
CREATE INDEX sessions_start_time_epoch_s_idx ON sessions (start_time_epoch_s);
CREATE INDEX sessions_end_time_epoch_s_idx ON sessions (end_time_epoch_s);
INSERT INTO sessions
SELECT id, host, user, start_time_epoch_s, end_time_epoch_s
FROM sessions_backup;
DROP TABLE sessions_backup;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TEMPORARY TABLE sessions_backup
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    host               TEXT NOT NULL,
    user               TEXT NOT NULL,
    start_time_epoch_s INTEGER,
    end_time_epoch_s   INTEGER
);
INSERT INTO sessions_backup
SELECT id, host, user, start_time_epoch_s, end_time_epoch_s
FROM sessions;
DROP TABLE sessions;

CREATE TABLE sessions
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    host               TEXT NOT NULL,
    user               TEXT NOT NULL,
    start_time_epoch_s INTEGER,
    end_time_epoch_s   INTEGER
);

CREATE INDEX sessions_host_idx ON sessions (host);
CREATE INDEX sessions_user_idx ON sessions (user);
CREATE INDEX sessions_start_time_epoch_s_idx ON sessions (start_time_epoch_s);
CREATE INDEX sessions_end_time_epoch_s_idx ON sessions (end_time_epoch_s);
INSERT INTO sessions
SELECT id, host, user, start_time_epoch_s, end_time_epoch_s
FROM sessions_backup;
DROP TABLE sessions_backup;
-- +goose StatementEnd
