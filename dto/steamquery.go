package dto

// 这里对标前端的结构

type SteamQuery struct {
	Addr   string `json:"addr"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Map    string `json:"map"`
}

type Server struct {
	ID                  int    `json:"id"`
	ServerName          string `json:"serverName"`
	Address             string `json:"address"`
	Map                 string `json:"map"`
	Version             string `json:"version"`
	OnlinePlayers       int    `json:"onlinePlayers"`
	MaxPlayers          int    `json:"maxPlayers"`
	LastQueryTimeString string `json:"lastQueryTimeString"`
	Rank                int    `json:"rank"`
	Tags                []Tag  `json:"tags"`
}
