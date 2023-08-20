package types

import "lets-go-framework/lets"

// Default grpc configuration
const (
	CLIENT_TCP_NAME = "lets-connect-tcp"
	CLIENT_TCP_HOST = "127.0.0.1"
	CLIENT_TCP_PORT = "5050"
	CLIENT_TCP_MODE = "debug"
)

// Interface for grpc method
type ITcpClient interface {
	GetName() string
	GetHost() string
	GetPort() string
	GetMode() string
	GetClients() []ITcpServiceClient
}

// Client information
type TcpClient struct {
	Name    string
	Host    string
	Port    string
	Mode    string
	Clients []ITcpServiceClient
}

// Get Name
func (tcp *TcpClient) GetName() string {
	if tcp.Name == "" {
		lets.LogW("Configs: CLIENT_TCP_NAME is not set, using default configuration.")

		return CLIENT_TCP_NAME
	}

	return tcp.Name
}

// Get Mode
func (tcp *TcpClient) GetMode() string {
	if tcp.Mode == "" {
		lets.LogW("Config: CLIENT_TCP_MODE is not set, using default configuration.")

		return CLIENT_TCP_MODE
	}

	return tcp.Mode
}

// Get Host
func (tcp *TcpClient) GetHost() string {
	if tcp.Host == "" {
		lets.LogW("Configs: CLIENT_TCP_HOST is not set, using default configuration.")

		return CLIENT_TCP_HOST
	}

	return tcp.Host
}

// Get Port
func (tcp *TcpClient) GetPort() string {
	if tcp.Port == "" {
		lets.LogW("Configs: CLIENT_TCP_PORT is not set, using default configuration.")

		return CLIENT_TCP_PORT
	}

	return tcp.Port
}

// Get Clients
func (tcp *TcpClient) GetClients() []ITcpServiceClient {
	return tcp.Clients
}
