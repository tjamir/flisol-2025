package test

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const BufSize = 1024 * 1024

type TestGRPCServer struct {
	Listener *bufconn.Listener
	Server   *grpc.Server
}

func StartGRPCServer(registerFn func(*grpc.Server)) *TestGRPCServer {
	listener := bufconn.Listen(BufSize)
	server := grpc.NewServer()

	registerFn(server)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Erro no servidor de teste: %v", err)
		}
	}()

	return &TestGRPCServer{
		Listener: listener,
		Server:   server,
	}
}

func (t *TestGRPCServer) Dialer() func(ctx context.Context, s string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return t.Listener.Dial()
	}
}
