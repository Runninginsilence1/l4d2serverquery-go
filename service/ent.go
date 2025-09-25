package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/pkg/sqlite_driver"
)

var once sync.Once
var _client *ent.Client

func Client() *ent.Client {
	once.Do(func() {
		sqlite_driver.Import()
		var err error
		dataSourceName := fmt.Sprintf("file:%s?cache=shared&_fk=1", dbPath)
		fmt.Println("使用的数据库文件路径为：", dbPath)
		_client, err = ent.Open("sqlite3", dataSourceName)
		if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
		}

		if err := _client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	})

	return _client
}
