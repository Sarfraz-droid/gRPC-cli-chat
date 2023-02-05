package client

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Sarfraz-droid/goBasics/chat"
	"github.com/Sarfraz-droid/goBasics/db"
	"github.com/fatih/color"
)

var pool = Pool{}

func Input() {
	var port int

	color.Red("Enter the port number of the server you want to connect to: ")
	fmt.Scanln(&port)

	if(pool.isConnected) {
		log.Printf("Already connected to a server")
		return
	}

	color.Red("Connecting to port %d", port)

	conn := db.GetDB().GetConnectionByPort(port)


	InitializeChat(conn.Port, conn.Name)
}


func InitializeChat(
	port int,
	name string,
) {
	if pool.DoesConnectionExist(port) {
		log.Printf("Connection already exists")
		return
	}

	pool.Connect(port, name)

	go pool.SendMessage()

	// defer pool.conn.Close()
}

type Pool struct {
	isConnected bool
	port int
	name string
	conn *grpc.ClientConn
	client chat.HelloServiceClient
	connectedPorts []int
}

func (d *Pool) DoesConnectionExist(port int) bool {
	for _, c := range d.connectedPorts {
		if c == port {
			return true
		}
	}
	return false
}

func (d *Pool) SendMessage()  {
	var message string

	color.Red("Enter the message you want to send: ")
	fmt.Scanln(&message)


	msg := chat.HelloRequest{
		Name: message,
	}

	response, err := d.client.SayHello(context.Background(), &msg)

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Message)

	go d.SendMessage()
}

func (d *Pool) Connect(port int, name string) {
	d.port = port
	d.name = name

	log.Printf("Dialing to port %d", port)
	
	conn, err := grpc.Dial(":" + fmt.Sprint(port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}


	_db := db.GetDB()

	d.conn = conn
	d.client = chat.NewHelloServiceClient(conn)
	message := chat.HandShakeRequest{
		Port: int32(_db.GetCurrConnection().Port),
		Name: _db.GetCurrConnection().Name,
	}

	d.connectedPorts = append(d.connectedPorts, port)
	d.isConnected = true
	log.Printf("Sending handshake to server: %d %s", message.Port, message.Name)

	response, err := d.client.HandShake(context.Background(), &message)

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Message)
}
