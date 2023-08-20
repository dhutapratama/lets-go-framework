package frameworks

import (
	"fmt"
	"lets-go-framework/lets"
	"lets-go-framework/lets/types"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gRPC framework configurations
var GrpcConfig types.IGrpcConfig

// gRPC Server
type grpcServer struct {
	dsn    string
	opts   []grpc.ServerOption
	engine *grpc.Server
	router func(*grpc.Server)
}

// Internal function for initialize gRPC server
func (g *grpcServer) init(config types.IGrpcServer) {
	g.dsn = fmt.Sprintf(":%s", config.GetPort())
	g.engine = grpc.NewServer(g.opts...)
	g.router = config.GetRouter()
}

// Internal function for starting gRPC server
func (rpc *grpcServer) serve() {
	listener, err := net.Listen("tcp", rpc.dsn)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpc.engine.Serve(listener)
}

type grpcClient struct {
	name   string
	dsn    string
	engine *grpc.ClientConn
}

func (rpc *grpcClient) init(config types.IGrpcClient) {
	rpc.name = config.GetName()
	rpc.dsn = fmt.Sprintf("%s:%s", config.GetHost(), config.GetPort())
}

func (rpc *grpcClient) connect() (err error) {
	rpc.engine, err = grpc.Dial(rpc.dsn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

// Run gRPC server and client
func Grpc() {
	if GrpcConfig == nil {
		return
	}

	// Running gRPC server
	if config := GrpcConfig.GetServer(); GrpcConfig.GetServer() != nil {
		lets.LogI("gRPC Server Starting ...")

		var rpcServer grpcServer
		rpcServer.init(config)
		rpcServer.router(rpcServer.engine)
		go rpcServer.serve()
	}

	// Running gRPC client
	if clients := GrpcConfig.GetClients(); len(clients) != 0 {
		lets.LogI("gRPC Client Starting ...")

		for _, config := range clients {
			var rpcClient grpcClient

			lets.LogI("gRPC Client: %s", config.GetName())
			rpcClient.init(config)

			if err := rpcClient.connect(); err != nil {
				lets.LogE("gRPC Client: %s", err.Error())
				lets.LogE("Cannot connect gRPC to %s", config.GetName())
				continue
			}

			for _, isc := range config.GetClients() {
				isc.SetConnection(rpcClient.engine)
			}
		}
	}
}
