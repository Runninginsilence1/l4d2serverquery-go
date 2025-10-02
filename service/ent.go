package service

import (
	"context"
	"fmt"
	"log"

	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/pkg/sqlite_driver"

	"github.com/duke-git/lancet/v2/fileutil"
)

var _client *ent.Client

func Client() *ent.Client {

	return _client
}

func InitClient() {
	sqlite_driver.Import()
	var err error
	dataSourceName := fmt.Sprintf("file:%s?cache=shared&_fk=1", dbPath)
	fmt.Println("使用的数据库文件路径为：", dbPath)

	if fileutil.IsExist(dbPath) {
		fmt.Println("数据库文件存在")
	} else {
		fmt.Println("数据库文件不存在, 注意!")
	}

	_client, err = ent.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := _client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
