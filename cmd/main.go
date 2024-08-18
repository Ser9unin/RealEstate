package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ser9unin/RealEstate/internal/config"
	"github.com/Ser9unin/RealEstate/internal/logger"
	"github.com/Ser9unin/RealEstate/internal/server"
	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := logger.NewLogger()
	serverCfg := config.NewSeverCfg()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	dsn := config.NewDBCfg()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to start db: %v", zap.Error(err)))
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("context cancelled: %v", zap.Error(err)))
	}

	storage := repository.New(db)

	server := server.NewServer(serverCfg, logger, storage)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server.Run()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Stop(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
