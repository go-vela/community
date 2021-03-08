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
	Driver           string
	Config           string
	CompressionLevel int
	EncryptionKey    string
	Connection       *connection

	BuildLimit       int
	ConcurrencyLimit int
	SecretLimit      int

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
func (d *db) Exec(c *cli.Context) error {
	logrus.Debug("executing workload from provided configuration")

	// alter required tables in the database
	err := d.Alter()
	if err != nil {
		return err
	}

	// drop unused indexes in the database
	err = d.Drop()
	if err != nil {
		return err
	}

	// create new database service
	err = d.New(c)
	if err != nil {
		return err
	}

	// compress all log entries in the database
	err = d.Compress()
	if err != nil {
		return err
	}

	// encrypt all secret values in the database
	err = d.Encrypt()
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

	// enforce AES-256, so check explicitly for 32 bytes on the key
	//
	// nolint: gomnd // ignore magic number
	if len(d.EncryptionKey) != 32 {
		return fmt.Errorf("VELA_DATABASE_ENCRYPTION_KEY invalid length specified: %d", len(d.EncryptionKey))
	}

	// check if the database build limit is set
	if d.BuildLimit < 0 {
		return fmt.Errorf("VELA_BUILD_LIMIT is not properly configured")
	}

	// check if the database concurrency limit is set
	if d.ConcurrencyLimit < 1 {
		return fmt.Errorf("VELA_CONCURRENCY_LIMIT is not properly configured")
	}

	// check if the database secret limit is set
	if d.SecretLimit < 0 {
		return fmt.Errorf("VELA_SECRET_LIMIT is not properly configured")
	}

	// check if the compression level is valid
	switch d.CompressionLevel {
	case constants.CompressionNegOne:
		fallthrough
	case constants.CompressionZero:
		fallthrough
	case constants.CompressionOne:
		fallthrough
	case constants.CompressionTwo:
		fallthrough
	case constants.CompressionThree:
		fallthrough
	case constants.CompressionFour:
		fallthrough
	case constants.CompressionFive:
		fallthrough
	case constants.CompressionSix:
		fallthrough
	case constants.CompressionSeven:
		fallthrough
	case constants.CompressionEight:
		fallthrough
	case constants.CompressionNine:
		break
	default:
		return fmt.Errorf("database compression level of '%d' is unsupported", d.CompressionLevel)
	}

	return nil
}
