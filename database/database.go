package database

import (
	"context"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/chpool"
	"github.com/gofiber/fiber/v2/log"
	"void-studio.net/fiesta/fapi/config"
)

var ctx = context.Background()
var pool *chpool.Pool

func init() {
	var err error = nil
	pool, err = chpool.Dial(ctx, chpool.Options{
		MaxConns: 4,
		ClientOptions: ch.Options{
			Address:  config.Config.Database.Address,
			Database: config.Config.Database.Username,
			User:     config.Config.Database.Username,
			Password: config.Config.Database.Password,
		},
	})

	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}
}
