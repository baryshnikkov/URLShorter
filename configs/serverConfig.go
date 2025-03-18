package configs

import (
	"flag"
)

type ServerConfig struct {
	Port string
	Ip   string
}

func LoadServerConfig() *ServerConfig {
	port := flag.String("port", ":8080", "Port to listen on")
	ip := flag.String("ip", "127.0.0.1", "IP to listen on")
	flag.Parse()

	return &ServerConfig{
		Port: *port,
		Ip:   *ip,
	}
}
