package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Connection struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

func (c Connection) GetPort() string {
	return c.Port
}

func (c Connection) GetName() string {
	return c.Name
}

type DB struct {
	list []Connection
	curr Connection
}


func (d *DB) UpdateJSON() {
	// update json file
	conn := db.list

	file, err := json.Marshal(conn)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(file))

	_ = ioutil.WriteFile("db.json", file, 0644)
}

func (d *DB) Connect(port int, name string) {
	// connect to db
}

func (d *DB) Disconnect(port int, name string) {
	// disconnect from db
}

func (d *DB) GetConnections() []Connection {
	return d.list
}

func (d *DB) GetConnection(port string, name string) Connection {
	for _, c := range d.list {
		if c.Port == port && c.Name == name {
			return c
		}
	}
	return Connection{}
}

func (d *DB) AddConnection(port string, name string) {
	d.list = append(d.list, Connection{port, name})

	log.Printf("Added connection: %v", d.list)

	d.UpdateJSON()
}

func (d *DB) RemoveConnection(port string, name string) {
	for i, c := range d.list {
		if c.Port == port && c.Name == name {
			d.list = append(d.list[:i], d.list[i+1:]...)
		}
	}
}

func (d *DB) fetchJSON() {
	jsonFile,err := os.Open("db.json")

	if err != nil {
		log.Fatalf("Failed to open db.json: %v", err)
	}

	defer jsonFile.Close()

	// read json file and add connections to db
	value, err := ioutil.ReadAll(jsonFile)
	
	var conn []Connection

	json.Unmarshal(value, &conn)

	if err != nil {
		log.Fatalf("Failed to read db.json: %v", err)
	}

	for i := 0; i < len(conn); i++ {
		db.AddConnection(conn[i].Port, conn[i].Name)
	}

}

var db = DB{}

func GetDB() *DB {
	return &db
}

func Initialize() {
	db.fetchJSON()
}

func (d *DB) GetCurrConnection() Connection {
	return d.curr
}

func (d *DB) SetCurrConnection(port string, name string) {
	d.curr = d.GetConnection(port, name)
}

func (d *DB) GetConnectionByPort(port string) Connection {
	for _, c := range d.list {
		if c.Port == port {
			return c
		}
	}
	return Connection{}
}