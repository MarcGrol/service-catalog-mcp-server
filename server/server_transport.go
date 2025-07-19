package server

type ServerTransport interface {
	Start(addr string) error
}