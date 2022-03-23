package db

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wesen/raza/lib/raza"
)

type sessionIterator struct {
	rows *sql.Rows
}

func (d *DB) ListSessions() (*sessionIterator, error) {
	sb := sqlbuilder.Select(
		"id",
		"host",
		"user",
		"start_time_epoch_s",
		"end_time_epoch_s",
	).From("sessions")
	sql_, args := sb.Build()

	debugLog(sql_, args, "list sessions")

	rows, err := d.db.Query(sql_, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting sessions from db")
	}
	return &sessionIterator{
		rows: rows,
	}, nil
}

func (s *sessionIterator) Next() bool {
	return s.rows.Next()
}

func (s *sessionIterator) Close() {
	if err := s.rows.Close(); err != nil {
		log.Warn().Err(err).Send()
	}
}

func (s *sessionIterator) Get() (*raza.Session, error) {
	var session raza.Session
	err := s.rows.Scan(
		&session.SessionId,
		&session.Host,
		&session.User,
		&session.StartTimeEpochS,
		&session.EndTimeEpochS)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (d *DB) StartSession(ctx context.Context, session *raza.Session) (sessionId int64, err error) {
	sb := sqlbuilder.InsertInto("sessions").
		Cols("host", "user", "start_time_epoch_s", "end_time_epoch_s").
		Values(session.Host, session.User, session.StartTimeEpochS, session.EndTimeEpochS)

	sql_, args := sb.Build()
	debugLog(sql_, args, "start session")

	result, err := d.db.ExecContext(ctx, sql_, args...)
	if err != nil {
		log.Error().Err(err).Stack().Msg("error inserting new session")
		return -1, err
	}
	sessionId, err = result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Stack().Msg("could not get last session insert id")
		return -1, err
	}

	return sessionId, err
}

func (d *DB) GetSession(ctx context.Context, sessionId int64) (*raza.Session, error) {
	sb := sqlbuilder.Select(
		"id",
		"host",
		"user",
		"start_time_epoch_s",
		"end_time_epoch_s",
	).From("sessions")
	sb.Where(sb.Equal("id", sessionId))

	sql_, args := sb.Build()
	debugLog(sql_, args, "get session")
	row := d.db.QueryRowContext(ctx, sql_, args...)

	var session raza.Session
	err := row.Scan(
		&session.SessionId,
		&session.Host,
		&session.User,
		&session.StartTimeEpochS,
		&session.EndTimeEpochS)
	if err != nil {
		return nil, err
	}

	// TODO(manuel) query metadata for the session as well

	return &session, nil
}

func (d *DB) EndSession(ctx context.Context, sessionId int64, endTimeEpochS int64) (*raza.Session, error) {
	sb := sqlbuilder.Update("sessions")
	sb.Set(sb.Assign("end_time_epoch_s", endTimeEpochS))
	sb.Where(sb.Equal("id", sessionId))

	sql_, args := sb.Build()
	debugLog(sql_, args, "update session")

	return d.GetSession(ctx, sessionId)
}
