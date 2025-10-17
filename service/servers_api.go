package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"l4d2serverquery-go/dto"
	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/ent/favoriteserver"
	"l4d2serverquery-go/ent/predicate"
	"l4d2serverquery-go/ent/tag"
	"l4d2serverquery-go/logger"
	"l4d2serverquery-go/pkg/steamquery/parse_data"
	"l4d2serverquery-go/singleflight"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
)

func AddServer(addr string) error {
	ctx := context.Background()
	client := Client()

	_, err := client.FavoriteServer.Create().SetAddr(addr).Save(ctx)
	return err
}

func CreateExampleServer() {
	servers := []string{
		"202.189.9.212:10207",
		"202.189.9.212:10208",
		"110.42.44.3:27022",
	}

	for _, addr := range servers {
		_ = AddServer(addr)
	}
}

func CreateExampleTag() {
	servers := []string{
		"芙芙",
		"HN",
		"萌聚",
		"萌新聚集地",
	}

	for _, addr := range servers {
		_ = AddTag(addr, 7) // 默认先用7来代替
	}
}

func Servers() ([]*ent.FavoriteServer, error) {
	ctx := context.Background()
	client := Client()

	servers, err := client.FavoriteServer.Query().All(ctx)
	return servers, err
}

func DeleteServer(id int) error {
	ctx := context.Background()

	client := Client()
	err := client.FavoriteServer.DeleteOneID(id).Exec(ctx)
	return err
}

func PatchServer(id int, item dto.Server) error {
	client := Client()
	ctx := context.Background()

	// 修改的字段在这里标出
	// 现在暂时只有修改rank的需求

	_, err := client.FavoriteServer.UpdateOneID(id).
		SetRank(item.Rank).
		Save(ctx)
	return err
}

func QueryServers(tagIds []int, name string) ([]dto.Server, error) {
	name = strings.TrimSpace(name)
	logger.Log.Info("QueryServers arg:", "name", lo.Ternary(name == "", "empty", name))
	result, err, _ := singleflight.Sf().Do("servers", func() (interface{}, error) {
		var dtos []dto.Server
		wg := conc.NewWaitGroup()
		client := Client()
		ctx := context.Background()
		serverCond := []predicate.FavoriteServer{}

		if name != "" {
			serverCond = append(serverCond, favoriteserver.NameContainsFold(name))
		}

		if len(tagIds) > 0 {
			serverCond = append(serverCond, favoriteserver.HasTagsWith(tag.IDIn(tagIds...)))
		}

		// 查询所有在tagIds中的服务器
		all, err := client.FavoriteServer.Query().Where(
			serverCond...,
		).
			WithTags().
			All(ctx)

		if err != nil {
			return nil, err
		}

		fmt.Printf("共查询到%v个结果! 即将对他们进行测试连接...\n", len(all))

		all = slice.UniqueBy(all, func(item *ent.FavoriteServer) string {
			return item.Addr
		})

		for _, server := range all {
			wg.Go(func() {
				info, err := Query(server.Addr)
				if err != nil {
					return
				}
				_, err = client.FavoriteServer.UpdateOneID(server.ID).SetName(info.Name).Save(ctx)
				if err != nil {
					logger.Log.Error("更新服务器名称失败", err)
					return
				}
				d := newServerDto(server, info)
				dtos = append(dtos, d)
			})
		}

		wg.Wait()

		if dtos == nil {
			dtos = make([]dto.Server, 0)
		}
		logger.Log.Info("所有服务器测试成功")
		slice.SortBy(dtos, func(a, b dto.Server) bool {
			priorityA := 0
			priorityB := 0

			rankA := a.Rank
			rankB := b.Rank

			priorityA = mathutil.Abs(rankA - a.OnlinePlayers)
			priorityB = mathutil.Abs(rankB - b.OnlinePlayers)

			return priorityA < priorityB
		})
		return dtos, nil
	})

	return result.([]dto.Server), err
}

func QuerySingleServer(serverID int) (any, error) {
	logger.Log.Info("single server id:", serverID)
	result, err, _ := singleflight.Sf().Do("servers", func() (interface{}, error) {
		client := Client()
		ctx := context.Background()

		server, err := client.FavoriteServer.Query().Where(favoriteserver.ID(serverID)).Only(ctx)

		if err != nil {
			return nil, err
		}

		fmt.Println("成功查询到指定ID的服务器")

		info, err := Query(server.Addr)
		if err != nil {
			return nil, err
		}
		_, err = client.FavoriteServer.UpdateOneID(server.ID).SetName(info.Name).Save(ctx)
		if err != nil {
			logger.Log.Error("更新服务器名称失败", err)
			return nil, err
		}

		d := newServerDto(server, info)

		return d, nil
	})

	return result, err
}

func newServerDto(item *ent.FavoriteServer, info L4d2SeverInfo) dto.Server {
	var r dto.Server
	r.ID = item.ID
	r.Address = item.Addr
	r.ServerName = info.Name
	r.Map = info.Map
	r.Version = info.Version
	r.OnlinePlayers = info.OnlinePlayers
	r.MaxPlayers = info.MaxPlayers
	//r.BotPlayers = info.BotPlayers
	r.Rank = item.Rank

	// 最后查询时间
	l := item.LastQueryTime
	if l.IsZero() {
		r.LastQueryTimeString = "从未连接"
	} else {
		r.LastQueryTimeString = item.LastQueryTime.Format("2006-01-02 15:04:05")
	}
	return r
}

func CleanServers() error {
	wg := conc.NewWaitGroup()

	client := Client()
	ctx := context.Background()

	all, err := client.FavoriteServer.Query().All(ctx)
	if err != nil {
		return err
	}

	for _, server := range all {
		wg.Go(func() {
			_, err := Query(server.Addr)
			if err == nil {
				return
			}

			client.FavoriteServer.DeleteOne(server)
		})
	}

	wg.Wait()

	return nil
}

func GroupServers() error {
	wg := conc.NewWaitGroup()

	client := Client()
	ctx := context.Background()

	servers, err := client.FavoriteServer.Query().All(ctx)
	if err != nil {
		return err
	}

	tags, err := client.Tag.Query().All(ctx)
	if err != nil {
		return err
	}

	for _, server := range servers {
		wg.Go(func() {
			info, err := Query(server.Addr)
			if err != nil {
				return
			}

			for _, t := range tags {
				serverName := strings.ToLower(info.Name)
				tagName := strings.ToLower(t.Name)

				if strings.Contains(serverName, tagName) {
					fmt.Printf("将 %v 关联到标签 %v\n", info.Name, t.Name)
					server.Update().AddTags(t).Save(ctx)
				}
			}
		})
	}

	wg.Wait()

	return nil
}

func UpdateLastCopyTime(id int) error {
	ctx := context.Background()
	client := Client()

	_, err := client.FavoriteServer.UpdateOneID(id).SetLastQueryTime(time.Now()).Save(ctx)
	return err
}

func DeleteServerById(id int) error {
	ctx := context.Background()
	client := Client()

	return client.FavoriteServer.DeleteOneID(id).Exec(ctx)
}

func RemoveDuplicateServers(servers []parse_data.Server) []parse_data.Server {
	ctx := context.Background()
	client := Client()

	servers = lo.Filter(servers, func(item parse_data.Server, index int) bool {
		_, err := client.FavoriteServer.Query().Where(
			favoriteserver.AddrEQ(fmt.Sprintf("%s:%d", item.IpAddress, item.Port)),
		).OnlyID(ctx)

		return err != nil
	})

	return servers
}
