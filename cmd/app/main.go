package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/handler"
	"github.com/nightlord189/docklogkeeper/pkg/log"
	"github.com/rs/zerolog"
	stdLog "log"
)

func main() {
	fmt.Println("start #1")

	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		stdLog.Fatalf("error on load config: %v", err)
	}

	if err := log.InitLogger(cfg.LogLevel); err != nil {
		stdLog.Fatalf("error on init logger: %v", err)
	}

	ctx := context.Background()

	zerolog.Ctx(ctx).Debug().Msg("start #2")

	dock := docker2.New(cfg)

	fmt.Println("containers", dock.GetContainers())

	handlerInst := handler.New(cfg.HTTP, dock)

	if err := handlerInst.Run(); err != nil {
		stdLog.Fatalf("run router error: %v", err)
	}
}
