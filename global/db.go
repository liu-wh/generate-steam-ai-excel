package global

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log/slog"
)

var (
	DB     *gorm.DB
	Logger *slog.Logger
	R      *redis.Client
	CTX    = context.Background()
	F      *excelize.File
)
