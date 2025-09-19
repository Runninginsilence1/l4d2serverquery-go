package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FavoriteServer holds the schema definition for the FavoriteServer entity.
type FavoriteServer struct {
	ent.Schema
}

// Fields of the FavoriteServer.
func (FavoriteServer) Fields() []ent.Field {
	return []ent.Field{
		field.String("addr").
			Unique(),
		field.String("name").Optional(),
		field.String("desc").Optional(),
		//field.Time("created_at").Default(time.Now()),
		field.Time("last_query_time").Optional(), // 最后查询时间
		field.Int("rank").Default(7),             // 排序的权重, 在线玩家数
	}
}

// Edges of the FavoriteServer.
func (FavoriteServer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", Tag.Type).
			Ref("servers"),
	}
}
