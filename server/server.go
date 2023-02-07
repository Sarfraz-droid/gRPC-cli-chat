package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Sarfraz-droid/goBasics/db"
	"github.com/fatih/color"
)

type ServerConnections struct {
	connections []db.Connection
}

func (s *ServerConnections) AddConnection(c db.Connection) {
	s.connections = append(s.connections, c)
}

func (s *ServerConnections) DoesConnectionExist(port string) bool {
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
		color.Green("Port: %s", c.Port)
		fmt.Println()
	} 
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", ":80")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.UDPAddr)
	fmt.Println(remoteAddr.IP.String())
	return remoteAddr.IP.String()
}

func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
    var (
        ief      *net.Interface
        addrs    []net.Addr
        ipv4Addr net.IP
    )
    if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
        return
    }
    if addrs, err = ief.Addrs(); err != nil { // get addresses
        return
    }
    for _, addr := range addrs { // get ipv4 address
        if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
            break
        }
    }
    if ipv4Addr == nil {
        return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
    }
    return ipv4Addr.String(), nil
}

// LocalIP get the host machine local IP address
func LocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if IsPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

func IsPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
