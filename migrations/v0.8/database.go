// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/postgres"
	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// actions represents all potential actions
// this utility can invoke with the database.
type actions struct {
	All          bool
	AlterTables  bool
	EncryptUsers bool
	SyncCounter  bool
}

// connection represents all connection related
// information used to communicate with the
// provided database.
type connection struct {
	Idle int
	Life time.Duration
	Open int
}

// db represents all database related
// information used to communicate
// with the database.
type db struct {
	Driver        string
	Config        string
	EncryptionKey string

	Actions    *actions
	Connection *connection

	ConcurrencyLimit int

	Client database.Service
}

// New creates a client with the provided database.
func (d *db) New(c *cli.Context) error {
	logrus.Debug("creating database from provided configuration")

	switch d.Driver {
	case constants.DriverPostgres:
		// create new Postgres database service
		//
		// https://pkg.go.dev/github.com/go-vela/server/database/postgres?tab=doc#New
		_database, err := postgres.New(
			postgres.WithAddress(c.String("database.config")),
			postgres.WithCompressionLevel(c.Int("database.compression.level")),
			postgres.WithConnectionLife(c.Duration("database.connection.life")),
			postgres.WithConnectionIdle(c.Int("database.connection.idle")),
			postgres.WithConnectionOpen(c.Int("database.connection.open")),
			postgres.WithEncryptionKey(c.String("database.encryption.key")),
		)
		if err != nil {
			return err
		}

		// update client with postgres database service
		d.Client = _database
	case constants.DriverSqlite:
		// create new Sqlite database service
		//
		// https://pkg.go.dev/github.com/go-vela/server/database/sqlite?tab=doc#New
		_database, err := sqlite.New(
			sqlite.WithAddress(c.String("database.config")),
			sqlite.WithCompressionLevel(c.Int("database.compression.level")),
			sqlite.WithConnectionLife(c.Duration("database.connection.life")),
			sqlite.WithConnectionIdle(c.Int("database.connection.idle")),
			sqlite.WithConnectionOpen(c.Int("database.connection.open")),
			sqlite.WithEncryptionKey(c.String("database.encryption.key")),
		)
		if err != nil {
			return err
		}

		// update client with sqlite database service
		d.Client = _database
	default:
		return fmt.Errorf("invalid database driver: %s", d.Driver)
	}

	return nil
}
