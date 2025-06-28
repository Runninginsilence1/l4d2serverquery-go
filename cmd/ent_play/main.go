package main

import (
	"context"
	"fmt"

	"l4d2serverquery-go/service"
)

func main() {
	client := service.Client()
	ctx := context.TODO()

	//tagsHN := client.Tag.Query().Where(tag.NameContains("HN")).OnlyX(ctx)

	tags := client.Tag.Query().AllX(ctx)

	for _, t := range tags {
		servers := t.QueryServers().AllX(ctx)
		fmt.Println(t.Name, len(servers))
	}

	//fmt.Println(tagsHN)

	//servers := tagsHN.QueryServers().AllX(ctx)
	//fmt.Println(servers)
	//fmt.Println(len(servers))
}
