// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

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
	Driver     string
	Config     string
	Connection *connection

	Client database.Service
}

// New creates a client with the provided database.
func (d *db) New(c *cli.Context) error {
	logrus.Debug("creating database from provided configuration")

	switch d.Driver {
	case constants.DriverPostgres:
		// creating database from provided configuration
		_database, err := database.New(c)
		if err != nil {
			return err
		}

		// update client with database service
		d.Client = _database
	case constants.DriverSqlite:
		// creating database from provided configuration
		_database, err := database.New(c)
		if err != nil {
			return err
		}

		// update client with database service
		d.Client = _database
	default:
		return fmt.Errorf("invalid database driver: %s", d.Driver)
	}

	return nil
}

// Exec takes the provided configuration and attempts to
// run a series of different functions that will
// manipulate indexes, tables and columns.
func (d *db) Exec() error {
	logrus.Debug("executing workload from provided configuration")

	// compress all log entries in the database
	err := d.Compress()
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the provided database is configured properly.
func (d *db) Validate() error {
	logrus.Debug("validating provided database configuration")

	// check if the database driver is set
	if len(d.Driver) == 0 {
		return fmt.Errorf("VELA_DATABASE_DRIVER is not properly configured")
	}

	// check if the database configuration is set
	if len(d.Config) == 0 {
		return fmt.Errorf("VELA_DATABASE_CONFIG is not properly configured")
	}

	return nil
}
