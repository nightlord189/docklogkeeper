package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/handler"
	"github.com/nightlord189/docklogkeeper/internal/log"
	pkgLog "github.com/nightlord189/docklogkeeper/pkg/log"
	"github.com/rs/zerolog"
	stdLog "log"
)

func main() {
	fmt.Println("start #1")

	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		stdLog.Fatalf("error on load config: %v", err)
	}

	if err := pkgLog.InitLogger(cfg.LogLevel); err != nil {
		stdLog.Fatalf("error on init log: %v", err)
	}

	ctx := context.Background()

	zerolog.Ctx(ctx).Debug().Msg("start #2")

	logAdapter, err := log.New(cfg.Log)
	if err != nil {
		stdLog.Fatalf("error on init log adapter: %v", err)
	}

	defer logAdapter.Close()

	dock, err := docker2.New(ctx, cfg, logAdapter)
	if err != nil {
		stdLog.Fatalf("error on init docker: %v", err)
	}

	defer dock.Close()

	go dock.Run(ctx)

	handlerInst := handler.New(cfg.HTTP, dock)

	if err := handlerInst.Run(); err != nil {
		stdLog.Fatalf("run router error: %v", err)
	}

	// TODO: on start ContainerLogs, get right Since param (based on last written file)
	// TODO: retention by date
	// TODO: retention by file size

	// TODO: auth
	// TODO: setting for beautifying container name
	// TODO: search
	// TODO: get
	// TODO: frontend
}
