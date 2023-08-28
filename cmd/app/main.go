package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/handler"
	"github.com/nightlord189/docklogkeeper/internal/log"
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

	logAdapter, err := log.New(cfg.Log)
	if err != nil {
		stdLog.Fatalf("error init log adapter: %v", err)
	}

	dock, err := docker2.New(ctx, cfg, logAdapter)
	if err != nil {
		stdLog.Fatalf("error on init docker: %v", err)
	}

	defer dock.Close()

	go dock.Run(ctx)

	usecaseInst := usecase.New(dock, logAdapter)

	handlerInst := handler.New(cfg, usecaseInst, logAdapter)

	/*go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("delayed test func")

		req := log.SearchRequest{Contains: "gin"}

		lines, err := logAdapter.SearchLines(ctx, "remindmenow", req)
		fmt.Printf("search %s: found %d lines, err %v\n", req.Contains, len(lines), err)
	}()*/

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
			if err := logAdapter.DeleteContainersWithoutLogs(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("delete containers without logs error")
			}
		}
	}()

	if err := handlerInst.Run(); err != nil {
		stdLog.Fatalf("run router error: %v", err)
	}

	// TODO: regular update of containers
	// TODO: regular update of logs
}
