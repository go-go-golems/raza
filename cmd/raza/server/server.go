// Package raza implements the commands to interact with the daemon.
package server

import (
	context "context"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/wesen/raza/lib/raza"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type server struct {
	raza.UnimplementedRazaShellWrapperServer
	raza.UnimplementedRazaQueryServer
	raza.UnimplementedRazaUserServer

	sessionId       int64
	commandId       int64
	sessions        map[int64]*raza.Session
	sessionCommands map[int64][]*raza.Command
}

func (s *server) GetSessions(request *raza.GetSessionsRequest, stream raza.RazaQuery_GetSessionsServer) error {
	for _, session := range s.sessions {
		err := stream.Send(session)
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

func NewServer() *server {
	return &server{
		commandId:       0,
		sessions:        make(map[int64]*raza.Session),
		sessionCommands: make(map[int64][]*raza.Command),
	}
}

func (s *server) StartSession(ctx context.Context, request *raza.StartSessionRequest) (*raza.StartSessionResponse, error) {
	sessionId := s.sessionId
	s.sessions[sessionId] = &raza.Session{
		SessionId:       sessionId,
		Host:            request.Host,
		User:            request.User,
		StartTimeEpochS: request.StartTimeEpochS,
	}
	s.sessionCommands[sessionId] = []*raza.Command{}
	s.sessionId += 1

	return &raza.StartSessionResponse{
		SessionId: sessionId,
	}, nil
}

func (s *server) EndSession(ctx context.Context, request *raza.EndSessionRequest) (*raza.EndSessionResponse, error) {
	session, ok := s.sessions[request.SessionId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Could not find session with id %d", request.SessionId)
	}

	session.EndTimeEpochS = request.EndTimeEpochS

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
		server := NewServer()
		raza.RegisterRazaShellWrapperServer(s, server)
		raza.RegisterRazaQueryServer(s, server)
		raza.RegisterRazaUserServer(s, server)

		log.Printf("Starting gRPC listener on %s", address)
		if err := s.Serve(lis); err != nil {
			cmd.PrintErrf("Failed to serve: %v\n", err)
			return
		}
	},
}

func init() {
	ServerCmd.AddCommand(&StartCmd)
}
