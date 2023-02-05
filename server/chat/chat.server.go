package chatServer

import (
	"log"

	"github.com/Sarfraz-droid/goBasics/chat"
	"golang.org/x/net/context"
)

type Server struct {
	port int
	name string
}

type ConnectClient interface {
	
}

type Pool struct {
	cc func(int, string)
}

var pool = Pool{}

func UpdatePool(f func(int, string)) {
	pool.cc = f
}

func (s *Server) Setup(port int, name string) {
	s.port = port
	s.name = name
}

func (s *Server) mustEmbedUnimplementedHelloServiceServer() {
	panic("implement me")
}

func (s *Server) SayHello(ctx context.Context, message *chat.HelloRequest) (*chat.HelloReply, error) {
	log.Printf("Received message body from client: %s", message.Name)
	return &chat.HelloReply{Message: "Message Sent!"}, nil
}

func (s *Server) HandShake(ctx context.Context, message *chat.HandShakeRequest) (*chat.HelloReply, error) {
	log.Printf("Handshake to : %d %s", message.Port, message.Name)

	pool.cc(int(message.Port), message.Name)

	return &chat.HelloReply{
	Message: "Handshake successful",
	}, nil
}