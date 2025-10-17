package service

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"l4d2serverquery-go/dto"
	"l4d2serverquery-go/ent"
	"l4d2serverquery-go/ent/tag"
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

	tag, err := client.Tag.Create().SetName(name).SetRank(rank).Save(ctx)
	return tag, err
}

func DeleteTag(id int) error {
	ctx := context.Background()
	client := Client()

	return client.Tag.DeleteOneID(id).Exec(ctx)
}
