package services

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

//CrudServer CRUD API Server
type CrudServer struct {
	listener net.Listener
	server   *grpc.Server
}

//InitServer ...
func (c *CrudServer) InitServer(listner net.Listener) {
	c.listener = listner
	opts := []grpc.ServerOption{}
	c.server = grpc.NewServer(opts...)

}

//StartServer ...
func (c *CrudServer) StartServer() {
	if err := c.server.Serve(c.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

//StopServer ...
func (c *CrudServer) StopServer() {
	c.server.Stop()
	c.listener.Close()

}

//GetGRPCServer ...
func (c *CrudServer) GetGRPCServer() *grpc.Server {
	return c.server
}
