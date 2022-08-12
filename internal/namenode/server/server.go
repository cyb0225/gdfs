package server

import (
	"log"
	"net"

	pb "github.com/cyb0225/gdfs/proto/namenode"
	"google.golang.org/grpc"
)

var _ pb.NameNodeServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedNameNodeServer
}

// start rpc server
func RunServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterNameNodeServer(s, &Server{})

	log.Printf("server start listening at %s", port)
	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}
