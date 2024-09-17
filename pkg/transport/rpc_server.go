package transport

import (
	"crypto/tls"
	"log"
	"net"
	"net/rpc"
)

type RPCServer struct {
	address string
}

func NewRPCServer(address string) *RPCServer {
	return &RPCServer{address: address}
}

func (s *RPCServer) Start() error {
	rpc.Register(s)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("RPC server started on %s", s.address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (s *RPCServer) Upload(args *string, reply *string) error {
	*reply = "File uploaded successfully"
	return nil
}

func StartSecureRPCServer(address string, certFile string, keyFile string) error {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Secure RPC server started on %s", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

type RPCHandler struct {
}

func (h *RPCHandler) RegisterRPCHandler() {
	rpc.Register(h)
}
