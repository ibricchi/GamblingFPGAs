package main

import (
	"context"
	"flag"
	"log"

	"github.com/IBricchi/GamblingFPGAs/server"
	"go.uber.org/zap"
)

func main() {
	var serverDBFilePath = flag.String("db", "serverDB.db", "SQLite DB file name")
	flag.Parse()

	dbDSN := "db/" + *serverDBFilePath

	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("credentials: failed to create zap logger: %v\n", err)
	}
	defer logger.Sync()

	db, err := server.OpenSQLiteDB(ctx, logger, dbDSN)
	if err != nil {
		logger.Fatal("credentials: failed to open SQLite server database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("credentials: opened sqlite3 DB")

	if err := server.AddCredential(ctx, db); err != nil {
		logger.Fatal("credentials: failed to add credentials", zap.Error(err))
	}
	logger.Info("credentials: finished adding credentials")
}
