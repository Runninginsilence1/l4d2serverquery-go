package parse_data

import (
	"fmt"
	"time"
)

type Server struct {
	Id              int       `json:"id"`
	IpAddress       string    `json:"ipAddress"`
	Port            int       `json:"port"`
	Name            string    `json:"name"`
	CurrentPlayers  int       `json:"currentPlayers"`
	MaxPlayers      int       `json:"maxPlayers"`
	Map             string    `json:"map"`
	IsSecure        bool      `json:"isSecure"`
	CurrentBots     int       `json:"currentBots"`
	LastQueryTime   time.Time `json:"lastQueryTime"`
	OperatingSystem int       `json:"operatingSystem"`
	ServerCategory  int       `json:"serverCategory"`
	Version         string    `json:"version"`
	IsOnline        bool      `json:"isOnline"`
	IsFavourite     bool      `json:"isFavourite"`
	IsOutdated      bool      `json:"isOutdated"`
}

func (receiver Server) EasyInfo() string {
	return fmt.Sprintf("%s:%d (%s) - %d/%d players - %s", receiver.IpAddress, receiver.Port, receiver.Name, receiver.CurrentPlayers, receiver.MaxPlayers, receiver.Map)
}
