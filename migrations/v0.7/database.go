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
// capture all logs from the database. After capturing all
// logs from the database, we'll update the log entry with
// a compressed data field.
func (d *db) Exec() error {
	logrus.Debug("executing workload from provided configuration")

	// TODO: remove this hack
	//
	// this allows us to "ignore" the error messages
	// returned from GetBuildLogs()
	//
	// capture current log level
	currentLevel := logrus.GetLevel()
	// only output panic level logs
	logrus.SetLevel(logrus.PanicLevel)

	logrus.Info("capturing all builds from the database")
	// capture all builds from the database
	builds, err := d.Client.GetBuildList()
	if err != nil {
		return err
	}

	// TODO: remove this hack
	//
	// this allows us to "ignore" the error messages
	// returned from GetBuildLogs()
	//
	// output intended level of logs
	logrus.SetLevel(currentLevel)

	// iterate through all builds from the database
	for _, build := range builds {
		logrus.Infof("capturing all logs for build %d", build.GetID())
		// capture all logs for the build from the database
		logs, err := d.Client.GetBuildLogs(build.GetID())
		if err != nil {
			return err
		}

		// iterate through all logs for the build from the database
		for _, log := range logs {
			// update log entry with compression in the database
			err = d.Client.UpdateLog(log)
			if err != nil {
				return err
			}
		}

		logrus.Debugf("all logs updated for build %d", build.GetID())
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
