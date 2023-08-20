package types

// GRPC configuration interface
type ITcpConfig interface {
	GetServers() []ITcpServer
	GetClients() []ITcpClient
}

// GRPC configuration struct
type TcpConfig struct {
	Servers []ITcpServer
	Clients []ITcpClient
}

// Get gRPC server configuration
func (tcp *TcpConfig) GetServers() []ITcpServer {
	return tcp.Servers
}

// Get gRPC client configuration
func (tcp *TcpConfig) GetClients() []ITcpClient {
	return tcp.Clients
}

type ITcpServiceClient interface {
	OnConnect()
	OnDisconnect()
}

type ITcpServiceHandler interface {
	OnConnect()
	OnDisconnect()
}
