package database

import (
	"SGMS/config/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// Init 函数用于初始化数据库连接
func Init() error {
	// 从配置中获取数据库连接所需的参数
	user := config.Config.GetString("GaussDB.User")
	pass := config.Config.GetString("GaussDB.Pass")
	port := config.Config.GetString("GaussDB.Port")
	host := config.Config.GetString("GaussDB.Host")
	name := config.Config.GetString("GaussDB.DBName")
	path := config.Config.GetString("GaussDB.SearchPath")

	// 构建 DSN
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v search_path=%v sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, pass, name, path)

	// 使用 GORM 打开数据库连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return fmt.Errorf("database connect failed: %w", err)
	}

	DB = db
	return nil
}
