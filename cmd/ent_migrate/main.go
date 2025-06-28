package main

import (
	"context"
	"log"

	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/pkg/sqlite_driver"
)

// 迁移数据库数据

func main() {
	sqlite_driver.Import()
	var err error
	var client *ent.Client
	client, err = ent.Open("sqlite3", "file:db.sqlite3?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	} else {
		log.Println("schema resources created successfully")
	}
}
