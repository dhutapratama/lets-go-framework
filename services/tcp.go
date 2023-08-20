package services

import (
	"lets-go-framework/lets"
)

type TcpHandlerExample struct{}

func (*TcpHandlerExample) OnConnect() {
	lets.LogD("Incoming Connection: Connected")
}

func (*TcpHandlerExample) OnDisconnect() {
	lets.LogD("Incoming Connection: Disconnected")
}
