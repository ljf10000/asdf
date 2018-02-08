package asdf

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGrpcServer(port string, register func(server *grpc.Server)) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	register(server)
	reflection.Register(server)

	if err := server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func GrpcCall(server string, call func(conn *grpc.ClientConn) error) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}

	err = call(conn)

	conn.Close()

	return err
}
