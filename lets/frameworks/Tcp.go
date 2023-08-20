package frameworks

import (
	"fmt"
	"lets-go-framework/lets"
	"lets-go-framework/lets/types"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

// TCP Configurations
var TcpConfig types.ITcpConfig

// TCP service struct
type tcpServer struct {
	server  string
	engine  net.Listener
	pool    []net.Conn
	logger  func(param gin.LogFormatterParams) string
	handler types.ITcpServiceHandler
}

// Initialize service
func (tcp *tcpServer) init(config types.ITcpServer) {
	var err error

	gin.SetMode(config.GetMode())

	tcp.server = fmt.Sprintf(":%s", config.GetPort())
	if tcp.engine, err = net.Listen("tcp", tcp.server); err != nil {
		lets.LogE("listen: %v", err)
		return
	}

	tcp.logger = func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("%s[HTTP]%s %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			"\x1b[32m", resetColor,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}

	tcp.handler = config.GetHandler()
}

// Run service
func (tcp *tcpServer) serve() (err error) {
	lets.LogI("TCP Server: Listening on %s", tcp.server)

	for {
		// Accept incoming connection
		if connection, err := tcp.engine.Accept(); err != nil {
			lets.LogE("TCP Client: %v", err)

			continue
		} else {
			lets.LogI("TCP Client: connected: %v", connection.RemoteAddr().String())
			tcp.handler.OnConnect()

			tcp.pool = append(tcp.pool, connection)

			// Handle every accepted connection so can used by multiple client
			go func(conn net.Conn, onDisconnect func()) {
				defer conn.Close()
				buf := make([]byte, 1024)

				for {
					if i, err := conn.Read(buf); err != nil {
						lets.LogE("TCP Server: %v", err)
						onDisconnect()
						break
					} else if i > 0 {
						lets.LogW("TCP Server: Unexpected data: %s", buf[:i])
					}
				}
			}(connection, tcp.handler.OnDisconnect)
		}
	}
}

// TCP service struct
type tcpClient struct {
	dsn        string
	connection net.Conn
	logger     func(param gin.LogFormatterParams) string
	clients    []types.ITcpServiceClient
	mode       string
}

// Initialize service
func (tcp *tcpClient) init(config types.ITcpClient) {

	gin.SetMode(config.GetMode())

	tcp.dsn = fmt.Sprintf("%s:%s", config.GetHost(), config.GetPort())

	tcp.logger = func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("%s[HTTP]%s %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			"\x1b[32m", resetColor,
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}

	tcp.clients = config.GetClients()
	tcp.mode = config.GetMode()
}

// Connect
func (tcp *tcpClient) connect() (err error) {
	if tcp.connection, err = net.Dial("tcp", tcp.dsn); err != nil {
		lets.LogE("TCP Client: %v", err)
		return
	}

	if err = tcp.connection.(*net.TCPConn).SetKeepAlive(true); err != nil {
		lets.LogE("TCP Client: set keep alive: %v", err)
		return
	}

	if err = tcp.connection.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second); err != nil {
		lets.LogE("TCP Client: set keep alive period: %v", err)
		return
	}

	// Call event handler
	for _, client := range tcp.clients {
		client.OnConnect()
	}

	// Check connection still alive
	buf := make([]byte, 1024)
	for {
		if i, err := tcp.connection.Read(buf); err != nil {
			tcp.connection.Close() // Close Connection
			lets.LogE("TCP Client: %v", err)

			// Call disconnected
			for _, client := range tcp.clients {
				client.OnDisconnect()
			}
			break
		} else if i > 0 { // Display unexpected incoming data
			lets.LogD("TCP Client: Unexpected data: %s", buf[:i])
		}
	}

	lets.LogW("TCP Client: stopped, auto connect in 10 seconds ...")
	return
}

// Start http service
func Tcp() {
	if TcpConfig == nil {
		return
	}

	// Running TCP server
	if servers := TcpConfig.GetServers(); len(servers) > 0 {
		lets.LogI("TCP Server Starting ...")

		for _, config := range servers {
			var server tcpServer
			server.init(config)

			go func(server tcpServer) {
				for {
					if err := server.serve(); err != nil {
						lets.LogE("TCP Server: %s", err.Error())
					}

					// Auto listen in 10 seconds
					time.Sleep(time.Second * 10)
					lets.LogI("TCP Client: Restarting %s ...", server.server)
				}
			}(server)
		}
	}

	// Running gRPC client
	if clients := TcpConfig.GetClients(); len(clients) != 0 {
		lets.LogI("TCP Client Starting ...")

		for _, config := range clients {
			var client tcpClient

			lets.LogI("TCP Client: %s", config.GetName())
			client.init(config)

			// Connection with auto reconnection
			go func(client tcpClient, config types.ITcpClient) {
				for {
					if err := client.connect(); err != nil {
						lets.LogE("TCP Client: %s", err.Error())
					}

					// Auto Reconnect every 10 seconds
					time.Sleep(time.Second * 10)

					lets.LogI("TCP Client: Restarting %s ...", config.GetName())
				}
			}(client, config)
		}
	}
}
