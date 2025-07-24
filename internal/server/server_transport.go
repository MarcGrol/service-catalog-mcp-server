package server

// Transport defines the interface for a server transport.
type Transport interface {
	Start(addr string) error
}
