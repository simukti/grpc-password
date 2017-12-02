// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package mysql

import (
	"log"
	"time"

	"github.com/GuiaBolso/darwin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/simukti/grpc-password/repository"
	"github.com/simukti/grpc-password/storage"
)

type mysqlStorage struct {
	db *dbr.Session
}

// Storage create mysql storage object
func Storage(driver, dsn string) storage.Storage {
	con, _ := dbr.Open(driver, dsn, nil)
	sess := con.NewSession(nil)
	sess.SetMaxIdleConns(0)
	sess.SetMaxOpenConns(8)
	sess.SetConnMaxLifetime(time.Minute * 5)

	st := &mysqlStorage{
		db: sess,
	}

	if err := st.initDB(); err != nil {
		log.Fatalf("MySQL storage init error: %s", err.Error())
	}

	return st
}

func (c *mysqlStorage) initDB() error {
	if err := c.db.Ping(); err != nil {
		return err
	}

	drv := darwin.NewGenericDriver(c.db.DB, darwin.MySQLDialect{})
	m := darwin.New(drv, migrations(), nil)

	if err := m.Migrate(); err != nil {
		return err
	}

	info, err := m.Info()
	if err != nil {
		return err
	}

	log.Printf("MySQL storage")
	for _, mi := range info {
		log.Printf("v%d: %s", int64(mi.Migration.Version), mi.Migration.Description)
	}

	return nil
}

// Password mysql PasswordRepository
func (c *mysqlStorage) Password() repository.PasswordRepository {
	return repository.PasswordMysqlStorage(c.db)
}
