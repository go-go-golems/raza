-- +goose Up
-- +goose StatementBegin
ALTER TABLE commands ADD COLUMN retval INTEGER NOT NULL DEFAULT -1;
CREATE INDEX commands_retval_idx ON commands (retval);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TEMPORARY TABLE commands_backup
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id         INTEGER,
    cmd                TEXT    NOT NULL,
    pwd                TEXT    NOT NULL,
    start_time_epoch_s INTEGER NOT NULL,
    end_time_epoch_s INTEGER NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions (id)
);

INSERT INTO commands_backup
SELECT id, session_id, cmd, pwd, start_time_epoch_s, end_time_epoch_s
FROM commands;

DROP TABLE commands;

CREATE TABLE commands
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id         INTEGER,
    cmd                TEXT    NOT NULL,
    pwd                TEXT    NOT NULL,
    start_time_epoch_s INTEGER NOT NULL,
    end_time_epoch_s INTEGER NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions (id)
);
CREATE INDEX commands_session_idx ON commands (session_id);
CREATE INDEX commands_cmd_idx ON commands (cmd);
CREATE INDEX commands_pwd_idx ON commands (pwd);
CREATE INDEX commands_start_time_epoch_s_idx ON commands (start_time_epoch_s);
CREATE INDEX commands_end_time_epoch_s_idx ON commands (end_time_epoch_s);

INSERT INTO commands
SELECT id, session_id, cmd, pwd, start_time_epoch_s, end_time_epoch_s
FROM commands_backup;
DROP TABLE commands_backup;
-- +goose StatementEnd
