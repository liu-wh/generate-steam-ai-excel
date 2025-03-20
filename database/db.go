package database

import (
	"context"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func SetUpDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SetUpRedis() error {
	global.R = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("Redis_Server"),
		DB:           0,
		Password:     os.Getenv("Redis_Password"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
	})

	if _, err := global.R.Ping(context.TODO()).Result(); err != nil {
		global.Logger.Error("redis连接失败", code.ERROR, err)
		return err
	}
	return nil

}
