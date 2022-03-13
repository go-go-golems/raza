package helpers

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/wesen/raza/lib/raza"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func ConnectToRazaServer(cmd *cobra.Command) (*grpc.ClientConn, error) {
	address, _ := cmd.Flags().GetString("address")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cmd.PrintErrf("Could not connect to grpc server at address: %s\n", address)
		return nil, err
	}
	return conn, nil
}

func NewRazaQueryClient(cmd *cobra.Command) (raza.RazaQueryClient, context.Context, func(), error) {
	conn, err := ConnectToRazaServer(cmd)
	if err != nil {
		return nil, nil, func() {}, err
	}

	c := raza.NewRazaQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return c, ctx, func() {
		cancel()
		conn.Close()
	}, nil
}

func NewRazaShellWrapperClient(cmd *cobra.Command) (raza.RazaShellWrapperClient, context.Context, func(), error) {
	conn, err := ConnectToRazaServer(cmd)
	if err != nil {
		return nil, nil, func() {}, err
	}

	c := raza.NewRazaShellWrapperClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return c, ctx, func() {
		cancel()
		conn.Close()
	}, nil
}

func NewRazaUserClient(cmd *cobra.Command) (raza.RazaUserClient, context.Context, func(), error) {
	conn, err := ConnectToRazaServer(cmd)
	if err != nil {
		return nil, nil, func() {}, err
	}

	c := raza.NewRazaUserClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return c, ctx, func() {
		cancel()
		conn.Close()
	}, nil
}
