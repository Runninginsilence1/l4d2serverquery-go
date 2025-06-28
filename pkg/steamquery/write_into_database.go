package steamquery

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// gorm 的原生 sql

// GormMysqlDsnConfig 配置
// example: root:123456@tcp(192.168.3.102:3306)/sdic-vault2?charset=utf8mb4&parseTime=True&loc=Local
type GormMysqlDsnConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

func GormMysqlDsn(config GormMysqlDsnConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DbName)
}

var dbOnce sync.Once
var _db *gorm.DB

// mysql 连接
func GetMysqlDB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := GormMysqlDsn(GormMysqlDsnConfig{
			Username: "root",
			Password: "Test654321",
			Host:     "121.37.157.126",
			Port:     "14118",
			DbName:   "l4d2",
		})
		var err error
		_db, err = gorm.Open(
			mysql.Open(dsn),
			&gorm.Config{
				Logger: logger.Discard,
			},
		)
		if err != nil {
			panic(err)
		}
	})
	return _db
}

func GetPgDB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := "host=121.37.157.126 user=root password=5D47jyKt dbname=l4d2 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		var err error
		_db, err = gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{
				Logger: logger.Discard,
				NamingStrategy: schema.NamingStrategy{
					NoLowerCase: true,
				},
			},
		)
		if err != nil {
			panic(err)
		}
	})
	return _db
}

func WriteIntoDatabase(sqls []string) {
	db := GetPgDB()

	fmt.Println("服务器数量:", len(sqls))

	errCount := 0

	for _, sql := range sqls {
		err := db.Exec(sql).Error
		handleErrorWithLog(err, sql, &errCount)
	}

	if errCount > 0 {
		//fmt.Println("执行原始sql出错数量:", errCount)
		fmt.Fprintln(os.Stderr, "执行原始sql出错数量:", errCount)
	}
	return
}

func handleErrorWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrorWithLog(err error, sql string, count *int) {
	if err != nil {
		*count += 1
		fmt.Fprintln(os.Stderr, "执行原始sql出错:", sql, "\n", err)
	}
}
