// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/simukti/grpc-password/server"
	"github.com/simukti/grpc-password/storage"
	"github.com/simukti/grpc-password/storage/memory"
	"github.com/simukti/grpc-password/storage/mysql"
)

func main() {
	hostPort := os.Getenv("APP_HOST_PORT")
	if hostPort == "" {
		hostPort = ":54123"
	}

	lst, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalln(err)
	}

	db := storageFactory(
		os.Getenv("APP_DB_TYPE"),
		os.Getenv("APP_DB_DSN"),
	)

	grpcServer := server.New(db)
	go func() {
		if err := grpcServer.Serve(lst); err != nil {
			log.Fatalln(err)
		}
	}()

	select {}
}

func storageFactory(driver, dsn string) storage.Storage {
	if driver == "" {
		driver = "memory"
	}

	switch strings.ToLower(driver) {
	case "memory":
		return memory.Storage()
	case "mysql":
		if dsn == "" {
			log.Fatalln("MySQL driver required dsn string env var")
		}
		return mysql.Storage(driver, dsn)
	default:
		log.Fatalf("Invalid storage driver: %s", driver)
	}

	return nil
}
