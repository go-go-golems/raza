// Package raza implements the commands to interact with the daemon.
package server

import (
	context "context"
	"database/sql"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/wesen/raza/lib/db"
	"github.com/wesen/raza/lib/raza"
	"google.golang.org/grpc"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	sqlbuilder.DefaultFieldMapper = sqlbuilder.SnakeCaseMapper
}

func (s *server) GetSessions(request *raza.GetSessionsRequest, stream raza.RazaQuery_GetSessionsServer) error {
	// TODO(manuel) actually add query parameters, pagination, etc.
	_ = request
	i, err := s.db.ListSessions()
	if err != nil {
		log.Error().Err(err).Msg("could not query sessions")
		return err
	}
	defer i.Close()
	for i.Next() {
		session, err := i.Get()
		if session == nil {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("could not query session")
			break
		}
		err = stream.Send(session)
		if err != nil {
			return errors.Wrapf(err, "error sending session %d message to stream", session.SessionId)
		}
	}
	return nil
}

func (s *server) StartSession(ctx context.Context, request *raza.StartSessionRequest) (*raza.StartSessionResponse, error) {
	session := &raza.Session{
		SessionId:       -1,
		Host:            request.Host,
		User:            request.User,
		StartTimeEpochS: request.StartTimeEpochS,
		// we close the session immediately so that we don't have open intervals in the DB
		EndTimeEpochS: request.StartTimeEpochS,
	}
	sessionId, err := s.db.StartSession(ctx, session)
	if err != nil {
		log.Error().Err(err).Msg("could not start session")
		return nil, err
	}

	return &raza.StartSessionResponse{
		SessionId: sessionId,
	}, nil
}

func (s *server) EndSession(ctx context.Context, request *raza.EndSessionRequest) (*raza.EndSessionResponse, error) {
	session, err := s.db.EndSession(ctx, request.SessionId, request.EndTimeEpochS)
	if err != nil {
		log.Error().Err(err).Msg("could not end session")
		return nil, err
	}

	return &raza.EndSessionResponse{
		Session: session,
	}, nil
}

func (s *server) GetCommands(request *raza.GetCommandsRequest, stream raza.RazaQuery_GetCommandsServer) error {
	// TODO(manuel) actually add query parameters, pagination, etc...
	_ = request
	i, err := s.db.ListCommands()
	if err != nil {
		log.Error().Err(err).Msg("could not query commands")
		return err
	}
	defer i.Close()

	for i.Next() {
		command, err := i.Get()
		if err != nil {
			log.Error().Err(err).Msg("could not get command")
		}
		err = stream.Send(command)
		if err != nil {
			return errors.Wrapf(err, "error sending command %d message to stream", command.Id)
		}
	}
	return nil
}

func (s *server) StartCommand(ctx context.Context, request *raza.StartCommandRequest) (*raza.StartCommandResponse, error) {
	command := &raza.Command{
		Id:              -1,
		SessionId:       request.SessionId,
		Cmd:             request.Cmd,
		Pwd:             request.Pwd,
		StartTimeEpochS: request.StartTimeEpochS,
		// we store the command as having immediately ended for now, to avoid having open intervals
		EndTimeEpochS: request.StartTimeEpochS,
		Metadata:      request.Metadata,
	}

	commandId, err := s.db.StartCommand(ctx, command)
	if err != nil {
		log.Error().Err(err).Stack().Msg("error starting command")
		return nil, err
	}

	return &raza.StartCommandResponse{
		CommandId: commandId,
	}, nil
}

func (s *server) EndSessionsLastCommand(ctx context.Context, request *raza.EndSessionsLastCommandRequest) (*raza.EndCommandResponse, error) {
	command, err := s.db.EndSessionsLastCommand(ctx, request.SessionId, request.Retval, request.EndTimeEpochS)
	if err != nil {
		log.Error().Err(err).Msg("could not end command")
		return nil, err
	}

	return &raza.EndCommandResponse{Command: command}, nil
}

type server struct {
	raza.UnimplementedRazaShellWrapperServer
	raza.UnimplementedRazaQueryServer
	raza.UnimplementedRazaUserServer

	db *db.DB
}

func NewServer(db_ *sql.DB) *server {
	return &server{
		db: db.NewDB(db_),
	}
}

func (s *server) Close() {
	s.db.Close()
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
		defer server.Close()
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
