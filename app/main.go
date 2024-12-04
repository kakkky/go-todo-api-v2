package main

import (
	"context"
	"log"

	"github.com/kakkky/app/config"
	"github.com/kakkky/app/infrastructure/db/mysql"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read config : %v", err)
	}
	if err := run(ctx, cfg); err != nil {
		// エラー処理
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	// DB接続
	close := mysql.NewDB(ctx, cfg)
	defer close()
	// サーバー
	return nil
}
