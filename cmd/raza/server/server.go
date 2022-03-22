// Package raza implements the commands to interact with the daemon.
package server

import (
	context "context"
	"database/sql"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/wesen/raza/lib/raza"
	"google.golang.org/grpc"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	raza.UnimplementedRazaShellWrapperServer
	raza.UnimplementedRazaQueryServer
	raza.UnimplementedRazaUserServer

	commandId       int64
	sessionCommands map[int64][]*raza.Command
	db              *sql.DB
}

func init() {
	sqlbuilder.DefaultFieldMapper = sqlbuilder.SnakeCaseMapper
}

func (s *server) GetSessions(request *raza.GetSessionsRequest, stream raza.RazaQuery_GetSessionsServer) error {
	sb := sqlbuilder.Select(
		"id",
		"host",
		"user",
		"start_time_epoch_s",
		"end_time_epoch_s",
	).From("sessions")
	sql_, args := sb.Build()

	if log.Debug().Enabled() {
		log.Debug().Str("sql", sql_).Interface("args", args).Msg("GetSessions")
		la := zerolog.Arr()
		for _, arg := range args {
			la = la.Interface(arg)
		}
	}

	rows, err := s.db.Query(sql_, args...)
	if err != nil {
		return errors.Wrapf(err, "error getting sessions from db")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warn().Err(err)
		}
	}()
	for rows.Next() {
		var session raza.Session
		err = rows.Scan(
			&session.SessionId,
			&session.Host,
			&session.User,
			&session.StartTimeEpochS,
			&session.EndTimeEpochS)
		if err != nil {
			log.Error().Err(err).Msg("could not query session")
			break
		}
		err = stream.Send(&session)
		if err != nil {
			return errors.Wrapf(err, "error sending session %d message to stream", session.SessionId)
		}
	}
	return nil
}

func (s *server) GetCommands(request *raza.GetCommandsRequest, stream raza.RazaQuery_GetCommandsServer) error {
	for _, commands := range s.sessionCommands {
		for _, command := range commands {
			err := stream.Send(command)
			if err != nil {
				return errors.Wrapf(err, "error sending command %d message to stream", command.Id)
			}
		}
	}
	return nil
}

func NewServer(db *sql.DB) *server {
	return &server{
		commandId:       0,
		db:              db,
		sessionCommands: make(map[int64][]*raza.Command),
	}
}

func (s *server) StartSession(ctx context.Context, request *raza.StartSessionRequest) (*raza.StartSessionResponse, error) {
	sb := sqlbuilder.InsertInto("sessions").
		Cols("host", "user", "start_time_epoch_s").
		Values(request.Host, request.User, time.Now().Unix())

	sql_, args := sb.Build()

	if log.Debug().Enabled() {
		la := zerolog.Arr()
		for _, arg := range args {
			la = la.Interface(arg)
		}
		log.Debug().Str("sql", sql_).Array("args", la).Msg("insert new session")
	}

	result, err := s.db.ExecContext(ctx, sql_, args...)
	if err != nil {
		log.Error().Err(err).Stack().Msg("error inserting new session")
		return nil, err
	}
	sessionId, err := result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Stack().Msg("could not get last session insert id")
		return nil, err
	}

	return &raza.StartSessionResponse{
		SessionId: sessionId,
	}, nil
}

func (s *server) EndSession(ctx context.Context, request *raza.EndSessionRequest) (*raza.EndSessionResponse, error) {
	sb := sqlbuilder.Update("sessions")
	sb.Set(sb.Assign("end_time_epoch_s", request.EndTimeEpochS))
	sb.Where(sb.Equal("session_id", request.SessionId))

	sql_, args := sb.Build()
	if log.Debug().Enabled() {
		la := zerolog.Arr()
		for _, arg := range args {
			la = la.Interface(arg)
		}
		log.Debug().Str("sql", sql_).Array("args", la).Msg("update session")
	}

	session := &raza.Session{}

	return &raza.EndSessionResponse{
		Session: session,
	}, nil
}

func (s *server) StartCommand(ctx context.Context, request *raza.StartCommandRequest) (*raza.StartCommandResponse, error) {
	commands, ok := s.sessionCommands[request.SessionId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Could not find session with id %d", request.SessionId)
	}

	commandId := s.commandId
	s.commandId += 1

	commands = append(commands,
		&raza.Command{
			Id:              commandId,
			SessionId:       request.SessionId,
			Cmd:             request.Cmd,
			Pwd:             request.Pwd,
			StartTimeEpochS: request.StartTimeEpochS,
			Metadata:        request.Metadata,
		})
	s.sessionCommands[request.SessionId] = commands

	return &raza.StartCommandResponse{
		CommandId: commandId,
	}, nil
}

func (s *server) EndCommand(ctx context.Context, request *raza.EndCommandRequest) (*raza.EndCommandResponse, error) {
	commands, ok := s.sessionCommands[request.SessionId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Could not find session with id %d", request.SessionId)
	}

	if len(commands) <= 0 {
		return nil, status.Errorf(codes.NotFound, "Could not find last command for session %d", request.SessionId)
	}

	command := commands[len(commands)-1]
	command.Retval = request.Retval
	command.EndTimeEpochS = request.EndTimeEpochS

	return &raza.EndCommandResponse{Command: command}, nil
}

var ServerCmd = cobra.Command{
	Use:   "server",
	Short: "Run the raza server",
}

var StartCmd = cobra.Command{
	Use:   "start",
	Short: "Start the raza server",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		lis, err := net.Listen("tcp", address)
		if err != nil {
			cmd.PrintErrf("Failed to listen: %v\n", err)
			return
		}

		s := grpc.NewServer()
		db, err := sql.Open("sqlite3", "./test.sqlite")
		if err != nil {
			log.Fatal().Err(err).Str("file", "./test.sqlite").
				Stack().Msg("failed to open database")
		}
		server := NewServer(db)
		raza.RegisterRazaShellWrapperServer(s, server)
		raza.RegisterRazaQueryServer(s, server)
		raza.RegisterRazaUserServer(s, server)

		log.Info().Str("address", address).Msgf("Starting gRPC listener on %s", address)
		if err := s.Serve(lis); err != nil {
			cmd.PrintErrf("Failed to serve: %v\n", err)
			return
		}
	},
}

func init() {
	ServerCmd.AddCommand(&StartCmd)
}
