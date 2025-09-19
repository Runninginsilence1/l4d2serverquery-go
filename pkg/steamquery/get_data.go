package steamquery

import (
	"fmt"
	"io"
	"net/http"

	"l4d2serverquery-go/pkg/steamquery/queryitem"
	"l4d2serverquery-go/pkg/steamquery/steamserverbrowser"
)

func GetDataWithName(serverName string, page int, pageSize int) (string, error) {
	// 很可惜url字符串不是通过 query参数来拼接的;
	url := fmt.Sprintf("%v/%v/%v", steamserverbrowser.GetBaseURL(), page, pageSize)
	resp, err := http.Post(url, "application/json",
		queryitem.BuildQueryConditions(
			queryitem.NewServerNameQueryItem(serverName), // server name limit,
		),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(all), nil
}

func GetData() (string, error) {
	// 很可惜url字符串不是通过 query参数来拼接的;
	url := "https://api.steamserverbrowser.com/v2/games/550/servers/query/AS/1/25"
	resp, err := http.Post(url, "application/json",
		queryitem.BuildQueryConditions(
			queryitem.NewServerNameQueryItem("芙芙"), // server name limit,
			//queryitem.NewCurrentPlayerCountQueryItem(8, queryitem.Equal),
			// ... add more query items here
		),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(all), nil
}
