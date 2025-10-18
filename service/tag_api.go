package service

import (
	"context"

	"l4d2serverquery-go/dto"
	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/ent/favoriteserver"
	"l4d2serverquery-go/ent/tag"

	"github.com/duke-git/lancet/v2/slice"
)

func Tags() ([]dto.Tag, error) {
	ctx := context.Background()
	client := Client()

	all, err := client.Tag.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	tags := slice.Map(all, func(index int, item *ent.Tag) dto.Tag {
		return dto.Tag{
			ID:   item.ID,
			Name: item.Name,
		}
	})

	return tags, nil
}

func TagsByIDs(ids []int) ([]dto.Tag, error) {
	ctx := context.Background()
	client := Client()

	all, err := client.Tag.Query().Where(tag.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, err
	}

	tags := slice.Map(all, func(index int, item *ent.Tag) dto.Tag {
		return dto.Tag{
			ID:   item.ID,
			Name: item.Name,
		}
	})

	return tags, nil
}

func AddTag(name string, rank int) error {
	ctx := context.Background()
	client := Client()

	_, err := client.Tag.Create().SetName(name).SetRank(rank).Save(ctx)
	return err
}

func CreateTag(name string, rank int) (*ent.Tag, error) {
	ctx := context.Background()
	client := Client()

	tagItem, err := client.Tag.Create().SetName(name).SetRank(rank).Save(ctx)
	return tagItem, err
}

func DeleteTag(id int) error {
	ctx := context.Background()
	client := Client()

	return client.Tag.DeleteOneID(id).Exec(ctx)
}

// BindServerTags 绑定服务器和标签
func BindServerTags(serverID int, tagIDs []int) error {
	ctx := context.Background()
	client := Client()

	// 查询服务器是否存在
	server, err := client.FavoriteServer.Get(ctx, serverID)
	if err != nil {
		return err
	}

	// 查询标签是否存在
	tags, err := client.Tag.Query().Where(tag.IDIn(tagIDs...)).All(ctx)
	if err != nil {
		return err
	}

	// 绑定标签（会清除原有的标签并添加新的）
	_, err = server.Update().ClearTags().AddTags(tags...).Save(ctx)
	return err
}

// GetServerTags 获取服务器的所有标签
func GetServerTags(serverID int) ([]dto.Tag, error) {
	ctx := context.Background()
	client := Client()

	// 查询服务器及其标签
	server, err := client.FavoriteServer.Query().
		Where(favoriteserver.ID(serverID)).
		WithTags().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	tags := slice.Map(server.Edges.Tags, func(index int, item *ent.Tag) dto.Tag {
		return dto.Tag{
			ID:   item.ID,
			Name: item.Name,
		}
	})

	if tags == nil {
		tags = make([]dto.Tag, 0)
	}

	return tags, nil
}
