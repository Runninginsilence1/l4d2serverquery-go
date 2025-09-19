package steamquery

import (
	"log"

	"l4d2serverquery-go/pkg/steamquery/parse_data"
)

func QueryMasterServer(serverName string, page, pageSize int) []parse_data.Server {
	data, err := GetDataWithName(serverName, page, pageSize)
	if err != nil {
		log.Println("查询服务器数据失败:", err)
		return nil
	}

	servers := Servers(data)

	return servers
}
