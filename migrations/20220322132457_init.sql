-- +goose Up
-- +goose StatementBegin
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

CREATE TABLE commands
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id         INTEGER,
    cmd                TEXT    NOT NULL,
    pwd                TEXT    NOT NULL,
    start_time_epoch_s INTEGER NOT NULL,
    end_time_epoch_s INTEGER,
    FOREIGN KEY (session_id) REFERENCES sessions (id)
);
CREATE INDEX commands_session_idx ON commands (session_id);
CREATE INDEX commands_cmd_idx ON commands (cmd);
CREATE INDEX commands_pwd_idx ON commands (pwd);
CREATE INDEX commands_start_time_epoch_s_idx ON commands (start_time_epoch_s);
CREATE INDEX commands_end_time_epoch_s_idx ON commands (end_time_epoch_s);

CREATE TABLE metadata_entries
(
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    key   TEXT NOT NULL,
    value TEXT NOT NULL
);
CREATE INDEX metadata_entries_key_idx ON metadata_entries (key);
CREATE INDEX metadata_entries_value_idx ON metadata_entries (value);

CREATE TABLE command_metadata_entries
(
    command_id        INTEGER NOT NULL,
    metadata_entry_id INTEGER NOT NULL,
    FOREIGN KEY (command_id) REFERENCES commands (id),
    FOREIGN KEY (metadata_entry_id) REFERENCES metadata_entries (id)
);
CREATE INDEX command_metadata_entries_command_idx ON command_metadata_entries (command_id);
CREATE INDEX command_metadata_entries_metadata_entry_idx ON command_metadata_entries (metadata_entry_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE command_metadata_entries;
DROP TABLE metadata_entries;
DROP TABLE commands;
DROP TABLE sessions;
-- +goose StatementEnd
