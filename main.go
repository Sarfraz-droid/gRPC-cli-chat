package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/Sarfraz-droid/goBasics/chat"
	"github.com/Sarfraz-droid/goBasics/client"
	"github.com/Sarfraz-droid/goBasics/db"
	"github.com/Sarfraz-droid/goBasics/server"
	chatServer "github.com/Sarfraz-droid/goBasics/server/chat"
	tm "github.com/buger/goterm"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
)



func main() {
	port, err := freeport.GetFreePort()
	db.Initialize();
	_db := db.GetDB()
	

	if err != nil {
		log.Fatalf("Failed to get free port: %v", err)
	}

	lis, err := net.Listen("tcp", ":" + fmt.Sprintf("%d", port))

	if err != nil { 
		log.Fatalf("Failed to listen to %d: %v",port, err)
	}
	// fmt.Printf("IP address: %s",add)
	// add, err := server.GetInterfaceIpv4Addr("en0")

	// if err != nil {
	// 	log.Fatalf("Failed to get IP address: %v", err)
	// }

	// fmt.Printf("IP address: %s",add)
	ip, err := server.LocalIP()
	if err != nil {
		log.Fatalf("Failed to get IP address: %v", err)
	}
	fmt.Printf("IP address: %s",ip)
	fmt.Println()
	full_address := fmt.Sprintf("%s:%d", ip ,port)

	// fmt.Println(server.)
	// fmt.Println(server.GetOutboundIP())
	// fmt.Println(full_address)
	fmt.Print("Enter your name: ")

	var name string;
	fmt.Scanln(&name)
	_db.AddConnection(full_address, name)
	_db.SetCurrConnection(full_address, name)
	
	chatServer.UpdatePool(client.InitializeChat)
	s := chatServer.Server{}
	s.Setup(port, name)

	grpcServer := grpc.NewServer()
	tm.Clear();
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()

	log.Printf("Server listening on port %d" , port)
	go chat.RegisterHelloServiceServer(grpcServer,&s);
	server.DisplayPorts();
	go client.Input();

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}