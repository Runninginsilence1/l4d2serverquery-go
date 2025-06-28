package steamquery

import "time"

type FavoriteServer struct {
	Id uint `gorm:"primarykey"`

	CreateAt time.Time
	Host     string `gorm:"type:varchar(255);not null"`
	Port     int    `gorm:"not null"`
	Desc     string `gorm:"type:varchar(255);"`
	Addr     string `gorm:"type:varchar(255);not null"`
}

func (f FavoriteServer) TableName() string {
	return "FavoriteServers"
}
