package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port     string
	DataFile string
}

func Load() (Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be an integer between 1 and 65535")
	}

	dataFile := strings.TrimSpace(os.Getenv("DATA_FILE"))
	if dataFile == "" {
		dataFile = "data/bills.json"
	}

	return Config{
		Port:     strconv.Itoa(portNumber),
		DataFile: dataFile,
	}, nil
}
