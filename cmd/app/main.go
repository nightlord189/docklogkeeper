package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/nightlord189/docklogkeeper/internal/handler"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/nightlord189/docklogkeeper/internal/repo"
	"github.com/nightlord189/docklogkeeper/internal/trigger"
	"github.com/nightlord189/docklogkeeper/internal/usecase"
	pkgLog "github.com/nightlord189/docklogkeeper/pkg/log"
	"github.com/rs/zerolog"
	stdLog "log"
	"time"
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

	repoInst, err := repo.New(cfg.DB)
	if err != nil {
		stdLog.Fatalf("error init repo: %v", err)
	}

	triggerAdapter := trigger.New(repoInst)

	go triggerAdapter.Run(ctx)

	logAdapter := log.New(cfg.Log, repoInst, []chan entity.LogDataDB{triggerAdapter.LogsChan})

	dock, err := docker2.New(ctx, cfg, logAdapter)
	if err != nil {
		stdLog.Fatalf("error on init docker: %v", err)
	}

	defer dock.Close(ctx)

	go dock.Run(ctx)

	usecaseInst := usecase.New(repoInst, dock, logAdapter, triggerAdapter)

	handlerInst := handler.New(cfg, repoInst, usecaseInst, logAdapter)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		logAdapter.ClearOldFiles(ctx)
		for range ticker.C {
			logAdapter.ClearOldFiles(ctx)
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			if err := repoInst.DeleteContainersWithoutLogs(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("delete containers without logs error")
			}
		}
	}()

	if err := handlerInst.Run(); err != nil {
		stdLog.Fatalf("run router error: %v", err)
	}

	// TODO: regular update of logs
}
