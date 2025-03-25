package configs

import (
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type ServerConfig struct {
	Port string
	Ip   string
}

func LoadServerConfig() *ServerConfig {
	portFlag := flag.String("port", "", "Port to listen on")
	ipFlag := flag.String("ip", "", "IP to listen on")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Fatal("Error loading .env file")
	}

	portEnv := os.Getenv("PORT")
	ipEnv := os.Getenv("IP")

	var port string
	switch {
	case *portFlag != "":
		port = *portFlag
	case portEnv != "":
		port = portEnv
	default:
		port = ":8080"
	}

	var ip string
	switch {
	case *ipFlag != "":
		ip = *ipFlag
	case ipEnv != "":
		ip = ipEnv
	default:
		ip = "127.0.0.1"
	}

	return &ServerConfig{
		Port: port,
		Ip:   ip,
	}
}
