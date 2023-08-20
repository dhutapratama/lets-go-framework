package types

import (
	"lets-go-framework/lets"
)

// Default gRPC configuration
const (
	SERVER_TCP_PORT = "5050"
	SERVER_TCP_MODE = "debug"
)

// Interface for accessable method
type ITcpServer interface {
	GetPort() string
	GetMode() string
	GetHandler() ITcpServiceHandler
}

// Serve information
type TcpServer struct {
	Port    string
	Mode    string
	Handler ITcpServiceHandler
}

// Get Port
func (ts *TcpServer) GetPort() string {
	if ts.Port == "" {
		lets.LogW("Config: SERVER_TCP_PORT is not set, using default configuration.")

		return SERVER_TCP_PORT
	}

	return ts.Port
}

// Get Mode
func (ts *TcpServer) GetMode() string {
	if ts.Mode == "" {
		lets.LogW("Config: SERVER_TCP_MODE is not set, using default configuration.")

		return SERVER_TCP_MODE
	}

	return ts.Mode
}

// Get Clients
func (tcp *TcpServer) GetHandler() ITcpServiceHandler {
	return tcp.Handler
}
