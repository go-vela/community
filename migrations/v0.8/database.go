// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/server/database"

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
	Address       string
	EncryptionKey string

	Actions    *actions
	Connection *connection

	ConcurrencyLimit int

	Client database.Service
}

// New creates a client with the provided database.
func (d *db) New(c *cli.Context) error {
	logrus.Debug("creating database from provided configuration")

	// database configuration
	//
	// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup
	_setup := &database.Setup{
		Driver:           c.String("database.driver"),
		Address:          c.String("database.addr"),
		CompressionLevel: c.Int("database.compression.level"),
		ConnectionLife:   c.Duration("database.connection.life"),
		ConnectionIdle:   c.Int("database.connection.idle"),
		ConnectionOpen:   c.Int("database.connection.open"),
		EncryptionKey:    c.String("database.encryption.key"),
	}

	// setup the database
	//
	// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#New
	_database, err := database.New(_setup)
	if err != nil {
		return err
	}

	// update client with database service
	d.Client = _database

	return nil
}
