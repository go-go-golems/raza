package db

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wesen/raza/lib/raza"
)

type commandIterator struct {
	rows *sql.Rows
}

func (s *commandIterator) Close() {
	if err := s.rows.Close(); err != nil {
		log.Warn().Err(err).Send()
	}
}

func (d *DB) ListCommands() (*commandIterator, error) {
	sb := sqlbuilder.Select(
		"id",
		"session_id",
		"cmd",
		"pwd",
		"start_time_epoch_s",
		"end_time_epoch_s",
	).From("commands")
	sql_, args := sb.Build()

	if log.Debug().Enabled() {
		log.Debug().Str("sql", sql_).Interface("args", args).Msg("GetCommands")
		la := zerolog.Arr()
		for _, arg := range args {
			la = la.Interface(arg)
		}
	}

	rows, err := d.db.Query(sql_, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting sessions from db")
	}

	return &commandIterator{rows: rows}, nil
}

func (s *commandIterator) Next() bool {
	return s.rows.Next()
}

func (s *commandIterator) Get() (*raza.Command, error) {
	var command raza.Command
	err := s.rows.Scan(
		&command.Id,
		&command.SessionId,
		&command.Cmd,
		&command.Pwd,
		&command.StartTimeEpochS,
		&command.EndTimeEpochS)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func (d *DB) StartCommand(ctx context.Context, command *raza.Command) (commandId int64, err error) {
	sb := sqlbuilder.InsertInto("commands").
		Cols("session_id", "cmd", "pwd", "start_time_epoch_s", "end_time_epoch_s")
	sb.Values(command.SessionId, command.Cmd, command.Pwd, command.StartTimeEpochS, command.EndTimeEpochS)

	sql_, args := sb.Build()
	debugLog(sql_, args, "start command")

	result, err := d.db.ExecContext(ctx, sql_, args...)

	if err != nil {
		log.Error().Err(err).Stack().Msg("error inserting new command")
		return -1, err
	}
	commandId, err = result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Stack().Msg("could not get last command insert id")
		return -1, err
	}

	return commandId, nil
}

func (d *DB) GetCommand(ctx context.Context, commandId int64) (*raza.Command, error) {
	sb := sqlbuilder.Select(
		"id",
		"session_id",
		"cmd", "pwd",
		"retval",
		"start_time_epoch_s",
		"end_time_epoch_s",
	).From("commands")
	sb.Where(sb.Equal("id", commandId))

	sql_, args := sb.Build()
	debugLog(sql_, args, "get command")
	row := d.db.QueryRowContext(ctx, sql_, args...)

	var command raza.Command
	err := row.Scan(
		&command.Id,
		&command.SessionId,
		&command.Cmd,
		&command.Pwd,
		&command.Retval,
		&command.StartTimeEpochS,
		&command.EndTimeEpochS)
	if err != nil {
		return nil, err
	}

	// TODO(manuel) query metadata for the command

	return &command, err
}

func (d *DB) EndSessionsLastCommand(ctx context.Context, sessionId int64, retVal int32, endTimeEpochS int64) (*raza.Command, error) {
	sb := sqlbuilder.Select("id").From("commands")
	sb.Where(sb.Equal("session_id", sessionId)).
		OrderBy("start_time_epoch_s ASC").
		Limit(1)

	sql_, args := sb.Build()
	debugLog(sql_, args, "get session's last command")
	row := d.db.QueryRowContext(ctx, sql_, args...)
	var commandId int64
	err := row.Scan(&commandId)
	if err != nil {
		return nil, err
	}

	ub := sqlbuilder.Update("commands")
	ub.Set(ub.Assign("end_time_epoch_s", endTimeEpochS))
	ub.Set(ub.Assign("retval", retVal))
	ub.Where(ub.Equal("id", commandId))

	sql_, args = sb.Build()
	debugLog(sql_, args, "update command")

	return d.GetCommand(ctx, commandId)
}
