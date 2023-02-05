package server

import (
	"fmt"

	"github.com/Sarfraz-droid/goBasics/db"
	"github.com/fatih/color"
)

type ServerConnections struct {
	connections []db.Connection
}

func (s *ServerConnections) AddConnection(c db.Connection) {
	s.connections = append(s.connections, c)
}

func (s *ServerConnections) DoesConnectionExist(port int) bool {
	for _, c := range s.connections {
		if c.Port == port {
			return true
		}
	}
	return false
}

func DisplayPorts() {
	// display all the ports
	_db :=db.GetDB()
	connections := _db.GetConnections()

	for _, c := range connections {
		color.Red("Name: %s", c.Name)
		color.Green("Port: %d", c.Port)
		fmt.Println()
	} 
}